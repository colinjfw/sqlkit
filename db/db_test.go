// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func wrap(t *testing.T, cb func(db DB)) {
	testdbDriver := os.Getenv("SQLKIT_DRIVER")
	if testdbDriver == "" {
		testdbDriver = "sqlite3"
	}
	testdbConn := os.Getenv("SQLKIT_CONN")
	if testdbConn == "" {
		testdbConn = "file::memory:?mode=memory&cache=shared"
	}

	fmt.Printf("driver=%s conn=%s\n", testdbDriver, testdbConn)

	testdb, err := Open(testdbDriver, testdbConn, WithLogger(StdLogger))
	require.Nil(t, err)

	testdb.Exec(context.Background(), Raw("drop table users"))
	err = testdb.Exec(context.Background(), Raw("create table users (id int primary key)")).Err()
	require.Nil(t, err)

	cb(testdb)

	require.Nil(t, testdb.Close())
}

func testSQL(t *testing.T, expected string, values []interface{}, sql SQL) {
	spew.Dump(sql)
	str, vals, err := sql.SQL()
	require.Nil(t, err)
	require.Equal(t, expected, strings.TrimSpace(str))
	require.EqualValues(t, values, vals)
}

func TestDB_Query(t *testing.T) {
	wrap(t, func(db DB) {
		ctx := context.Background()

		err := db.Query(ctx, db.Select("*").From("users")).Err()
		require.Nil(t, err)
	})
}

func TestDB_Insert(t *testing.T) {
	wrap(t, func(db DB) {
		ctx := context.Background()

		err := db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
		require.Nil(t, err)

		err = db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
		require.NotNil(t, err)
	})
}

func TestDB_CreateUpdateDelete(t *testing.T) {
	wrap(t, func(db DB) {
		ctx := context.Background()

		err := db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
		require.Nil(t, err)

		err = db.Exec(ctx, db.Update("users").Value("id", 3).Where("id = ?", 1)).Err()
		require.Nil(t, err)

		err = db.Exec(ctx, db.Delete().From("users").Where("id = ?", 3)).Err()
		require.Nil(t, err)
	})
}

func TestDB_Unmarshal(t *testing.T) {
	wrap(t, func(db DB) {
		ctx := context.Background()

		obj := []*struct {
			ID int `sql:"id"`
		}{}

		err := db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
		require.Nil(t, err)

		err = db.Query(ctx, db.Select("users.*").From("users")).Decode(&obj)
		require.Nil(t, err)

		require.Equal(t, 1, obj[0].ID)
	})
}

func TestDB_Marshal(t *testing.T) {
	wrap(t, func(db DB) {
		ctx := context.Background()

		obj := struct {
			ID int `sql:"id"`
		}{ID: 1}

		err := db.Exec(ctx, db.Insert().Into("users").Record(obj)).Err()
		require.Nil(t, err)

		out := struct {
			Count int `sql:"count"`
		}{}
		err = db.
			Query(ctx, db.Select("count(*) as count").From("users")).
			Decode(&out)
		require.Equal(t, 1, out.Count)
		require.Nil(t, err)
	})
}

func TestDB_TxBegin(t *testing.T) {
	wrap(t, func(db DB) {
		ctx, err := db.Begin(context.Background())
		require.Nil(t, err)

		defer func() { require.Nil(t, ctx.Rollback()) }()

		err = db.Exec(ctx, db.Select("*").From("users")).Err()
		require.Nil(t, err)

		require.Nil(t, ctx.Commit())
	})
}

func TestDB_TxRollback(t *testing.T) {
	wrap(t, func(db DB) {
		ctx, err := db.Begin(context.Background())
		require.Nil(t, err)

		err = db.Exec(ctx, db.Insert().Into("users").Value("id", 1)).Err()
		require.Nil(t, err)

		require.Nil(t, ctx.Rollback())

		var count int
		err = db.Query(context.Background(), Raw("SELECT COUNT(*) FROM users")).Decode(&count)
		require.Nil(t, err)
		require.Equal(t, 0, count)
	})
}

func TestDB_TxAlreadyDone(t *testing.T) {
	wrap(t, func(db DB) {
		ctx, err := db.Begin(context.Background())
		require.Nil(t, err)

		err = db.Exec(ctx, Raw("INSERT INTO users (id) VALUES (1)")).Err()
		require.Nil(t, err)

		require.Nil(t, ctx.Rollback())

		err = db.Exec(ctx, Raw("INSERT INTO users (id) VALUES (2)")).Err()
		require.NotNil(t, err)
	})
}

func TestDB_TxNested(t *testing.T) {
	wrap(t, func(db DB) {
		parent, err := db.Begin(context.Background())
		require.Nil(t, err)

		defer func() { require.Nil(t, parent.Rollback()) }()

		func() {
			var ctx TX
			ctx, err = db.Begin(parent)
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
			var ctx TX
			ctx, err = db.Begin(parent)
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
	})
}

func TestDB_TxNestedDouble(t *testing.T) {
	wrap(t, func(db DB) {
		parent, err := db.Begin(context.Background())
		require.Nil(t, err)

		defer func() { require.Nil(t, parent.Rollback()) }()

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
	})
}

func TestDB_TxCancel(t *testing.T) {
	wrap(t, func(db DB) {
		parent, err := db.Begin(context.Background())
		require.Nil(t, err)

		defer func() { require.Nil(t, parent.Rollback()) }()

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
	})
}
