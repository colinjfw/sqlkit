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
func (q DeleteStmt) Where(where string, values ...interface{}) DeleteStmt {
	q.sel = q.sel.Where(where, values...)
	return q
}

// Join adds a join statement of a specific kind.
func (q DeleteStmt) Join(kind, table, on string, values ...interface{}) DeleteStmt {
	q.sel = q.sel.Join(kind, table, on, values...)
	return q
}

// InnerJoin adds a join of type INNER.
func (q DeleteStmt) InnerJoin(table, on string) DeleteStmt {
	return q.Join("INNER", table, on)
}

// LeftJoin adds a join of type LEFT.
func (q DeleteStmt) LeftJoin(table, on string) DeleteStmt {
	return q.Join("LEFT", table, on)
}

// RightJoin adds a join of type RIGHT.
func (q DeleteStmt) RightJoin(table, on string) DeleteStmt {
	return q.Join("RIGHT", table, on)
}

// SQL implements the SQL interface.
func (q DeleteStmt) SQL() (string, []interface{}, error) {
	sql := dialects[q.dialect].delete(q)
	return sql, q.sel.values, q.sel.err
}
