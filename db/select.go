package db

import (
	"strconv"
)

// Select returns a new SelectStmt.
func Select(cols ...string) SelectStmt { return SelectStmt{columns: cols} }

// SelectStmt represents a SELECT in sql.
type SelectStmt struct {
	dialect Dialect
	columns []string
	table   string
	groupBy []string
	orderBy []string
	offset  string
	limit   string
	join    [][]string
	where   string
	values  []interface{}
	err     error
}

// Select configures the columns to select.
func (q SelectStmt) Select(cols ...string) SelectStmt {
	q.columns = cols
	return q
}

// From configures the table.
func (q SelectStmt) From(table string) SelectStmt {
	q.table = table
	return q
}

// Where configures the WHERE clause.
func (q SelectStmt) Where(where string, values ...interface{}) SelectStmt {
	// TODO: handle IN () queries.
	q.where = where
	q.values = append(q.values, values...)
	return q
}

// GroupBy configures the GROUP BY clause.
func (q SelectStmt) GroupBy(groupBy ...string) SelectStmt {
	q.groupBy = groupBy
	return q
}

// OrderBy configures the ORDER BY clause.
func (q SelectStmt) OrderBy(orderBy ...string) SelectStmt {
	q.orderBy = orderBy
	return q
}

// Offset configures the OFFSET clause.
func (q SelectStmt) Offset(offset int) SelectStmt {
	q.offset = strconv.FormatInt(int64(offset), 10)
	return q
}

// Limit configures the LIMIT clause.
func (q SelectStmt) Limit(limit int) SelectStmt {
	q.limit = strconv.FormatInt(int64(limit), 10)
	return q
}

// Join adds a join statement of a specific kind.
func (q SelectStmt) Join(kind, table, on string, values ...interface{}) SelectStmt {
	q.join = append(q.join, []string{kind, table, on})
	q.values = append(q.values, values...)
	return q
}

// InnerJoin adds a join of type INNER.
func (q SelectStmt) InnerJoin(table, on string) SelectStmt {
	return q.Join("INNER", table, on)
}

// LeftJoin adds a join of type LEFT.
func (q SelectStmt) LeftJoin(table, on string) SelectStmt {
	return q.Join("LEFT", table, on)
}

// RightJoin adds a join of type RIGHT.
func (q SelectStmt) RightJoin(table, on string) SelectStmt {
	return q.Join("RIGHT", table, on)
}

// SQL implements the SQL interface.
func (q SelectStmt) SQL() (string, []interface{}, error) {
	sql := dialects[q.dialect].query(q)
	return sql, q.values, q.err
}
