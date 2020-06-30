/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package hvs

import (
	"github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/model"
)

// Flavor sourced from the lib/flavor - this is a external request/response on the HVS API
type Flavor = model.Flavor

// FlavorCollection is a list of Flavor objects in response to a Flavor Search query
type FlavorCollection struct {
	Flavors []Flavor `json:"flavors" xml:"flavors"`
}

// SignedFlavor sourced from the lib/flavor - this is a external request/response on the HVS API
type SignedFlavor = model.SignedFlavor

// SignedFlavorCollection is a list of SignedFlavor objects
type SignedFlavorCollection struct {
	SignedFlavors []*SignedFlavor `json:"signed_flavors" xml:"signed_flavors"`
}
