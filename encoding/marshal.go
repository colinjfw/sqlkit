// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package encoding

import (
	"errors"
	"reflect"
)

var (
	// ErrDuplicateValue is returned when duplicate values exist in the struct.
	ErrDuplicateValue = errors.New("sqlkit/marshal: duplicate values")
)

// Encode will encode to a set of fields and values using the Encoder's
// settings. It will return an error if there are duplicate fields and unsafe
// is not set.
//
// By default Encode walks through all fields in a struct. If the fields
// argument is provided however it will only walk through the fields provided.
// An error will be returned if the field doesn't exist.
func (e Encoder) Encode(obj interface{}, fields ...string) ([]string, []interface{}, error) {
	m := DefaultMapper
	if e.mapper != nil {
		m = e.mapper
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	tm := m.TypeMap(t)

	values := make([]interface{}, 0, len(tm.Index))
	names := make([]string, 0, len(tm.Index))
	for _, field := range tm.Index {
		if field.Embedded {
			continue
		}
		if len(fields) > 0 && !inStr(fields, field.Name) {
			continue
		}
		if inStr(names, field.Name) {
			if e.unsafe {
				continue
			}
			return nil, nil, ErrDuplicateValue
		}
		names = append(names, field.Name)
		f := v
		for _, i := range field.Index {
			f = f.Field(i)
		}
		values = append(values, f.Interface())
	}
	return names, values, nil
}

// Marshal runs the default encoder.
func Marshal(obj interface{}, fields ...string) ([]string, []interface{}, error) {
	return Encoder{}.Encode(obj, fields...)
}

func inStr(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}
