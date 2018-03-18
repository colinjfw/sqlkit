package marshal

import (
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

var (
	// ErrDuplicateValue is returned when duplicate values exist in the struct.
	ErrDuplicateValue = errors.New("sqlkit/marshal: duplicate values")
)

// NewEncoder returns an Encoder with the default settings which are blank.
func NewEncoder() Encoder { return Encoder{mapper: DefaultMapper} }

// Encoder manages options for encoding.
type Encoder struct {
	unsafe bool
	mapper *reflectx.Mapper
}

// Unsafe configures and returns a new Encoder which uses unsafe options.
// Specifically it will ignore duplicate fields and just take the first one.
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

// Encode will encode to a set of fields and values using the Encoder's
// settings. It will return and error if there are duplicate fields and unsafe
// is not set.
func (e Encoder) Encode(obj interface{}) ([]string, []interface{}, error) {
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
func Marshal(obj interface{}) ([]string, []interface{}, error) {
	return Encoder{}.Encode(obj)
}

func inStr(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}
