// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package encoding

import (
	"database/sql"
	"flag"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var (
	testdbDriver string
	testdbConn   string
)

func init() {
	flag.StringVar(&testdbDriver, "testdb.driver", "sqlite3", "Test database driver")
	flag.StringVar(&testdbConn, "testdb.conn", ":memory:", "Test database connection url")
	flag.Parse()
	DefaultMapper = reflectx.NewMapperFunc("db", strings.ToLower)
}

const defaultSchema = `
create table users (
	id int primary key,
	tint int,
	tfloat float,
	tbytes blob,
	tstring text,
	tbool tinyint,
	ttime timestamp,
	tjson blob
)
`

const defaultDrop = `
drop table users
`

func run(t *testing.T, schema, drop string, cb func(db *sql.DB)) {
	db, err := sql.Open(testdbDriver, testdbConn)
	require.Nil(t, err, "failed to open")

	_, err = db.Exec(schema)
	require.Nil(t, err, "failed to apply schema")

	cb(db)

	db.Exec(drop)
	db.Close()
}

func TestUnmarshal_NullTypes(t *testing.T) {
	type allTypes struct {
		ID      int
		TInt    int
		TFloat  float64
		TBytes  []byte
		TString string
		TBool   bool
		TTime   time.Time
		TJSON   types.JSONText
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(`insert into users (id) values (?)`, 1)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := []*allTypes{}
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		require.Equal(t, "", dest[0].TString)
		spew.Dump(dest)
	})
}

func TestUnmarshal_PtrNullTypes(t *testing.T) {
	type allTypes struct {
		ID      *int
		TInt    *int
		TFloat  *float64
		TBytes  *[]byte
		TString *string
		TBool   *bool
		TTime   *time.Time
		TJSON   *types.JSONText
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(`insert into users (id) values (?)`, 1)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := []*allTypes{}
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		spew.Dump(dest)
		require.Equal(t, 1, *dest[0].ID)
		require.Nil(t, dest[0].TString)
	})
}

func TestUnmarshal_PresentTypes(t *testing.T) {
	type allTypes struct {
		ID      int
		TInt    int
		TFloat  float64
		TBytes  []byte
		TString string
		TBool   bool
		TTime   time.Time
		TJSON   types.JSONText
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := []*allTypes{}
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		require.Equal(t, "hello", dest[0].TString)
		spew.Dump(dest)
	})
}

func TestUnmarshal_PresentPtrTypes(t *testing.T) {
	type allTypes struct {
		ID      *int
		TInt    *int
		TFloat  *float64
		TBytes  *[]byte
		TString *string
		TBool   *bool
		TTime   *time.Time
		TJSON   *types.JSONText
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := []*allTypes{}
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		require.Equal(t, "hello", *dest[0].TString)
		spew.Dump(dest)
	})
}

func TestUnmarshal_SingleRow(t *testing.T) {
	type allTypes struct {
		ID      *int
		TInt    *int
		TFloat  *float64
		TBytes  *[]byte
		TString *string
		TBool   *bool
		TTime   *time.Time
		TJSON   *types.JSONText
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := &allTypes{}
		err = Unmarshal(dest, rows)
		require.Nil(t, err)

		require.Equal(t, "hello", *dest.TString)
		spew.Dump(dest)
	})
}

func TestUnmarshal_ListNonPtrStruct(t *testing.T) {
	type allTypes struct {
		ID      *int
		TInt    *int
		TFloat  *float64
		TBytes  *[]byte
		TString *string
		TBool   *bool
		TTime   *time.Time
		TJSON   *types.JSONText
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := []allTypes{}
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		spew.Dump(dest)
	})
}

func TestUnmarshal_Embedded(t *testing.T) {
	type BaseType struct {
		ID int
	}
	type allTypes struct {
		BaseType
		TInt    int
		TFloat  float64
		TBytes  []byte
		TString string
		TBool   bool
		TTime   time.Time
		TJSON   types.JSONText
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := &allTypes{}
		err = Unmarshal(dest, rows)
		require.Nil(t, err)

		require.Equal(t, 1, dest.ID)
		require.Equal(t, "hello", dest.TString)
		spew.Dump(dest)
	})
}

func TestUnmarshal_OtherStruct(t *testing.T) {
	type BaseType struct {
		ID int
	}
	type OtherType struct {
		Boring string
	}
	type allTypes struct {
		BaseType
		TInt    int
		TFloat  float64
		TBytes  []byte
		TString string
		TBool   bool
		TTime   time.Time
		TJSON   types.JSONText
		Other   OtherType
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := &allTypes{}
		err = Unmarshal(dest, rows)
		require.Nil(t, err)

		require.Equal(t, 1, dest.ID)
		require.Equal(t, "hello", dest.TString)
		spew.Dump(dest)
	})
}

func TestUnmarshal_Missing(t *testing.T) {
	type allTypes struct {
		ID   int
		TInt int
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := &allTypes{}
		err = Unmarshal(dest, rows)
		require.Equal(t, ErrMissingDestination, err)
	})
}

func TestUnmarshal_MissingUnsafe(t *testing.T) {
	type allTypes struct {
		ID   int
		TInt int
	}

	const insert = `
	insert into users (id, tint, tfloat, tbytes, tstring, tbool, ttime, tjson)
	values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var args = []interface{}{
		1, 2, 2.3, []byte("hello"), "hello", true, time.Now(), types.JSONText("{}"),
	}

	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		_, err := db.Exec(insert, args...)
		require.Nil(t, err)

		rows, err := db.Query(`select * from users`)
		require.Nil(t, err)
		defer rows.Close()

		dest := &allTypes{}

		err = NewEncoder().Unsafe().Decode(dest, rows)
		require.Nil(t, err)
		require.Equal(t, 1, dest.ID)
		require.Equal(t, 2, dest.TInt)
		spew.Dump(dest)
	})
}

func TestUnmarshal_ListRaw(t *testing.T) {
	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		rows, err := db.Query(`select count(*) from users`)
		require.Nil(t, err)
		defer rows.Close()

		var dest []int
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		require.Equal(t, 0, dest[0])
		spew.Dump(dest)
	})
}

func TestUnmarshal_ListRawPtr(t *testing.T) {
	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		rows, err := db.Query(`select count(*) from users`)
		require.Nil(t, err)
		defer rows.Close()

		var dest []*int
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		require.Equal(t, 0, *dest[0])
		spew.Dump(dest)
	})
}

func TestUnmarshal_RawRow(t *testing.T) {
	run(t, defaultSchema, defaultDrop, func(db *sql.DB) {
		db.Exec(`insert into users (id) values (?), (?)`, 1, 2)
		rows, err := db.Query(`select count(*) from users`)
		require.Nil(t, err)
		defer rows.Close()

		var dest int
		err = Unmarshal(&dest, rows)
		require.Nil(t, err)

		require.Equal(t, 2, dest)
		spew.Dump(dest)
	})
}
