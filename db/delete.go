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

// Join adds a join.
func (q DeleteStmt) Join(table, on string) DeleteStmt {
	q.sel = q.sel.Join(table, on)
	return q
}

// InnerJoin adds a join of type INNER.
func (q DeleteStmt) InnerJoin(table, on string) DeleteStmt {
	q.sel = q.sel.InnerJoin(table, on)
	return q
}

// LeftJoin adds a join of type LEFT.
func (q DeleteStmt) LeftJoin(table, on string) DeleteStmt {
	q.sel = q.sel.LeftJoin(table, on)
	return q
}

// RightJoin adds a join of type RIGHT.
func (q DeleteStmt) RightJoin(table, on string) DeleteStmt {
	q.sel = q.sel.RightJoin(table, on)
	return q
}

// SQL implements the SQL interface.
func (q DeleteStmt) SQL() (string, []interface{}, error) {
	sql := dialects[q.dialect].delete(q)
	return sql, q.sel.values, q.sel.err
}
