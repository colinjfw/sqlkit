// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package encoding

import "github.com/jmoiron/sqlx/reflectx"

// NewEncoder returns an Encoder with the default settings which are blank.
func NewEncoder() Encoder { return Encoder{} }

// Encoder manages options for encoding.
type Encoder struct {
	unsafe bool
	mapper *reflectx.Mapper
}

// Unsafe configures and returns a new Encoder which uses unsafe options.
// Specifically it will ignore duplicate fields on marshalling and will ignore
// missing fields on unmarshalling.
func (e Encoder) Unsafe() Encoder {
	e.unsafe = true
	return e
}

// WithMapper configures the encoder with a reflectx.Mapper for configuring
// different fields to be encoded. The DefaultMapper is used if this is not set.
func (e Encoder) WithMapper(m *reflectx.Mapper) Encoder {
	e.mapper = m
	return e
}
