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
