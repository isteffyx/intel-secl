/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package model

import "github.com/google/uuid"

/**
 *
 * @author purvades
 */

// Meta holds metadata information related to the Flavor
type Meta struct {
	Schema      *Schema     `json:"schema,omitempty"`
	ID          uuid.UUID   `json:"id"`
	Realm       string      `json:"realm,omitempty"`
	Description Description `json:"description,omitempty"`
	Vendor      string      `json:"vendor,omitempty"`
}

// Schema defines the Uri of the schema
type Schema struct {
	Uri string `json:"uri,omitempty"`
}
