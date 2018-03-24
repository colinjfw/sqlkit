// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"fmt"
	"reflect"
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

// isSlice checks if the value is a slice, if it is, it returns the length and a
// representation of the slice as an []interface{}.
func isSlice(i interface{}) (bool, int, []interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice {
		l := v.Len()
		arr := make([]interface{}, l)
		for i := 0; i < l; i++ {
			arr[i] = v.Index(i).Interface()
		}
		return true, l, arr
	}
	return false, 0, nil
}

// insertQuestions will insert count questions in str at the question mark idx.
// This turns 'x in ?' into 'x in (?, ?)'.
func insertQuestions(str string, insertAt, count int) (string, error) {
	var qIdx int
	var sIdx int
	for i, char := range str {
		if char == '?' {
			if qIdx == insertAt {
				sIdx = i
				break
			}
			qIdx++
		}
	}
	if sIdx == 0 {
		// Couldn't find a question mark at this index.
		return str, fmt.Errorf(
			"sqlkit/db: could not find matching '?' at index %d", insertAt)
	}
	return str[:sIdx] + questions(count) + str[sIdx+1:], nil
}

// Where configures the WHERE clause. It expects values to be interpolated using
// the question (?) mark parameter. For values that are slices, the question
// mark will be transformed in the where query. This means that IN queries can
// be writted without knowing the specific number of arguments needed in the
// array.
func (q SelectStmt) Where(where string, values ...interface{}) SelectStmt {
	var err error
	for i, arg := range values {
		if ok, l, inVals := isSlice(arg); ok {
			where, err = insertQuestions(where, i, l)
			inVals = append(inVals, values[i+1:]...)
			values = append(values[:i], inVals...)
		}
	}
	if err != nil {
		q.err = err
		return q
	}
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
