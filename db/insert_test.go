// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsert_SQLValues(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (c1, c2) VALUES (?, ?)",
		[]interface{}{"1", "2"},
		Insert().
			Into("users").
			Columns("c1", "c2").
			Values("1", "2"),
	)
}

func TestInsert_SQLValue(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (c1) VALUES (?)",
		[]interface{}{"1"},
		Insert().
			Into("users").
			Value("c1", "1"),
	)
}

func TestInsert_SQLRecord(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (id, name) VALUES (?, ?)",
		[]interface{}{1, "name"},
		Insert().
			Into("users").
			Record(struct {
				ID   int    `sql:"id"`
				Name string `sql:"name"`
			}{
				ID:   1,
				Name: "name",
			}),
	)
}

func TestInsert_MultipleSQLRecord(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (id, name) VALUES (?, ?), (?, ?)",
		[]interface{}{1, "test", 2, "test"},
		Insert().
			Into("users").
			Record(struct {
				ID   int    `sql:"id"`
				Name string `sql:"name"`
			}{
				ID:   1,
				Name: "test",
			}).
			Record(struct {
				ID   int    `sql:"id"`
				Name string `sql:"name"`
			}{
				ID:   2,
				Name: "test",
			}),
	)
}

func TestInsert_InvalidMultipleSQLRecord(t *testing.T) {
	_, _, err := Insert().
		Into("users").
		Record(struct {
			ID   int    `sql:"id"`
			Name string `sql:"name"`
		}{
			ID:   1,
			Name: "test",
		}).
		Record(struct {
			ID    int    `sql:"id"`
			Name  string `sql:"name"`
			Other string `sql:"other"`
		}{
			ID:   2,
			Name: "test",
		}).
		SQL()
	require.Equal(t, ErrStatementInvalid, err)
}

func TestInsert_InvalidMultipleSQLRecordCols(t *testing.T) {
	_, _, err := Insert().
		Into("users").
		Record(struct {
			ID   int    `sql:"id"`
			Name string `sql:"name"`
		}{
			ID:   1,
			Name: "test",
		}).
		Record(struct {
			ID    int    `sql:"id"`
			Other string `sql:"other"`
		}{
			ID:    2,
			Other: "test",
		}).
		SQL()
	require.Equal(t, ErrStatementInvalid, err)
}

func TestInsert_InvalidArgs(t *testing.T) {
	_, _, err := Insert().
		Into("users").
		Columns("id", "extra").
		Values(1).
		SQL()
	require.Equal(t, ErrStatementInvalid, err)
}

func TestInsert_Postgres(t *testing.T) {
	testSQL(t,
		"INSERT INTO users (c1, c2) VALUES ($1, $2)",
		[]interface{}{"1", "2"},
		InsertStmt{dialect: Postgres}.
			Into("users").
			Columns("c1", "c2").
			Values("1", "2"),
	)
}
