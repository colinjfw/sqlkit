package marshal

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx/reflectx"
)

var (
	ErrRequiresPtr        = errors.New("sql/marshal: pointer passed to Unmarshal")
	ErrMustNotBeNil       = errors.New("sql/marshal: nil value passed to Unmarshal")
	ErrMissingDestination = errors.New("sql/marshal: missing destination")
)

// NameMapper is used to map column names to struct field names.  By default,
// it uses strings.ToLower to lowercase struct field names.  It can be set
// to whatever you want, but it is encouraged to be set before sqlx is used
// as name-to-field mappings are cached after first use on a type.
var NameMapper = strings.ToLower
var origMapper = reflect.ValueOf(NameMapper)

// Rather than creating on init, this is created when necessary so that
// importers have time to customize the NameMapper.
var mpr *reflectx.Mapper

// mprMu protects mpr.
var mprMu sync.Mutex

// mapper returns a valid mapper using the configured NameMapper func.
func mapper() *reflectx.Mapper {
	mprMu.Lock()
	defer mprMu.Unlock()

	if mpr == nil {
		mpr = reflectx.NewMapperFunc("db", NameMapper)
	} else if origMapper != reflect.ValueOf(NameMapper) {
		// if NameMapper has changed, create a new mapper
		mpr = reflectx.NewMapperFunc("db", NameMapper)
		origMapper = reflect.ValueOf(NameMapper)
	}
	return mpr
}

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

	// it's not important that we use the right mapper for this particular object,
	// we're only concerned on how many exported fields this struct has
	m := mapper()
	if len(m.TypeMap(t).Index) == 0 {
		return true
	}
	return false
}

type nilSafety struct {
	dest interface{}
}

func (n *nilSafety) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return convertAssign(n.dest, src)
}

// fieldsByName fills a values interface with fields from the passed value based
// on the traversals in int.  If ptrs is true, return addresses instead of values.
// We write this instead of using FieldsByName to save allocations and map lookups
// when iterating over many rows.  Empty traversals will get an interface pointer.
// Because of the necessity of requesting ptrs or values, it's considered a bit too
// specialized for inclusion in reflectx itself.
func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		var f reflect.Value
		for _, i2 := range traversal {
			f = reflect.Indirect(v).Field(i2)
		}
		values[i] = &nilSafety{dest: f.Addr().Interface()}
	}
	return nil
}

func Unmarshal(dest interface{}, rows *sql.Rows) error {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return ErrRequiresPtr
	}
	if value.IsNil() {
		return ErrRequiresPtr
	}
	var slice bool
	t := reflectx.Deref(value.Type())
	if t.Kind() == reflect.Slice {
		slice = true
	}

	if slice {
		if err := scanAll(t, value, rows); err != nil {
			return err
		}
	} else {
		if err := scanRow(value, rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

func scanAll(slice reflect.Type, value reflect.Value, rows *sql.Rows) error {
	direct := reflect.Indirect(value)
	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := reflectx.Deref(slice.Elem())
	scannable := isScannable(base)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// if it's a base type make sure it only has 1 column;  if not return an error
	if scannable && len(columns) > 1 {
		return fmt.Errorf("non-struct dest type %s with >1 columns (%d)", base.Kind(), len(columns))
	}

	if !scannable {
		var values []interface{}
		m := mapper()
		fields := m.TraversalsByName(base, columns)
		values = make([]interface{}, len(columns))

		for rows.Next() {
			// create a new struct type (which returns PtrTo) and indirect it
			vp := reflect.New(base)
			v := reflect.Indirect(vp)

			err = fieldsByTraversal(v, fields, values)
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

func scanRow(value reflect.Value, rows *sql.Rows) error {
	value = reflect.Indirect(value)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	m := mapper()
	fields := m.TraversalsByName(value.Type(), columns)
	values := make([]interface{}, len(columns))
	err = fieldsByTraversal(value, fields, values)
	if err != nil {
		return err
	}

	// Check if we can iterate if not we return the standard sql ErrNoRows.
	if !rows.Next() {
		return sql.ErrNoRows
	}

	// scan into the struct field pointers and append to our results
	err = rows.Scan(values...)
	if err != nil {
		return err
	}
	return nil
}
