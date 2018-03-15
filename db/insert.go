package sql

import (
	"github.com/coldog/sqlkit/marshal"
)

func Insert() InsertStmt { return InsertStmt{} }

type InsertStmt struct {
	dialect Dialect
	table   string
	columns []string
	values  []interface{}
	err     error
}

func (i InsertStmt) Into(table string) InsertStmt {
	i.table = table
	return i
}

func (i InsertStmt) Columns(cols ...string) InsertStmt {
	i.columns = cols
	return i
}

func (i InsertStmt) Values(vals ...interface{}) InsertStmt {
	i.values = vals
	return i
}

func (i InsertStmt) Value(name string, val interface{}) InsertStmt {
	i.values = append(i.values, val)
	i.columns = append(i.columns, name)
	return i
}

func (i InsertStmt) Record(obj interface{}) InsertStmt {
	cols, vals := marshal.Marshal(obj)
	i.columns = append(i.columns, cols...)
	i.values = append(i.values, vals...)
	return i
}

func (i InsertStmt) SQL() (string, []interface{}, error) {
	if len(i.columns) != len(i.values) {
		return "", nil, ErrStatementInvalid
	}
	sql := dialects[i.dialect].insert(i)
	return sql, i.values, i.err
}
