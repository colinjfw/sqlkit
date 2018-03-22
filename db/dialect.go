// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

// Dialect represents the SQL dialect type.
type Dialect int

// Dialect selections.
const (
	Generic Dialect = iota
	Postgres
	MySQL
)

var dialects = map[Dialect]dialectMapper{
	Generic:  genericMapper{bindType: bindQuestion},
	Postgres: genericMapper{bindType: bindDollar},
	MySQL:    genericMapper{bindType: bindQuestion},
}

type dialectMapper interface {
	query(q SelectStmt) string
	insert(i InsertStmt) string
	update(q UpdateStmt) string
}
