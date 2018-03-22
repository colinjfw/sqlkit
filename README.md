# SQLKit

[![travis](https://travis-ci.org/ColDog/sqlkit.svg?branch=master)](https://travis-ci.org/ColDog/sqlkit.svg?branch=master)
[![codecov](https://codecov.io/gh/ColDog/sqlkit/branch/master/graph/badge.svg)](https://codecov.io/gh/ColDog/sqlkit)
[![goreport](https://goreportcard.com/badge/github.com/ColDog/sqlkit)](https://goreportcard.com/report/github.com/ColDog/sqlkit)

Multipurpose SQL packages for GO programs.

## [`sqlkit/encoding`](encoding)

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

## [`sqlkit/db`](db)

Basic SQL query builder for SQL which handles marshalling and unmarshalling values.

```go
err = d.Exec(ctx, d.Insert().Into("test").Value("id", 2)).Err()

var rows []int
err = d.Query(ctx, d.Select("*").From("test")).Decode(&rows)
fmt.Printf("%v\n", rows)

var count int
err = d.Query(ctx, d.Select("count(*)").From("test")).Decode(&count)
fmt.Println(count)
```
