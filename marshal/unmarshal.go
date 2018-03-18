package marshal

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
)

var (
	ErrRequiresPtr        = errors.New("sql/marshal: pointer passed to Unmarshal")
	ErrMustNotBeNil       = errors.New("sql/marshal: nil value passed to Unmarshal")
	ErrMissingDestination = errors.New("sql/marshal: missing destination")
	ErrTooManyColumns     = errors.New("sql/marshal: too many columns to scan")
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

	// it's not important that we use the right mapper for this particular object,
	// we're only concerned on how many exported fields this struct has
	if len(DefaultMapper.TypeMap(t).Index) == 0 {
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

func Unmarshal(dest interface{}, rows *sql.Rows) error {
	return Decoder{}.Decode(dest, rows)
}

func NewDecoder() Decoder { return Decoder{} }

type Decoder struct {
	unsafe bool
}

func (e Decoder) Unsafe() Decoder {
	e.unsafe = true
	return e
}

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

func (e Decoder) scanAll(
	slice reflect.Type, value reflect.Value, rows *sql.Rows) error {
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
		return ErrTooManyColumns
	}

	if !scannable {
		var values []interface{}
		m := DefaultMapper
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

func (e Decoder) scanRow(
	base reflect.Type, value reflect.Value, dest interface{},
	rows *sql.Rows) error {
	// Do this early so we don't have to waste type reflecting or traversing if
	// there isn't anything to scan.
	if !rows.Next() {
		return sql.ErrNoRows
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

	m := DefaultMapper
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
