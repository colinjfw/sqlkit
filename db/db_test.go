// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"context"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func testSQL(t *testing.T, expected string, values []interface{}, sql SQL) {
	str, vals, err := sql.SQL()
	require.Nil(t, err)
	require.Equal(t, expected, strings.TrimSpace(str))
	require.EqualValues(t, values, vals)
}

func TestDB_Open(t *testing.T) {
	db, err := Open("sqlite3", ":memory:")
	require.Nil(t, err)
	require.Nil(t, db.Close())
}

func TestDB_Query(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	ctx := context.Background()

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	err = db.Query(ctx, db.Select("*").From("users")).Err()
	require.Nil(t, err)
	require.Nil(t, db.Close())
}

func TestDB_Insert(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	ctx := context.Background()

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	err = db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
	require.Nil(t, err)

	err = db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
	require.NotNil(t, err)

	require.Nil(t, db.Close())
}

func TestDB_Update(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	ctx := context.Background()

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	err = db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
	require.Nil(t, err)

	err = db.Exec(ctx, db.Update("users").Value("id", 3).Where("id = ?", 1)).Err()
	require.Nil(t, err)

	require.Nil(t, db.Close())
}

func TestDB_Unmarshal(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	obj := []*struct {
		ID int `sql:"id"`
	}{}

	ctx := context.Background()

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	err = db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
	require.Nil(t, err)

	err = db.Query(ctx, db.Select("users.*").From("users")).Decode(&obj)
	require.Nil(t, err)

	require.Equal(t, 1, obj[0].ID)

	require.Nil(t, db.Close())
}

func TestDB_Marshal(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	obj := struct {
		ID int `sql:"id"`
	}{ID: 1}

	ctx := context.Background()

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	err = db.Exec(ctx, db.Insert().Into("users").Record(obj)).Err()
	require.Nil(t, err)

	out := struct {
		Count int `sql:"count"`
	}{}
	err = db.Query(ctx,
		db.Select("count(*) as count").
			From("users"),
	).Decode(&out)
	require.Equal(t, 1, out.Count)
	require.Nil(t, err)

	require.Nil(t, db.Close())
}

func TestDB_TxBegin(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	defer db.Close()

	ctx, err := db.Begin(context.Background())
	require.Nil(t, err)

	defer func() { require.Nil(t, ctx.Rollback()) }()

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	require.Nil(t, ctx.Commit())
}

func TestDB_TxRollback(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	defer db.Close()

	ctx, err := db.Begin(context.Background())
	require.Nil(t, err)

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	require.Nil(t, ctx.Rollback())
}

func TestDB_TxAlreadyDone(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	defer db.Close()

	ctx, err := db.Begin(context.Background())
	require.Nil(t, err)

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	require.Nil(t, ctx.Rollback())

	err = db.Exec(ctx,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Error(t, err)
}

func TestDB_TxNested(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	defer db.Close()

	parent, err := db.Begin(context.Background())
	require.Nil(t, err)

	defer func() { require.Nil(t, parent.Rollback()) }()

	err = db.Exec(parent,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	func() {
		ctx, err := db.Begin(parent)
		require.Nil(t, err)

		defer func() { require.Nil(t, ctx.Rollback()) }()

		err = db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
		require.Nil(t, err)

		var count int
		err = db.Query(ctx, db.Select("count(*)").From("users")).Decode(&count)
		require.Equal(t, 1, count)
		require.Nil(t, err)

		require.Nil(t, ctx.Commit())
	}()

	func() {
		ctx, err := db.Begin(parent)
		require.Nil(t, err)

		err = db.Exec(ctx, db.Insert().Into("users").Value("id", 2)).Err()
		require.Nil(t, err)

		var count int
		err = db.Query(ctx, db.Select("count(*)").From("users")).Decode(&count)
		require.Equal(t, 2, count)
		require.Nil(t, err)

		require.Nil(t, ctx.Rollback())
	}()

	var count int
	err = db.Query(parent, db.Select("count(*)").From("users")).Decode(&count)
	require.Equal(t, 1, count)
	require.Nil(t, err)

	require.Nil(t, parent.Commit())
}

func TestDB_TxNestedDouble(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	defer db.Close()

	parent, err := db.Begin(context.Background())
	require.Nil(t, err)

	defer func() { require.Nil(t, parent.Rollback()) }()

	err = db.Exec(parent,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	func() {
		outer, err := db.Begin(parent)
		require.Nil(t, err)

		defer func() { require.Nil(t, outer.Rollback()) }()

		err = db.Exec(outer, db.Insert().Into("users").Value("id", 1)).Err()
		require.Nil(t, err)

		var count int
		err = db.Query(outer, db.Select("count(*)").From("users")).Decode(&count)
		require.Equal(t, 1, count)
		require.Nil(t, err)

		func() {
			inner, err := db.Begin(outer)
			require.Nil(t, err)

			defer func() { require.Nil(t, inner.Rollback()) }()

			var count int
			err = db.Query(inner, db.Select("count(*)").From("users")).Decode(&count)
			require.Equal(t, 1, count)
			require.Nil(t, err)

			require.Nil(t, inner.Commit())
		}()

		require.Nil(t, outer.Commit())
	}()

	require.Nil(t, parent.Commit())
}

func TestDB_TxCancel(t *testing.T) {
	db, err := Open("sqlite3", ":memory:", WithLogger(StdLogger))
	require.Nil(t, err)

	defer db.Close()

	parent, err := db.Begin(context.Background())
	require.Nil(t, err)

	defer func() { require.Nil(t, parent.Rollback()) }()

	err = db.Exec(parent,
		Raw("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")).Err()
	require.Nil(t, err)

	func() {
		ctx, cancel := context.WithCancel(parent)
		defer cancel()

		tx, err := db.Begin(parent.WithContext(ctx))
		require.Nil(t, err)
		defer tx.Rollback()

		err = db.Exec(tx, db.Insert().Into("users").Value("id", 1)).Err()
		require.Nil(t, err)

		var count int
		err = db.Query(tx, db.Select("count(*)").From("users")).Decode(&count)
		require.Equal(t, 1, count)
		require.Nil(t, err)
	}()

	require.Nil(t, parent.Commit())
}
