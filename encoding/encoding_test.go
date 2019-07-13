// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package encoding

import (
	"strings"
	"testing"

	// Drivers for multi test.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/require"
)

var (
	testdbDriver string
	testdbConn   string
)

func init() {
	testdbDriver = "sqlite3"
	testdbConn = ":memory:"
	DefaultMapper = reflectx.NewMapperFunc("db", strings.ToLower)
}

func TestUnderscore(t *testing.T) {
	require.Equal(t, "created_at", Underscore("CreatedAt"))
	require.Equal(t, "created", Underscore("Created"))
	require.Equal(t, "api", Underscore("API"))
	require.Equal(t, "test_api", Underscore("Test_API"))
	require.Equal(t, "test_upper", Underscore("TestUPPER"))
}
