package marshal

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
)

var (
	// ErrRequiresPtr is returned if a non pointer value is passed to Decode.
	ErrRequiresPtr = errors.New("sqlkit/marshal: non pointer passed to Decode")
	// ErrMustNotBeNil is returned if a nil value is passed to Decode.
	ErrMustNotBeNil = errors.New("sqlkit/marshal: nil value passed to Decode")
	// ErrMissingDestination is returned if a destination value is missing and
	// unsafe is not configured.
	ErrMissingDestination = errors.New("sqlkit/marshal: missing destination")
	// ErrTooManyColumns is returned when too many columns are present to scan
	// into a value.
	ErrTooManyColumns = errors.New("sqlkit/marshal: too many columns to scan")
	// ErrNoRows is mirrored from the database/sql package.
	ErrNoRows = sql.ErrNoRows
)

// DefaultMapper is the default reflectx mapper used. This uses strings.ToLower
// to map field names.
var DefaultMapper = reflectx.NewMapperFunc("db", strings.ToLower)

var _scannerInterface = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

// isScannable takes the reflect.Type and the actual dest value and returns
// whether or not it's Scannable.  Something is scannable if:
//   * it is not a struct
//   * it implements sql.Scanner
//   * it has no exported fields
func isScannable(t reflect.Type) bool {
	if reflect.PtrTo(t).Implements(_scannerInterface) {
		return true
	}
	if t.Kind() != reflect.Struct {
		return true
	}
	if len(DefaultMapper.TypeMap(t).Index) == 0 {
		return true
	}
	return false
}

// nilSafety will not scan a value into the target if the src value to be
// scanned is nil. This means, the default zero value in the object will be left
// alone.
type nilSafety struct {
	dest interface{}
}

// Scan implements the nilSafety behaviour by returning early on nil and then
// using the original behaviour from database/sql.
func (n *nilSafety) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return convertAssign(n.dest, src)
}

// fieldsByTraversal fills a list of value interfaces. It will also return an
// error if unsafe is not set and a value is missing in the destination struct.
func fieldsByTraversal(
	v reflect.Value, traversals [][]int, values []interface{}, unsafe bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			if unsafe {
				values[i] = new(interface{})
				continue
			}
			return ErrMissingDestination
		}
		f := v
		for _, i2 := range traversal {
			f = reflect.Indirect(f).Field(i2)
		}
		values[i] = &nilSafety{dest: f.Addr().Interface()}
	}
	return nil
}

// Unmarshal will run Decode with the default Decoder configuration.
func Unmarshal(dest interface{}, rows *sql.Rows) error {
	return Decoder{}.Decode(dest, rows)
}

// NewDecoder returns a new Decoder with the default configuration.
func NewDecoder() Decoder { return Decoder{} }

// Decoder holds settings for decoding from *sql.Rows into struct objects.
type Decoder struct {
	unsafe bool
	mapper *reflectx.Mapper
}

// Unsafe configures the Decoder to be unsafe, meaning that the Decoder can lose
// information by skipping over fields that do not exist in the destination
// struct.
func (e Decoder) Unsafe() Decoder {
	e.unsafe = true
	return e
}

// WithMapper configures the decoder with a reflectx.Mapper for configuring
// different fields to be decoded. The DefaultMapper is used if this is not set.
func (e Decoder) WithMapper(m *reflectx.Mapper) Decoder {
	e.mapper = m
	return e
}

// Decode does the work of decoding an *sql.Rows into a struct, array or scalar
// value. Depending on the value passed in the decoder will perform the
// following actions:
//
// * If an array is passed in the decoder will work through all rows
//   initializing and scanning in a new instance of the value(s) for all rows.
// * If the values inside the array are scalar, then the decoder will check for
//   only a single column to scan and scan this in.
// * If a single struct or scalar is passed in the decoder will loop once over
//   the rows returning sql.ErrNoRows if this is improssible and scan the value.
//
// The rows object is not closed after iteration is completed. The Decode
// function is thread safe.
func (e Decoder) Decode(dest interface{}, rows *sql.Rows) error {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return ErrRequiresPtr
	}
	if value.IsNil() {
		return ErrRequiresPtr
	}
	var slice bool
	base := reflectx.Deref(value.Type())
	if base.Kind() == reflect.Slice {
		slice = true
	}

	if slice {
		if err := e.scanAll(base, value, rows); err != nil {
			return err
		}
	} else {
		if err := e.scanRow(base, value, dest, rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (e Decoder) scanAll(slice reflect.Type, value reflect.Value, rows *sql.Rows) error {
	direct := reflect.Indirect(value)
	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := reflectx.Deref(slice.Elem())
	scannable := isScannable(base)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	m := DefaultMapper
	if e.mapper != nil {
		m = e.mapper
	}

	// if it's a base type make sure it only has 1 column;  if not return an error
	if scannable && len(columns) > 1 {
		return ErrTooManyColumns
	}

	if !scannable {
		var values []interface{}
		fields := m.TraversalsByName(base, columns)
		values = make([]interface{}, len(columns))

		for rows.Next() {
			// create a new struct type (which returns PtrTo) and indirect it
			vp := reflect.New(base)
			v := reflect.Indirect(vp)

			err = fieldsByTraversal(v, fields, values, e.unsafe)
			if err != nil {
				return err
			}

			// scan into the struct field pointers and append to our results
			err = rows.Scan(values...)
			if err != nil {
				return err
			}

			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, v))
			}
		}
	} else {
		for rows.Next() {
			vp := reflect.New(base)
			err = rows.Scan(vp.Interface())
			if err != nil {
				return err
			}
			// append
			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
			}
		}
	}
	return nil
}

func (e Decoder) scanRow(base reflect.Type, value reflect.Value, dest interface{}, rows *sql.Rows) error {
	m := DefaultMapper
	if e.mapper != nil {
		m = e.mapper
	}

	// Do this early so we don't have to waste type reflecting or traversing if
	// there isn't anything to scan.
	if !rows.Next() {
		return ErrNoRows
	}

	value = reflect.Indirect(value)
	scannable := isScannable(base)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// If it's scannable, like a scalar type, we can scan directly right here.
	if scannable {
		// Ensure that only one column to scan into.
		if len(columns) > 1 {
			return ErrTooManyColumns
		}
		return rows.Scan(dest)
	}

	fields := m.TraversalsByName(value.Type(), columns)
	values := make([]interface{}, len(columns))
	err = fieldsByTraversal(value, fields, values, e.unsafe)
	if err != nil {
		return err
	}

	// scan into the struct field pointers and append to our results
	err = rows.Scan(values...)
	if err != nil {
		return err
	}
	return nil
}
