// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate_SQLValues(t *testing.T) {
	testSQL(t,
		"UPDATE users SET c1=?, c2=?",
		[]interface{}{"1", "2"},
		Update("users").
			Columns("c1", "c2").
			Values("1", "2"),
	)
}

func TestUpdate_SQLValue(t *testing.T) {
	testSQL(t,
		"UPDATE users SET c1=?",
		[]interface{}{"1"},
		Update("users").
			Table("users").
			Value("c1", "1"),
	)
}

func TestUpdate_SQLRecord(t *testing.T) {
	testSQL(t,
		"UPDATE users SET id=?, name=?",
		[]interface{}{1, "name"},
		Update("users").
			Record(struct {
				ID   int    `sql:"id"`
				Name string `sql:"name"`
			}{
				ID:   1,
				Name: "name",
			}),
	)
}

func TestUpdate_SQLInvalidArgs(t *testing.T) {
	_, _, err := Update("users").
		Columns("id", "extra").
		Values(1).
		SQL()
	require.NotNil(t, err)
}

func TestUpdate_SQLPostgres(t *testing.T) {
	testSQL(t,
		"UPDATE users SET c1=$1, c2=$2",
		[]interface{}{"1", "2"},
		New(WithDialect(Postgres)).
			Update("users").
			Columns("c1", "c2").
			Values("1", "2"),
	)
}

func TestUpdate_SQLWhereStmt(t *testing.T) {
	testSQL(t,
		"UPDATE users SET c1=?, c2=? WHERE (c1 = ?)",
		[]interface{}{"1", "2", 1},
		Update("users").
			Columns("c1", "c2").
			Values("1", "2").
			Where(Eq("c1", 1)),
	)
}
