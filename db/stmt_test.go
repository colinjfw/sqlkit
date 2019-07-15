// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEq(t *testing.T) {
	sql, args, err := Eq("a", "b").SQL()
	require.NoError(t, err)
	require.Equal(t, "(a = ?)", sql)
	require.Equal(t, []interface{}{"b"}, args)
}

func TestGt(t *testing.T) {
	sql, args, err := Gt("a", "b").SQL()
	require.NoError(t, err)
	require.Equal(t, "(a > ?)", sql)
	require.Equal(t, []interface{}{"b"}, args)
}

func TestGtEq(t *testing.T) {
	sql, args, err := GtEq("a", "b").SQL()
	require.NoError(t, err)
	require.Equal(t, "(a >= ?)", sql)
	require.Equal(t, []interface{}{"b"}, args)
}

func TestLt(t *testing.T) {
	sql, args, err := Lt("a", "b").SQL()
	require.NoError(t, err)
	require.Equal(t, "(a < ?)", sql)
	require.Equal(t, []interface{}{"b"}, args)
}

func TestLtEq(t *testing.T) {
	sql, args, err := LtEq("a", "b").SQL()
	require.NoError(t, err)
	require.Equal(t, "(a <= ?)", sql)
	require.Equal(t, []interface{}{"b"}, args)
}

func TestNotEq(t *testing.T) {
	sql, args, err := NotEq("a", "b").SQL()
	require.NoError(t, err)
	require.Equal(t, "(a != ?)", sql)
	require.Equal(t, []interface{}{"b"}, args)
}

func TestIn(t *testing.T) {
	sql, args, err := In("col", []int{1, 2}).SQL()
	require.NoError(t, err)
	require.Equal(t, "(col IN ?)", sql)
	require.Equal(t, []interface{}{[]int{1, 2}}, args)
}

func TestIs(t *testing.T) {
	sql, args, err := Is("col", Null).SQL()
	require.NoError(t, err)
	require.Equal(t, "(col IS NULL)", sql)
	require.Equal(t, []interface{}(nil), args)
}

func TestAnd(t *testing.T) {
	sql, args, err := Eq("a", 1).And(NotEq("b", 2)).SQL()
	require.NoError(t, err)
	require.Equal(t, "((a = ?) AND (b != ?))", sql)
	require.Equal(t, []interface{}{1, 2}, args)
}

func TestOr(t *testing.T) {
	sql, args, err := Eq("a", 1).Or(NotEq("b", 2)).SQL()
	require.NoError(t, err)
	require.Equal(t, "((a = ?) OR (b != ?))", sql)
	require.Equal(t, []interface{}{1, 2}, args)
}

func TestEqAllMap(t *testing.T) {
	sql, args, err := EqAllMap(map[string]interface{}{
		"x": 1,
		"y": 2,
	}).SQL()
	require.NoError(t, err)
	require.Equal(t, "((x = ?) AND (y = ?))", sql)
	require.Equal(t, []interface{}{1, 2}, args)
}

func TestEqAnyMap(t *testing.T) {
	sql, args, err := EqAnyMap(map[string]interface{}{
		"x": 1,
		"y": 2,
	}).SQL()
	require.NoError(t, err)
	require.Equal(t, "((x = ?) OR (y = ?))", sql)
	require.Equal(t, []interface{}{1, 2}, args)
}
