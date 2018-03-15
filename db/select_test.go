package sql

import "testing"

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
		"SELECT * FROM users WHERE ( name = ? )",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where("name = ?", "test"),
	)
}

func TestSelect_SQLGroupBy(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE ( name = ? ) GROUP BY id",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where("name = ?", "test").
			GroupBy("id"),
	)
}

func TestSelect_SQLOrderBy(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE ( name = ? ) ORDER BY id, name",
		[]interface{}{"test"},
		Select("*").
			From("users").
			Where("name = ?", "test").
			OrderBy("id", "name"),
	)
}

func TestSelect_SQLOffsetLimit(t *testing.T) {
	testSQL(t,
		"SELECT * FROM users WHERE ( name = ? ) LIMIT 20 OFFSET 10",
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
