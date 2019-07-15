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

func TestDelete_SQLJoin(t *testing.T) {
	testSQL(t,
		"DELETE FROM users INNER JOIN other ON other.user_id = users.id WHERE name = ?",
		[]interface{}{"test"},
		Delete().
			From("users").
			Where("name = ?", "test").
			InnerJoin("other", "other.user_id = users.id"),
	)
}

func TestDelete_SQLLLeftJoin(t *testing.T) {
	testSQL(t,
		"DELETE FROM users LEFT JOIN other ON other.user_id = users.id WHERE name = ?",
		[]interface{}{"test"},
		Delete().
			From("users").
			Where("name = ?", "test").
			LeftJoin("other", "other.user_id = users.id"),
	)
}

func TestDelete_SQLLRightJoin(t *testing.T) {
	testSQL(t,
		"DELETE FROM users RIGHT JOIN other ON other.user_id = users.id WHERE name = ?",
		[]interface{}{"test"},
		Delete().
			From("users").
			Where("name = ?", "test").
			RightJoin("other", "other.user_id = users.id"),
	)
}
