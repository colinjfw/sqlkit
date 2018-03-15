package sql

import (
	"strconv"
)

func Select(cols ...string) SelectStmt { return SelectStmt{columns: cols} }

type SelectStmt struct {
	dialect   Dialect
	columns   []string
	table     string
	groupBy   []string
	orderBy   []string
	offset    string
	limit     string
	join      [][]string
	where     string
	returning []string
	values    []interface{}
	err       error
}

func (q SelectStmt) Select(cols ...string) SelectStmt {
	q.columns = cols
	return q
}

func (q SelectStmt) From(table string) SelectStmt {
	q.table = table
	return q
}

func (q SelectStmt) Where(where string, values ...interface{}) SelectStmt {
	// TODO: handle IN () queries.
	q.where = where
	q.values = append(q.values, values...)
	return q
}

func (q SelectStmt) GroupBy(groupBy ...string) SelectStmt {
	q.groupBy = groupBy
	return q
}

func (q SelectStmt) OrderBy(orderBy ...string) SelectStmt {
	q.orderBy = orderBy
	return q
}

func (q SelectStmt) Offset(offset int) SelectStmt {
	q.offset = strconv.FormatInt(int64(offset), 10)
	return q
}

func (q SelectStmt) Limit(limit int) SelectStmt {
	q.limit = strconv.FormatInt(int64(limit), 10)
	return q
}

func (q SelectStmt) Join(kind, table, on string, values ...interface{}) SelectStmt {
	q.join = append(q.join, []string{kind, table, on})
	q.values = append(q.values, values...)
	return q
}

func (q SelectStmt) InnerJoin(table, on string) SelectStmt {
	return q.Join("INNER", table, on)
}

func (q SelectStmt) LeftJoin(table, on string) SelectStmt {
	return q.Join("LEFT", table, on)
}

func (q SelectStmt) RightJoin(table, on string) SelectStmt {
	return q.Join("RIGHT", table, on)
}

func (q SelectStmt) SQL() (string, []interface{}, error) {
	sql := dialects[q.dialect].query(q)
	return sql, q.values, q.err
}
