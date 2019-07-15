// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"testing"
)

func TestDelete_SQLWhere(t *testing.T) {
	testSQL(t,
		"DELETE FROM users WHERE name = ?",
		[]interface{}{"test"},
		Delete().
			From("users").
			Where("name = ?", "test"),
	)
}

func TestDelete_SQLWhereStmt(t *testing.T) {
	testSQL(t,
		"DELETE FROM users WHERE (name = ?)",
		[]interface{}{1},
		Delete().
			From("users").
			Where(Eq("name", 1)),
	)
}
