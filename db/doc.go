/*
Package db provides a wrapper around the base database/sql package with a
convenient sql builder, transaction handling and encoding features.

The simplest usage of db is below:

  ctx := context.Background()
  d, err := db.Open("sqlite3", ":memory:", db.WithLogger(db.StdLogger))

  err = d.Exec(ctx, d.Insert(). Into("test"). Value("id", 2)).Err()

  var rows []int
  err = d.Query(ctx, d.Select("*").From("test")).Decode(&rows)
  fmt.Printf("%v\n", rows)

  var count int
  err = d.Query(ctx, d.Select("count(*)").From("test")).Decode(&count)
  fmt.Println(count)

*/
package db
