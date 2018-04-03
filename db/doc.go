// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

/*
Basic SQL query builder for SQL which handles marshalling and unmarshalling values using the `sqlkit/encoding` package. Features include:

* Query builder for SQL statements.
* Nested transactions using savepoints.
* Support for Postgres, MySQL and other sql flavours.
* Extensible query logging.
* Expands placeholders for IN (?) queries.

An example of common API usage:

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
*/

package db
