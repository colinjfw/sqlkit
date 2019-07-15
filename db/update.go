// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"github.com/colinjfw/sqlkit/encoding"
)

// Update returns a new Update statement.
func Update(table string) UpdateStmt { return UpdateStmt{table: table} }

// UpdateStmt represents an UPDATE in SQL.
type UpdateStmt struct {
	dialect Dialect
	table   string
	columns []string
	values  []interface{}
	args    []interface{}
	where   string
	err     error
	encoder encoding.Encoder
}

// Table configures the table for the query.
func (i UpdateStmt) Table(table string) UpdateStmt {
	i.table = table
	return i
}

// Columns sets the colums for the update.
func (i UpdateStmt) Columns(cols ...string) UpdateStmt {
	i.columns = cols
	return i
}

// Where configures the WHERE block.
func (i UpdateStmt) Where(where string, args ...interface{}) UpdateStmt {
	i.where = where
	i.args = append(i.args, args...)
	return i
}

// Values sets the values for the update.
func (i UpdateStmt) Values(vals ...interface{}) UpdateStmt {
	i.values = vals
	return i
}

// Value configures a single value for the query.
func (i UpdateStmt) Value(name string, val interface{}) UpdateStmt {
	i.values = append(i.values, val)
	i.columns = append(i.columns, name)
	return i
}

// Record will encode the struct and append the columns and returned values.
func (i UpdateStmt) Record(obj interface{}, fields ...string) UpdateStmt {
	cols, vals, err := i.encoder.Encode(obj, fields...)
	if err != nil {
		i.err = err
		return i
	}
	i.columns = append(i.columns, cols...)
	i.values = append(i.values, vals...)
	return i
}

// SQL implements the SQL interface.
func (i UpdateStmt) SQL() (string, []interface{}, error) {
	if len(i.columns) != len(i.values) {
		return "", nil, ErrStatementInvalid
	}
	sql := dialects[i.dialect].update(i)
	return sql, append(i.values, i.args...), i.err
}
