// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"fmt"
	"strconv"
)

// Select returns a new SelectStmt.
func Select(cols ...string) SelectStmt { return SelectStmt{columns: cols} }

// SelectStmt represents a SELECT in sql.
type SelectStmt struct {
	dialect     Dialect
	columns     []string
	table       string
	groupBy     []string
	orderBy     []string
	offset      string
	limit       string
	join        [][]string
	whereClause where
	where       string
	values      []interface{}
	err         error
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
//
// The where parameter can take multiple types.
func (q SelectStmt) Where(where interface{}, values ...interface{}) SelectStmt {
	q.whereClause = q.whereClause.where(and, where, values...)
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

// join adds a join statement of a specific kind.
func (q SelectStmt) joins(kind, table, on string, values ...interface{}) SelectStmt {
	q.join = append(q.join, []string{kind, table, on})
	q.values = append(q.values, values...)
	return q
}

// Join adds a join.
func (q SelectStmt) Join(table, on string) SelectStmt {
	return q.joins("", table, on)
}

// InnerJoin adds a join of type INNER.
func (q SelectStmt) InnerJoin(table, on string) SelectStmt {
	return q.joins("INNER", table, on)
}

// LeftJoin adds a join of type LEFT.
func (q SelectStmt) LeftJoin(table, on string) SelectStmt {
	return q.joins("LEFT", table, on)
}

// RightJoin adds a join of type RIGHT.
func (q SelectStmt) RightJoin(table, on string) SelectStmt {
	return q.joins("RIGHT", table, on)
}

// SQL implements the SQL interface.
func (q SelectStmt) SQL() (string, []interface{}, error) {
	if q.err != nil {
		return "", nil, q.err
	}
	q = q.parseWhere()
	sql := dialects[q.dialect].query(q)
	return sql, q.values, q.err
}

func (q SelectStmt) parseWhere() SelectStmt {
	var values []interface{}
	q.where, values, q.err = q.whereClause.SQL()
	q.values = append(q.values, values...)
	return q
}
