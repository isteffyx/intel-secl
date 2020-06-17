/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package controllers_test

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/constants"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/controllers"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/controllers/mocks"
	hvsRoutes "github.com/intel-secl/intel-secl/v3/pkg/hvs/router"
	"github.com/intel-secl/intel-secl/v3/pkg/model/hvs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("HostStatusController", func() {
	var router *mux.Router
	var w *httptest.ResponseRecorder
	var hostStatusStore *mocks.MockHostStatusStore
	var hostStatusController *controllers.HostStatusController
	BeforeEach(func() {
		router = mux.NewRouter()
		hostStatusStore = mocks.NewFakeHostStatusStore()
		hostStatusController = &controllers.HostStatusController{Store: hostStatusStore}
	})

	// Specs for HTTP Get to "/host-status"
	Describe("Search HostStatus", func() {
		Context("When no filter arguments are passed", func() {
			It("All HostStatus records are returned", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(4))
			})
		})

		Context("When filtered by HostStatus id", func() {
			It("Should get a single HostStatus entry", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?id=afed7372-18c3-42af-bd9a-70b7f44c11ad", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(1))
			})
		})

		Context("When filtered by a non-existent hostId", func() {
			It("Should get an empty list", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostId=13885605-a0ee-41f2-b6fc-fd82edc487ad", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(hsCollection.HostStatuses).To(BeNil())
			})
		})

		Context("When filtered by an invalid hostId", func() {
			It("Should get an 400 error", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostId=13885605-a0ee-41f20000000000000000000000-b6fc-fd82edc487ad", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})

		Context("When filtered by hostId", func() {
			It("Should get a filtered list of HostStatuses by host-id", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostId=47a3b602-f321-4e03-b3b2-8f3ca3cde128", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(2))
			})
		})

		Context("When filtered by host-hardware-id", func() {
			It("Should get a filtered list of HostStatuses by host-hardware-id", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostHardwareId=1ad9c003-b0e0-4319-b2b3-06053dfd1407", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(2))
			})
		})

		Context("When filtered by an invalid host-hardware-id", func() {
			It("Should get an error response", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostHardwareId=1ad9c003-ABCABCABC-4319-b2b3-06053dfd1407", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})

		Context("When filtered by an non-existent host-hardware-id", func() {
			It("Should get an empty response", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostHardwareId=7f71bff0-3c12-4f92-9a77-d380eb9ad2e2", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).ToNot(HaveOccurred())
				Expect(hsCollection.HostStatuses).To(BeNil())
			})
		})

		Context("When searching for a list of CONNECTED hosts", func() {
			It("Should get a filtered list of HostStatuses with HostState = CONNECTED", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostStatus=connected", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(2))
			})
		})

		Context("When searching for a list of UNKNOWN hosts", func() {
			It("Should get a filtered list of HostStatuses with HostState = UNKNOWN", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostStatus=unknown", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(1))
			})
		})

		Context("When filtering HostStatus by an invalid HostState value", func() {
			It("Should see an 400 error", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?hostStatus=BADSTATUS", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})

		Context("When HostStatus filtered by numberOfDays old", func() {
			It("Should get a filtered list of HostStatuses by numberOfDays", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?numberOfDays=2", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(3))
			})
		})

		Context("When HostStatus filtered by an invalid value for numberOfDays", func() {
			It("Should get a 400 error", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?numberOfDays=-2", nil)
				Expect(err).ToNot(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})

		Context("When limiting the number of rows returned from HostStatus search", func() {
			It("Should get a truncated list of HostStatuses", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?limit=2", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(hsCollection.HostStatuses)).To(Equal(2))
			})
		})

		Context("When limiting the number of rows returned from HostStatus search with an invalid value", func() {
			It("Should get a 400 error", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?limit=-2", nil)
				Expect(err).ToNot(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})

		Context("Search HostStatuses from data store with invalid id", func() {
			It("Should return 400 error", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?id=e57e5ea0-d465-461e-882d-", nil)
				Expect(err).ToNot(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})

		Context("Search HostStatuses from data store with valid fromDate and toDate", func() {
			It("Should return a list of HostStatus", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?fromDate="+time.Now().Add(-mocks.TimeDuration12Hrs).Format(constants.HostStatusDateFormat)+"&toDate="+time.Now().Format(constants.HostStatusDateFormat), nil)
				Expect(err).ToNot(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).ToNot(HaveOccurred())
				Expect(hsCollection).ToNot(BeNil())
			})
		})

		Context("Search HostStatuses from data store with invalid fromDate and toDate", func() {
			It("Should return a list of HostStatus", func() {
				router.Handle("/host-status", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status?fromDate="+time.Now().Add(-mocks.TimeDuration12Hrs).Format(constants.HostStatusDateFormat)+"ABC"+"&toDate="+time.Now().Format(constants.HostStatusDateFormat), nil)
				Expect(err).ToNot(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})
	})

	// Specs for HTTP Get to "/host-status/{hoststatus_id}"
	Describe("Retrieve HostStatus", func() {
		Context("Retrieve HostStatus by valid ID from data store", func() {
			It("Should retrieve HostStatus", func() {
				router.Handle("/host-status/{id}", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status/afed7372-18c3-42af-bd9a-70b7f44c11ad", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
		Context("Try to retrieve HostStatus by non-existent ID from data store", func() {
			It("Should fail to retrieve HostStatus", func() {
				router.Handle("/host-status/{id}", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status/73755fda-c910-46be-821f-e8ddeab189e9", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))

				var hsCollection *hvs.HostStatusCollection
				err = json.Unmarshal(w.Body.Bytes(), &hsCollection)
				Expect(err).To(HaveOccurred())
				Expect(hsCollection).To(BeNil())
			})
		})
		Context("Try to retrieve HostStatus by invalid ID from data store", func() {
			It("Should fail to retrieve HostStatus", func() {
				router.Handle("/host-status/{id}", hvsRoutes.ErrorHandler(hvsRoutes.ResponseHandler(hostStatusController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/host-status/ee37c360-7eae-4250-a677", nil)
				Expect(err).ToNot(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})