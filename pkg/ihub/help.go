/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package ihub

import (
	"fmt"

	"github.com/intel-secl/intel-secl/v3/pkg/ihub/version"
)

const helpStr = `Usage:
	ihub <command> [arguments]
	
Available Commands:
	-h|--help                    Show this help message
	-v|--version                 Show the version of current ihub build
	setup <task>                 Run setup task
	start [-i <instance-name>]   Start ihub
		-i|--instance <instance-name>      the instance name to execute command against specific instance
	status [-i <instance-name>]  Show the status of ihub
		-i|--instance <instance-name>      the instance name to execute command against specific instance
	stop [-i <instance-name>]    Stop ihub
		-i|--instance <instance-name>      the instance name to execute command against specific instance
	uninstall [--purge] [--exec] [-i <instance-name>] Uninstall ihub
		-i|--instance <instance-name>      the instance name to execute command against specific instance
		--purge                            all configuration and data files will be removed for instance if this flag is set
		--exec                             executable will be removed which is common for all instances if this flag is set

Usage of ihub setup:
	ihub setup <task> [--help] [--force] [-f <answer-file>] [-i <instance-name>]
		-i|--instance <instance-name>      the instance name to execute command against specific instance
		--help                             show help message for setup task
		--force                            existing configuration will e overwritten if this flag is set
		-f|--file <answer-file>            the answer file with required arguments

Available Tasks for setup:
	all                                 Runs all setup tasks
	download-ca-cert                    Download CMS root CA certificate
	download-cert-tls                   Download CA certificate from CMS for tls
	attestation-service-connection      Establish Attestation service connection
	tenant-service-connection           Establish Tenant service connection
	create-signing-key                  Create signing key for IHUB
	download-saml-cert                  Download SAML certificate from Attestation service
	update-service-config               Sets or Updates the Service configuration                
`

func (app *App) printUsage() {
	fmt.Fprintln(app.consoleWriter(), helpStr)
}

func (app *App) printUsageWithError(err error) {
	fmt.Fprintln(app.errorWriter(), "Application returned with error:", err.Error())
	fmt.Fprintln(app.errorWriter(), helpStr)
}

func (app *App) printVersion() {
	fmt.Fprintf(app.consoleWriter(), version.GetVersion())
}
