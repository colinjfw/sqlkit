package db

import (
	"github.com/coldog/sqlkit/encoding"
)

// Insert constructs an InsertStmt.
func Insert() InsertStmt { return InsertStmt{} }

// InsertStmt represents an INSERT in SQL.
type InsertStmt struct {
	dialect Dialect
	table   string
	columns []string
	values  []interface{}
	err     error
	encoder encoding.Encoder
}

// Into configures the table name.
func (i InsertStmt) Into(table string) InsertStmt {
	i.table = table
	return i
}

// Columns configures the columns.
func (i InsertStmt) Columns(cols ...string) InsertStmt {
	i.columns = cols
	return i
}

// Values configures the values.
func (i InsertStmt) Values(vals ...interface{}) InsertStmt {
	i.values = vals
	return i
}

// Value sets one column and one value.
func (i InsertStmt) Value(name string, val interface{}) InsertStmt {
	i.values = append(i.values, val)
	i.columns = append(i.columns, name)
	return i
}

// Record will decode using the decoder into a list of fields and values.
func (i InsertStmt) Record(obj interface{}) InsertStmt {
	cols, vals, err := i.encoder.Encode(obj)
	if err != nil {
		i.err = err
		return i
	}
	i.columns = append(i.columns, cols...)
	i.values = append(i.values, vals...)
	return i
}

// SQL implements the SQL interface.
func (i InsertStmt) SQL() (string, []interface{}, error) {
	if len(i.columns) != len(i.values) {
		return "", nil, ErrStatementInvalid
	}
	sql := dialects[i.dialect].insert(i)
	return sql, i.values, i.err
}
