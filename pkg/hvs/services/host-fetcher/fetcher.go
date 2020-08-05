/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package hostfetcher

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/controllers"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/domain"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/common/chnlworkq"
	commLog "github.com/intel-secl/intel-secl/v3/pkg/lib/common/log"
	hc "github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector/types"
	"github.com/intel-secl/intel-secl/v3/pkg/model/hvs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	defaultRetryIntervalMins = 5
)

var defaultLog = commLog.GetDefaultLogger()

type retryRequest struct {
	retryTime time.Time
	hostId    uuid.UUID
}

type fetchRequest struct {
	ctx   context.Context
	host  hvs.Host
	rcvrs []domain.HostDataReceiver
}

type Service struct {
	Fetcher domain.HostDataFetcher
	// request channel is used to route requests into internal queue
	rqstChan chan interface{}
	// work items (their id) is pulled out of a queue and fed to the workers
	workChan chan interface{}

	retryRqstChan chan interface{}
	retryWorkChan chan interface{}

	// map that holds all hosts that need to be fetched.
	// The reason this is a map is that redundant requests can come in
	// that could theoretically be consolidated
	workMap map[uuid.UUID][]*fetchRequest
	// map is not protected access concurrent goroutine access. need a lock
	wmLock sync.Mutex
	// waitgroup used to wait for workers to finish up when signal for shutdown comes in
	wg sync.WaitGroup

	quit              chan struct{}
	serviceDone       bool
	retryIntervalMins int
	hcCfg             domain.HostConnectionConfig
	hcf               hc.HostConnectorFactory
	hs                domain.HostStatusStore
}

func NewService(cfg domain.HostDataFetcherConfig, workers int) (*Service, domain.HostDataFetcher, error) {
	defaultLog.Trace("hostfetcher/Service:New() Entering")
	defer defaultLog.Trace("hostfetcher/Service:New() Leaving")

	// setting size of channel to the same as number of workers.
	// this way, go routine can start work as soon as a current work is done
	svc := &Service{workMap: make(map[uuid.UUID][]*fetchRequest),
		quit:              make(chan struct{}),
		hcf:               cfg.HostConnectorFactory,
		retryIntervalMins: cfg.RetryTimeMinutes,
		hs:                cfg.HostStatusStore,
		hcCfg:             cfg.HostConnectionConfig,
	}
	if svc.hs == nil {
		return nil, nil, errors.New("host status store cannot be empty")
	}
	if svc.retryIntervalMins == 0 {
		svc.retryIntervalMins = defaultRetryIntervalMins
	}

	svc.Fetcher = svc
	var err error
	if svc.rqstChan, svc.workChan, err = chnlworkq.New(workers, workers, svc.addWorkToMap, nil, svc.quit, &svc.wg); err != nil {
		return nil, nil, errors.New("hostfetcher:NewService:error starting work queue")
	}
	if svc.retryRqstChan, svc.retryWorkChan, err = chnlworkq.New(workers, workers, nil, nil, svc.quit, &svc.wg); err != nil {
		return nil, nil, errors.New("hostfetcher:NewService:error starting retry queue")
	}

	// start workers.. individual workers are spawned as go routines
	svc.startWorkers(workers)
	svc.startRetryChannelProcessor(cfg.RetryTimeMinutes)
	return svc, svc.Fetcher, nil
}

// Function to Shutdown service. Will wait for pending host data fetch jobs to complete
// Will not process any further requests. Calling interface Async methods will result in error
func (svc *Service) Shutdown() error {
	defaultLog.Trace("hostfetcher/Service:Shutdown() Entering")
	defer defaultLog.Trace("hostfetcher/Service:Shutdown() Leaving")

	svc.serviceDone = true
	close(svc.quit)
	svc.wg.Wait()

	return nil
}

func (svc *Service) startRetryChannelProcessor(retryMins int) {
	defaultLog.Trace("hostfetcher/Service:startRetryChannelProcessor() Entering")
	defer defaultLog.Trace("hostfetcher/Service:startRetryChannelProcessor() Leaving")

	// start worker go routines
	svc.wg.Add(1)
	go func() {
		defer svc.wg.Done()
		for {
			select {
			case <-svc.quit:
				return
			case r := <-svc.retryWorkChan:
				retry := r.(retryRequest)
				select {
				case <-svc.quit:
					return
				case <-time.After(retry.retryTime.Sub(time.Now())):
					svc.workChan <- retry.hostId
				}
			}
		}
	}()
}

func (svc *Service) startWorkers(workers int) {
	defaultLog.Trace("hostfetcher/Service:startWorkers() Entering")
	defer defaultLog.Trace("hostfetcher/Service:startWorkers() Leaving")

	// start worker go routines
	for i := 0; i < workers; i++ {
		svc.wg.Add(1)
		go svc.doWork()
	}
}

// function used to add work to the map. If there is a current entry
// append the new request to the already queued up requests
func (svc *Service) addWorkToMap(wrk interface{}) interface{} {
	defaultLog.Trace("hostfetcher/Service:addWorkToMap() Entering")
	defer defaultLog.Trace("hostfetcher/Service:addWorkToMap() Leaving")

	switch v := wrk.(type) {
	case *fetchRequest:

		svc.wmLock.Lock()
		if _, ok := svc.workMap[v.host.Id]; !ok {
			svc.workMap[v.host.Id] = append(svc.workMap[v.host.Id], v)
		} else {
			svc.workMap[v.host.Id] = []*fetchRequest{v}
		}
		svc.wmLock.Unlock()
		return v.host.Id

	case retryRequest:
		return v.hostId
	default:
		log.Error("unexpected type in request channel")
		return nil
	}

}

// function that does the actual work. Receives id of host through work channel
// then pull records from the map, proceed to work unless requests are not already
// cancelled
func (svc *Service) doWork() {
	defaultLog.Trace("hostfetcher/Service:doWork() Entering")
	defer defaultLog.Trace("hostfetcher/Service:doWork() Leaving")

	defer svc.wg.Done()

	// receive id of queued work over the channel.
	// Fetch work context from the map.
	for {
		select {
		case <-svc.quit:
			// we have received a quit. Don't process anymore items - just return
			return
		case id := <-svc.workChan:
			hId, ok := id.(uuid.UUID)
			var connUrl string
			if !ok {
				defaultLog.Error("hostfetcher:doWork:expecting uuid from channel - but got different type")
			}
			// iterate through work requests for this host. Usually, there will only be a single element in the
			// work list.
			svc.wmLock.Lock()
			frs := svc.workMap[hId]
			connUrl = frs[0].host.ConnectionString
			getData := false
			for i, req := range frs {
				select {
				// remove the requests that have already been cancelled.
				case <-req.ctx.Done():
					frs = append(frs[:i], frs[i+1:]...)
					continue
				default:
					getData = true
				}
			}
			svc.workMap[hId] = frs
			svc.wmLock.Unlock()

			if getData {
				svc.FetchDataAndRespond(hId, connUrl)
			} else {
				defaultLog.Info("Fetch data for ", hId, "cancelled")
			}
		}
	}
}

func (svc *Service) Retrieve(ctx context.Context, host hvs.Host) (*types.HostManifest, error) {
	defaultLog.Trace("hostfetcher/Service:Retrieve() Entering")
	defer defaultLog.Trace("hostfetcher/Service:Retrieve() Leaving")

	hostData, err := svc.GetHostData(host.ConnectionString)
	hostStatus := &hvs.HostStatus{
		HostID: host.Id,
		HostStatusInformation: hvs.HostStatusInformation{
			LastTimeConnected: time.Now(),
		},
	}
	if err != nil {
		hostStatus.HostStatusInformation.HostState = hvs.HostStateUnknown
		if err := svc.hs.Persist(hostStatus); err != nil {
			defaultLog.Error("could not update host status to store")
		}
		return nil, err
	}
	hostStatus.HostStatusInformation.HostState = hvs.HostStateConnected
	hostStatus.HostManifest = *hostData

	if err := svc.hs.Persist(hostStatus); err != nil {
		defaultLog.Error("could not update host status and manifest to store")
	}

	return hostData, nil
}

func (svc *Service) RetriveAsync(ctx context.Context, host hvs.Host, rcvrs ...domain.HostDataReceiver) error {
	defaultLog.Trace("hostfetcher/Service:RetriveAsync() Entering")
	defer defaultLog.Trace("hostfetcher/Service:RetriveAsync() Leaving")

	if svc.serviceDone {
		return errors.New("Host Fetcher has been shut down - cannot accept any more requests")
	}
	fr := &fetchRequest{ctx, host, rcvrs}
	// queue up the request
	svc.rqstChan <- fr
	return nil
}

func (svc *Service) FetchDataAndRespond(hId uuid.UUID, connUrl string) {
	defaultLog.Trace("hostfetcher/Service:FetchDataAndRespond() Entering")
	defer defaultLog.Trace("hostfetcher/Service:FetchDataAndRespond() Leaving")

	//TODO: update the state in the context to reflect that we are about to start processing

	hostData, err := svc.GetHostData(connUrl)
	if err != nil {
		defaultLog.WithError(err).Errorf("hostfetcher/Service:FetchDataAndRespond() Failed to get data	")
		//TODO - presume that error is due to connection failure and we need to retry operation
		svc.retryRqstChan <- retryRequest{
			retryTime: time.Now().Add(time.Duration(svc.retryIntervalMins) * time.Minute),
			hostId:    hId,
		}
		_ = svc.hs.Persist(&hvs.HostStatus{
			HostID: hId,
			HostStatusInformation: hvs.HostStatusInformation{
				HostState:         hvs.HostStateUnknown,
				LastTimeConnected: time.Now(),
			},
		})
		return
	}
	//TODO: we need to check if the error is due to a connection failure.. In this case, we need to
	// add it to another channel to be resubmitted.
	log.Debug(" data for ", hId, "using connection string", connUrl)
	// work is done. get the list of callbacks and delete the entry.
	svc.wmLock.Lock()
	frs := svc.workMap[hId]
	delete(svc.workMap, hId)
	svc.wmLock.Unlock()
	_ = svc.hs.Persist(&hvs.HostStatus{
		HostID: hId,
		HostStatusInformation: hvs.HostStatusInformation{
			HostState:         hvs.HostStateConnected,
			LastTimeConnected: time.Now(),
		},
		HostManifest: *hostData,
	})

	for _, fr := range frs {
		select {
		case <-fr.ctx.Done():
			continue
		default:
		}
		for _, rcv := range fr.rcvrs {
			_ = rcv.ProcessHostData(fr.ctx, fr.host, hostData, err)
		}

	}

}

func (svc *Service) GetHostData(connUrl string) (*types.HostManifest, error) {
	defaultLog.Trace("hostfetcher/Service:GetHostData() Entering")
	defer defaultLog.Trace("hostfetcher/Service:GetHostData() Leaving")

	//get the host data
	connectionString, _, err := controllers.GenerateConnectionString(connUrl, svc.hcCfg.ServiceUsername,
		svc.hcCfg.ServicePassword,
		svc.hcCfg.HCStore)
	if err != nil {
		defaultLog.WithError(err).Error("hostfetcher/Service:GetHostData() Could not generate formatted connection string")
		return nil, err
	}

	connector, err := svc.hcf.NewHostConnector(connectionString)
	if err != nil {
		return nil, err
	}
	data, err := connector.GetHostManifest()
	return &data, err
}
