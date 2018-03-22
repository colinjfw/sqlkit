// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package encoding

import (
	"strings"
	"testing"

	"github.com/jmoiron/sqlx/reflectx"
	"github.com/stretchr/testify/require"
)

func TestMarshal_Obj(t *testing.T) {
	type obj struct {
		ID int `db:"dbid"`
	}
	names, args, err := Marshal(obj{})
	require.Nil(t, err)
	require.Equal(t, []string{"dbid"}, names)
	require.Equal(t, []interface{}{0}, args)
}

func TestMarshal_NestedObj(t *testing.T) {
	type Parent struct {
		ID int `db:"id"`
	}
	type obj struct {
		Parent
		Other int
	}
	names, args, err := Marshal(obj{})
	require.Nil(t, err)
	require.Equal(t, []string{"other", "id"}, names)
	require.Equal(t, []interface{}{0, 0}, args)
}

func TestMarshal_NestedDoubleObj(t *testing.T) {
	type Parent struct {
		ID int `db:"id"`
	}
	type Parent2 struct {
		ID int `db:"id"`
	}
	type obj struct {
		Parent
		Parent2
		Other int
	}
	_, _, err := Marshal(obj{})
	require.Equal(t, ErrDuplicateValue, err)
}

func TestMarshal_NestedDoubleUnsafe(t *testing.T) {
	type Parent struct {
		ID int `db:"id"`
	}
	type Parent2 struct {
		ID int `db:"id"`
	}
	type obj struct {
		Parent
		Parent2
		Other int
	}
	cols, vals, err := NewEncoder().Unsafe().Encode(obj{})
	require.Nil(t, err)
	require.Equal(t, []string{"other", "id"}, cols)
	require.Equal(t, []interface{}{0, 0}, vals)
}

func TestMarshal_CustomMapper(t *testing.T) {
	type Parent struct {
		ID int `db:"id"`
	}
	type obj struct {
		Parent
		Other int
	}
	m := NewEncoder().WithMapper(reflectx.NewMapperFunc("db", strings.ToLower))
	names, args, err := m.Encode(obj{})
	require.Nil(t, err)
	require.Equal(t, []string{"other", "id"}, names)
	require.Equal(t, []interface{}{0, 0}, args)
}
