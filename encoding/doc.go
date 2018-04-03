// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

/*
Package encoding provides marshalling values to and from SQL. Gracefully
handles null values.

The simplest way to use this is to take advantage of the base marshal and
unmarshal functions:

An example below:

  db, err := sql.Open("sqlite3", ":memory:")

  type user struct {
  ID int `db:"id"`
  }

  cols, vals, err := encoding.Marshal(user{1})
  _, err = db.Exec(
  "insert into users ("+strings.Join(cols, ",")+") values "+"(?)", vals...,
  )

  users := []user{}
  rows, err := db.Query("select * from users")
  err = encoding.Unmarsal(&users, rows)

*/
package encoding
