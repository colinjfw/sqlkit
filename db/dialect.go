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

// dialects define all available dialects. Currently only differences between
// them is rebinding the query from `?` to the desired variable placeholder.
var dialects = map[Dialect]dialectMapper{
	Generic:  genericMapper{bindType: bindQuestion},
	Postgres: genericMapper{bindType: bindDollar},
	MySQL:    genericMapper{bindType: bindQuestion},
}

// dialectMapper provides a mapper for different dialects.
type dialectMapper interface {
	query(q SelectStmt) string
	insert(i InsertStmt) string
	update(q UpdateStmt) string
	delete(q DeleteStmt) string

	beginSavepoint(name string) string
	releaseSavepoint(name string) string
	rollbackSavepoint(name string) string
}
