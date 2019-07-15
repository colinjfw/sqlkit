// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import "sort"

type operator int

const (
	none operator = iota
	and
	or
	in
	eq
	ne
	is
	gt
	lt
	lte
	gte
)

func (o operator) String() string {
	switch o {
	case and:
		return "AND"
	case or:
		return "OR"
	case in:
		return "IN"
	case eq:
		return "="
	case gt:
		return ">"
	case lt:
		return "<"
	case lte:
		return "<="
	case gte:
		return ">="
	case is:
		return "IS"
	case ne:
		return "!="
	default:
		return "<unknown>"
	}
}

// Statement represents an SQL statement.
type Statement struct {
	left     SQL
	operator operator
	right    SQL
}

func (s Statement) isZero() bool { return s.operator == none }

// And configures the And operator.
func (s Statement) And(c Statement) Statement {
	return Statement{left: s, operator: and, right: c}
}

// Or configures an or operator.
func (s Statement) Or(c Statement) Statement {
	return Statement{left: s, operator: or, right: c}
}

// SQL implements the SQL interface.
func (s Statement) SQL() (string, []interface{}, error) {
	left, argsLeft, err := s.left.SQL()
	if err != nil {
		return "", nil, err
	}
	right, argsRight, err := s.right.SQL()
	if err != nil {
		return "", nil, err
	}
	if left == "" || right == "" {
		return "", nil, ErrStatementInvalid
	}
	sql := "(" + left + " " + s.operator.String() + " " + right + ")"
	return sql, append(argsLeft, argsRight...), nil
}

// NotEq sets col != val.
func NotEq(col string, val interface{}) Statement { return stmt(ne, col, val) }

// Eq sets col = val.
func Eq(col string, val interface{}) Statement { return stmt(eq, col, val) }

// Gt sets col > val.
func Gt(col string, val interface{}) Statement { return stmt(gt, col, val) }

// GtEq sets col >= val.
func GtEq(col string, val interface{}) Statement { return stmt(gte, col, val) }

// Lt sets col < val.
func Lt(col string, val interface{}) Statement { return stmt(lt, col, val) }

// LtEq sets col <= val.
func LtEq(col string, val interface{}) Statement { return stmt(lte, col, val) }

// In sets col in (val).
func In(col string, val interface{}) Statement { return stmt(in, col, val) }

// Is sets col is val.
func Is(col string, val interface{}) Statement { return stmt(is, col, val) }

// EqAllMap sets (key = val) for every parameter joining with AND.
func EqAllMap(m map[string]interface{}) (s Statement) {
	for _, k := range mapKeys(m) {
		v := m[k]
		if s.operator == none {
			s = Eq(k, v)
		} else {
			s = s.And(Eq(k, v))
		}
	}
	return
}

// EqAnyMap sets (key = val) for every parameter joining with OR.
func EqAnyMap(m map[string]interface{}) (s Statement) {
	for _, k := range mapKeys(m) {
		v := m[k]
		if s.operator == none {
			s = Eq(k, v)
		} else {
			s = s.Or(Eq(k, v))
		}
	}
	return
}

// Null is a shorthand raw for SQL.
var Null = Raw("NULL")

func stmt(op operator, col string, value interface{}) Statement {
	var right SQL
	if s, ok := value.(SQL); ok {
		right = s
	} else {
		right = RawWithValues("?", value)
	}
	if s, ok := right.(SelectStmt); ok {
		right = parens{s}
	}
	return Statement{
		left:     Raw(col),
		operator: op,
		right:    right,
	}
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

type parens struct { sql SQL }

func (q parens) SQL() (string, []interface{}, error) {
	sql, values, err := q.sql.SQL()
	return "(" + sql + ")", values, err
}
