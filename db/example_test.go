// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/colinjfw/sqlkit/db"
)

func Example_open() {
	d, err := db.Open("sqlite3", ":memory:", db.WithLogger(db.StdLogger))
	if err != nil {
		panic(err)
	}

	ctx, err := d.Begin(context.Background())
	if err != nil {
		panic(err)
	}

	err = d.Exec(ctx, db.Raw("create table test (id int primary key)")).Err()
	if err != nil {
		panic(err)
	}

	err = d.Exec(
		ctx,
		d.Insert().
			Into("test").
			Value("id", 1).
			Value("id", 2),
	).Err()
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

func Example_transactions() {
	d, err := db.Open("sqlite3", ":memory:", db.WithLogger(db.StdLogger))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	var rows []int
	err = d.TX(ctx, func(ctx context.Context) error {
		err = d.Exec(ctx, db.Raw("create table test (id int primary key)")).Err()
		if err != nil {
			return err
		}

		// Will be executed using a savepoint.
		err = d.TX(ctx, func(ctx context.Context) error {
			return d.Exec(
				ctx,
				d.Insert().
					Into("test").
					Value("id", 1).
					Value("id", 2),
			).Err()
		})
		if err != nil {
			return err
		}

		// This savepoint will be rolled back.
		err = d.TX(ctx, func(ctx context.Context) error {
			err = d.Exec(
				ctx,
				d.Insert().
					Into("test").
					Value("id", 3),
			).Err()
			if err != nil {
				return err
			}
			return errors.New("fake an error")
		})
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		// We won't see number 3 in the count as the savepoint is rolled back.
		return d.Query(ctx, d.Select("*").From("test")).Decode(&rows)
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", rows)

	// Output:
	// err: fake an error
	// [1 2]
}
