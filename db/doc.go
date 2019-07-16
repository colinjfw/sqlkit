// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

/*
Package db provides a wrapper around database/sql DB with added functionality
and tooling. Main features include:

* Query builder for SQL statements.
* Nested transactions using savepoints.
* Support for Postgres, MySQL and other sql flavours.
* Extensible query logging.
* Expands placeholders for IN (?) queries.
*/
package db
