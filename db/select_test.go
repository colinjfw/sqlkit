// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelect_SQLSelect(t *testing.T) {
	testSQL(t,
		"SELECT users.* FROM users",
		nil,
		Select().Select("users.*").
			From("users"),
	)
}

func TestSelect_SQLWhere(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE name = ?",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where("name = ?", "test"),
	)
}

func TestSelect_SQLWhereStmt(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE (name = ?)",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where(Eq("name", "test")),
	)
}

func TestSelect_SQLGroupBy(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE name = ? GROUP BY id",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where("name = ?", "test").
			GroupBy("id"),
	)
}

func TestSelect_SQLOrderBy(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE name = ? ORDER BY id, name",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where("name = ?", "test").
			OrderBy("id", "name"),
	)
}

func TestSelect_SQLOffsetLimit(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE name = ? LIMIT 20 OFFSET 10",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where("name = ?", "test").
			Limit(20).
			Offset(10),
	)
}

func TestSelect_SQLLeftJoin(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users LEFT JOIN groups ON users.group_id = groups.id",
		nil,
		Select("*").
			From("users").
			LeftJoin("groups", "users.group_id = groups.id"),
	)
}

func TestSelect_SQLRightJoin(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users RIGHT JOIN groups ON users.group_id = groups.id",
		nil,
		Select("*").
			From("users").
			RightJoin("groups", "users.group_id = groups.id"),
	)
}

func TestSelect_SQLInnerJoin(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users INNER JOIN groups ON users.group_id = groups.id",
		nil,
		Select("*").
			From("users").
			InnerJoin("groups", "users.group_id = groups.id"),
	)
}

func TestSelect_SQLPostgres(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE name = $1",
		[]interface{}{"test"},
		SelectStmt{dialect: Postgres}.
			Select("*").
			From("users").
			Where("name = ?", "test"),
	)
}

func TestSelect_In(t *testing.T) {
	arr := []string{"a", "b"}
	st := Select().Where("a = ? and x in ?", "b", arr).parseWhere()

	require.Equal(t, "a = ? and x in (?, ?)", st.where)
	require.Equal(t, []interface{}{"b", "a", "b"}, st.values)
}

func TestSelect_InPtr(t *testing.T) {
	arr := []string{"a", "b"}
	st := Select().Where("a = ? and x in ?", "b", &arr).parseWhere()

	require.Equal(t, "a = ? and x in (?, ?)", st.where)
	require.Equal(t, []interface{}{"b", "a", "b"}, st.values)
}

func TestSelect_InNone(t *testing.T) {
	arr := []string{"a", "b"}
	st := Select().Where("x in", arr).parseWhere()

	require.NotNil(t, st.err)
	require.Equal(t,
		"sqlkit/db: could not find matching '?' at index 0", st.err.Error())
}

func TestSelect_SQLKitchenSink(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users  JOIN other ON users.id = other.user_id  JOIN other2 ON users.id = other2.user_id WHERE ((id IN (SELECT id FROM other3)) AND ((birthday = ?) AND (name = ?))) GROUP BY id, name LIMIT 10 OFFSET 20",
		[]interface{}{2019, "george"},
		Select("*").
			From("users").
			Join("other", "users.id = other.user_id").
			Join("other2", "users.id = other2.user_id").
			Where(
				In("id", Select("id").From("other3")).And(
					EqAllMap(map[string]interface{}{
						"name":     "george",
						"birthday": 2019,
					}),
				),
			).
			GroupBy("id", "name").
			Limit(10).
			Offset(20),
	)
}

var _bSQL string

func BenchmarkFib10(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_bSQL, _, _ = Select("*").
			From("users").
			Join("other", "users.id = other.user_id").
			Join("other2", "users.id = other2.user_id").
			Where(
				In("id", Select("id").From("other3")).And(
					EqAllMap(map[string]interface{}{
						"name":     "george",
						"birthday": 2019,
					}),
				),
			).
			GroupBy("id", "name").
			Limit(10).
			Offset(20).
			SQL()
	}
}
