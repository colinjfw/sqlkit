package db_test

import (
	"context"
	"fmt"

	"github.com/coldog/sqlkit/db"
)

func Example() {
	ctx := context.Background()
	d, err := db.Open("sqlite3", ":memory:", db.WithLogger(db.StdLogger))
	if err != nil {
		panic(err)
	}

	err = d.Exec(ctx, db.Raw("create table test (id int primary key)")).Err()
	if err != nil {
		panic(err)
	}

	err = d.Exec(ctx, d.Insert().Into("test").Value("id", 1)).Err()
	if err != nil {
		panic(err)
	}
	err = d.Exec(ctx, d.Insert().Into("test").Value("id", 2)).Err()
	if err != nil {
		panic(err)
	}

	var rows []int
	err = d.Query(ctx, d.Select("*").From("test")).Decode(&rows)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", rows)

	var count int
	err = d.Query(ctx, d.Select("count(*)").From("test")).Decode(&count)
	if err != nil {
		panic(err)
	}
	fmt.Println(count)

	// Output:
	// [1 2]
	// 2
}
