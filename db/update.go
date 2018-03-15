package sql

import (
	"github.com/coldog/sqlkit/marshal"
)

func Update(table string) UpdateStmt {
	return UpdateStmt{table: table}
}

type UpdateStmt struct {
	dialect Dialect
	table   string
	columns []string
	values  []interface{}
	args    []interface{}
	where   string
	err     error
}

func (i UpdateStmt) Table(table string) UpdateStmt {
	i.table = table
	return i
}

func (i UpdateStmt) Columns(cols ...string) UpdateStmt {
	i.columns = cols
	return i
}

func (i UpdateStmt) Where(where string, args ...interface{}) UpdateStmt {
	i.where = where
	i.args = append(i.args, args...)
	return i
}

func (i UpdateStmt) Values(vals ...interface{}) UpdateStmt {
	i.values = vals
	return i
}

func (i UpdateStmt) Value(name string, val interface{}) UpdateStmt {
	i.values = append(i.values, val)
	i.columns = append(i.columns, name)
	return i
}

func (i UpdateStmt) Record(obj interface{}) UpdateStmt {
	cols, vals := marshal.Marshal(obj)
	i.columns = append(i.columns, cols...)
	i.values = append(i.values, vals...)
	return i
}

func (i UpdateStmt) SQL() (string, []interface{}, error) {
	if len(i.columns) != len(i.values) {
		return "", nil, ErrStatementInvalid
	}
	sql := dialects[i.dialect].update(i)
	return sql, append(i.values, i.args...), i.err
}
