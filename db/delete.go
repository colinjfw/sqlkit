// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

// Delete returns a new DeleteStmt.
func Delete() DeleteStmt { return DeleteStmt{} }

// DeleteStmt represents a DELETE in sql.
type DeleteStmt struct {
	dialect Dialect
	sel     SelectStmt
}

// From configures the table.
func (q DeleteStmt) From(table string) DeleteStmt {
	q.sel = q.sel.From(table)
	return q
}

// Where configures the WHERE clause in a DELETE statement. It follows the same
// format as the select statement where statement.
func (q DeleteStmt) Where(where interface{}, values ...interface{}) DeleteStmt {
	q.sel = q.sel.Where(where, values...)
	return q
}

// SQL implements the SQL interface.
func (q DeleteStmt) SQL() (string, []interface{}, error) {
	if q.sel.err != nil {
		return "", nil, q.sel.err
	}
	q.sel = q.sel.parseWhere()
	sql := dialects[q.dialect].delete(q)
	return sql, q.sel.values, q.sel.err
}
