# SQLKit

[![travis](https://travis-ci.org/ColDog/sqlkit.svg?branch=master)](https://travis-ci.org/ColDog/sqlkit.svg?branch=master)
[![codecov](https://codecov.io/gh/ColDog/sqlkit/branch/master/graph/badge.svg)](https://codecov.io/gh/ColDog/sqlkit)
[![goreport](https://goreportcard.com/badge/github.com/ColDog/sqlkit)](https://goreportcard.com/report/github.com/ColDog/sqlkit)
[![godoc](https://godoc.org/github.com/ColDog/sqlkit?status.svg)](https://godoc.org/github.com/ColDog/sqlkit)

Multipurpose SQL packages for GO programs.

## Status

This project is currently in an ALPHA state. The api is relatively solid for the `encoding` package but may change for the `db` package.

## Goals

Working with SQL in golang is challenging in certain respects. Some of the main challenges I've encountered have been:

* The way NULL values are handled by the sql package requires using pointers in places where often just the zero value would be used.
* The lack of nested transactions which are very valuable when trying to wrap an entire test or work with more complex transactions.
* Lacking a simple and extendable SQL builder.

This project is designed to fix some of these issues.

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

```go
tx, err := d.Begin(context.Background())
defer tx.Rollback()

err = d.Exec(tx, d.Insert().Into("test").Value("id", 2)).Err()

var rows []int
err = d.Query(tx, d.Select("*").From("test")).Decode(&rows)
fmt.Printf("%v\n", rows)

var count int
err = d.Query(tx, d.Select("count(*)").From("test")).Decode(&count)
fmt.Println(count)
```
