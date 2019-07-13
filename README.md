# SQLKit

[![travis](https://travis-ci.org/colinjfw/sqlkit.svg?branch=master)](https://travis-ci.org/colinjfw/sqlkit.svg?branch=master)
[![codecov](https://codecov.io/gh/colinjfw/sqlkit/branch/master/graph/badge.svg)](https://codecov.io/gh/colinjfw/sqlkit)
[![goreport](https://goreportcard.com/badge/github.com/colinjfw/sqlkit)](https://goreportcard.com/report/github.com/colinjfw/sqlkit)
[![godoc](https://godoc.org/github.com/colinjfw/sqlkit?status.svg)](https://godoc.org/github.com/colinjfw/sqlkit)
[![release](https://img.shields.io/github/release/colinjfw/sqlkit.svg)](https://github.com/colinjfw/sqlkit/releases)
[![license](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Multipurpose SQL packages for GO programs.

Working with SQL in golang is challenging in certain respects. Some of the main challenges I've encountered have been:

* The way NULL values are handled by the sql package requires using pointers in places where often just the zero value would be used.
* The lack of nested transactions which are very valuable when trying to wrap an entire test or work with more complex transactions.
* Lacking a simple and extendable SQL builder.

This project is designed to fix some of these issues. It is heavily influenced by sqlx https://github.com/jmoiron/sqlx but with more opinions on how SQL should be used in projects.

View the [example](example) to see an example of using this project.

## Status

This project is currently in an ALPHA state. The api is relatively solid for the `encoding` package but may change for the `db` package.

## Versioning

This project follows semantic versioning. Best effort will be maintained to keep backwards compatibility as soon as the api stabilizes at 1.0.0.

## Packages

### [`encoding`](encoding)

Marshal to structs to and from SQL. Gracefully handles null values.

```go
cols, vals, err := encoding.Marshal(user{1})
_, err = db.Exec(
"insert into users ("+strings.Join(cols, ",")+") values "+"(?)", vals...,
)

users := []user{}
rows, err := db.Query("select * from users")
err = encoding.Unmarsal(&users, rows)
```

### [`db`](db)

Basic SQL query builder for SQL which handles marshalling and unmarshalling values using the `sqlkit/encoding` package. Features include:

* Query builder for SQL statements.
* Nested transactions using savepoints.
* Support for Postgres, MySQL and other sql flavours.
* Extensible query logging.
* Expands placeholders for IN (?) queries.

An example of common API usage:

```go
d, err := db.Open("sqlite3", ":memory:", db.WithLogger(db.StdLogger))

tx, err := d.Begin(context.Background()) // tx implements context.Context.
defer tx.Rollback()

err = d.Exec(tx, d.Insert().Into("test").Value("id", 2)).Err()

var rows []int
err = d.Query(tx, d.Select("*").From("test")).Decode(&rows)
fmt.Printf("%v\n", rows) // Can decode a slice of objects or scalars.

var count int
err = d.Query(tx, d.Select("count(*)").From("test")).Decode(&count)
fmt.Println(count) // Can decode single values.
```

## Overview

### Transactions

Transactions are implemented as a wrapper over `context.Context` which allows us
to build transaction unaware database logic in our applications. Instead of
calling `tx.Exec` on a transaction we call `db.Exec(ctx, ...)` where context
will contain transaction information that we can use.

### Nested Transactions

Nested transactions are a feature used heavily in Rails or other frameworks to
roll back test code to leave the database in a clean state. If you take a look
at the [example_test.go](example/example_test.go) file you can see how we can
use a setup function to rollback on error:

```go
func setup(t *testing.T) (db.TX, func()) {
	ctx, err := db.Begin(context.Background())
	require.NoError(t, err)
	return ctx, func() { ctx.Rollback() }
}
```

Transaction begin takes a `context.Context`, if this is already a transaction,
we will initialize a savepoint instead of an additional transaction. Note that
additional nesting will continue to initialize savepoints. This isn't a true
nesting of transactions and rollbacks will not bubble up. It's important to make
sure that your code handles rolling back transactions appropriately with error
handling if you need errors to bubble up.

Using `db.TX` is the safest way to manage transactions. The function handles
rollbacks and commit calls based on the error value from the callback function.

```go
db.TX(ctx, func(ctx context.Context) error {
	return db.TX(ctx, ...)
})
```

### Usage With Other SQL Generators

The `Query` and `Exec` methods take an `SQL` interface defined below:

```go
type SQL interface {
	SQL() (string, []interface{}, error)
}
```

This means that as long as a query generator implements this interface, you can
pass this into any for the `db` package functions to execute your query. This
means that you can extend query builders or bring your own query builder.

If you want to use `squirrel` for example, you can provide a simple little
helper to translate the interfaces:

```go
type SQ struct { squirrel.Sqlizer }

func (s SQ) SQL() (string, []interface{}, error) { return s.ToSql() }
```

### Encoding

The encoding package handles taking golang interfaces and converting them from
sql types to golang types. It uses `sqlx` under the hood to manage this.

One of the main differences is that it has a different philosophy from the core
golang `database/sql` package when it comes to nullable types. Mainly, nullable
types are converted into their default value in golang. Instead of requiring a
pointer.

For example:

```go
type user struct { ID string `db:"id"` }
user := &user{}

encoding.Unmarshal(user, row) // database/sql#Rows
```

In this case, if `id` in the database is `NULL` then it will scan simply as a
blank string into `user.ID`.

You can customize the mapping between struct fields and database field names by
specifiying a mapper function:

```go
encoding.NewEncoder().WithMapper(reflectx.NewMapperFunc("mytag", strings.ToLower))
```


