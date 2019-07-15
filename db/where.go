// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import "reflect"

type where struct {
	sql SQL // Bare where clause.
}

func (q where) where(op operator, where interface{}, values ...interface{}) where {
	var next SQL
	switch v := where.(type) {
	case string:
		if len(values) > 0 {
			next = RawWithValues(v, values...)
		} else {
			next = Raw(v) // More efficient implementation if no values.
		}
	case SQL:
		next = v
	default:
		panic("unknown type")
	}

	if q.sql == nil {
		q.sql =next
	} else {
		q.sql = Statement{left: q.sql, operator: op, right: next}
	}
	return q
}

func (q where) SQL() (sql string, values []interface{}, err error) {
	if q.sql != nil {
		sql, values, err = q.sql.SQL()
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
