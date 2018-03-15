package sql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testSQL(t *testing.T, expected string, values []interface{}, sql SQL) {
	str, vals, err := sql.SQL()
	require.Nil(t, err)
	require.Equal(t, expected, strings.TrimSpace(str))
	require.EqualValues(t, values, vals)
}
