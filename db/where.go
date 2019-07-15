// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import "reflect"

type where struct {
	bare SQL       // Bare where clause.
	stmt Statement // Statement for multiple calls of where.
}

func (q where) where(op operator, where interface{}, values ...interface{}) where {
	var sql SQL
	switch v := where.(type) {
	case string:
		if len(values) > 0 {
			sql = sqlHolder{sql: v, args: values}
		} else {
			sql = Raw(v)
		}
	case SQL:
		sql = v
	default:
		panic("unknown type")
	}

	if q.stmt.isZero() {
		// No statement yet to add to.
		if q.bare == nil {
			// Empty state just add the sql.
			q.bare = sql
		} else {
			// Bare is present, add a clause to the right.
			q.stmt = Statement{left: q.bare, operator: op, right: sql}
			q.bare = nil
		}
	} else {
		// We have a statement add another clause to the right.
		next := Statement{left: q.stmt, operator: op, right: sql}
		q.stmt = next
	}
	return q
}

func (q where) SQL() (sql string, values []interface{}, err error) {
	if q.bare != nil {
		sql, values, err = q.bare.SQL()
	} else if !q.stmt.isZero() {
		sql, values, err = q.stmt.SQL()
	} else {
		return "", nil, nil // No statements.
	}
	for i, arg := range values {
		if ok, l, inVals := isSlice(arg); ok {
			sql, err = insertQuestions(sql, i, l)
			inVals = append(inVals, values[i+1:]...)
			values = append(values[:i], inVals...)
		}
	}
	return
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
