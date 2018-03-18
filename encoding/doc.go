/*
Package encoding provides marshalling values to and from SQL. Gracefully
handles null values.

The simplest way to use this is to take advantage of the base marhsal and
unmarshal functions:

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
