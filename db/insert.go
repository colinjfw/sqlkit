// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"github.com/colinjfw/sqlkit/encoding"
)

// Insert constructs an InsertStmt.
func Insert() InsertStmt { return InsertStmt{} }

// InsertStmt represents an INSERT in SQL.
type InsertStmt struct {
	dialect Dialect
	table   string
	columns []string
	rows    [][]interface{}
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

// Values configures a single row of values.
func (i InsertStmt) Values(vals ...interface{}) InsertStmt {
	i.rows = append(i.rows, vals)
	return i
}

// Value configures a single value insert.
func (i InsertStmt) Value(column string, value interface{}) InsertStmt {
	cols := []string{column}
	vals := []interface{}{value}
	return i.Row(cols, vals)
}

// Record will decode using the decoder into a list of fields and values.
func (i InsertStmt) Record(obj interface{}, fields ...string) InsertStmt {
	cols, vals, err := i.encoder.Encode(obj, fields...)
	if err != nil {
		i.err = err
		return i
	}
	return i.Row(cols, vals)
}

// Row configures a single row into the insert statement. If the columns don't
// match previous insert statements then an error is forwarded.
func (i InsertStmt) Row(cols []string, vals []interface{}) InsertStmt {
	// Only write if nil to allow multiple record calls. Only the first will
	// configure the columns.
	if i.columns == nil {
		i.columns = cols
	}
	// Validate that columns are the same for this record.
	if len(i.columns) != len(cols) {
		i.err = ErrStatementInvalid
		return i
	}
	for _, c := range cols {
		var exists bool
		for _, c2 := range i.columns {
			if c == c2 {
				exists = true
			}
		}
		if !exists {
			i.err = ErrStatementInvalid
			return i
		}
	}
	i.rows = append(i.rows, vals)
	return i
}

// SQL implements the SQL interface.
func (i InsertStmt) SQL() (string, []interface{}, error) {
	values := []interface{}{}
	for _, row := range i.rows {
		if len(i.columns) != len(row) {
			return "", nil, ErrStatementInvalid
		}
		values = append(values, row...)
	}
	sql := dialects[i.dialect].insert(i)
	return sql, values, i.err
}
