package marshal

import (
	"testing"

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
