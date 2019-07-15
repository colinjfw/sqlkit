// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"strconv"
	"strings"
)

type genericMapper struct {
	bindType int
}

func (m genericMapper) beginSavepoint(name string) string {
	return "SAVEPOINT " + name
}

func (m genericMapper) releaseSavepoint(name string) string {
	return "RELEASE SAVEPOINT " + name
}

func (m genericMapper) rollbackSavepoint(name string) string {
	return "ROLLBACK TO SAVEPOINT " + name
}

func (m genericMapper) query(q SelectStmt) string {
	sql := strings.Builder{}
	sql.WriteString("SELECT ")
	sql.WriteString(strings.Join(q.columns, ","))
	sql.WriteString(" FROM ")
	sql.WriteString(q.table)

	for _, join := range q.join {
		sql.WriteString(" ")
		sql.WriteString(join[0])
		sql.WriteString(" JOIN ")
		sql.WriteString(join[1])
		sql.WriteString(" ON ")
		sql.WriteString(join[2])
	}
	if q.where != "" {
		sql.WriteString(" WHERE ")
		sql.WriteString(q.where)
	}
	if q.groupBy != nil {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(q.groupBy, ", "))
	}
	if q.orderBy != nil {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(q.orderBy, ", "))
	}
	if q.limit != "" {
		sql.WriteString(" LIMIT ")
		sql.WriteString(q.limit)
	}
	if q.offset != "" {
		sql.WriteString(" OFFSET ")
		sql.WriteString(q.offset)
	}
	return rebind(m.bindType, sql.String())
}

func (m genericMapper) delete(q DeleteStmt) string {
	sql := strings.Builder{}
	sql.WriteString("DELETE FROM ")
	sql.WriteString(q.sel.table)
	sql.WriteString(" ")
	for _, join := range q.sel.join {
		sql.WriteString(join[0])
		sql.WriteString(" JOIN ")
		sql.WriteString(join[1])
		sql.WriteString(" ON ")
		sql.WriteString(join[2])
		sql.WriteString(" ")
	}
	if q.sel.where != "" {
		sql.WriteString("WHERE ")
		sql.WriteString(q.sel.where)
		sql.WriteString(" ")
	}
	return rebind(m.bindType, sql.String())
}

func (m genericMapper) insert(q InsertStmt) string {
	sql := strings.Builder{}
	sql.WriteString("INSERT INTO ")
	sql.WriteString(q.table)
	sql.WriteString(" (")
	sql.WriteString(strings.Join(q.columns, ", "))
	sql.WriteString(") VALUES ")
	for i, row := range q.rows {
		sql.WriteString(questions(len(row)))
		if i != len(q.rows)-1 {
			sql.WriteString(", ")
		}
	}
	return rebind(m.bindType, sql.String())
}

func (m genericMapper) update(q UpdateStmt) string {
	sql := strings.Builder{}
	sql.WriteString("UPDATE ")
	sql.WriteString(q.table)
	sql.WriteString(" SET ")
	for i := range q.columns {
		sql.WriteString(q.columns[i])
		sql.WriteString("=?")
		if i == len(q.columns)-1 {
			sql.WriteString(" ")
		} else {
			sql.WriteString(", ")
		}
	}
	if q.sel.where != "" {
		sql.WriteString("WHERE ")
		sql.WriteString(q.sel.where)
	}
	return rebind(m.bindType, sql.String())
}

func questions(count int) string {
	qs := strings.Builder{}
	qs.WriteString("(")
	for i := 0; i < count; i++ {
		if i == count - 1 {
			qs.WriteString("?")
		} else {
			qs.WriteString("?, ")
		}
	}
	qs.WriteString(")")
	return qs.String()
}

const (
	bindUnknown int = iota
	bindDollar
	bindNamed
	bindQuestion
)

func rebind(bindType int, query string) string {
	switch bindType {
	case bindQuestion, bindUnknown:
		return query
	}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		switch bindType {
		case bindDollar:
			rqb = append(rqb, '$')
		case bindNamed:
			rqb = append(rqb, ':', 'a', 'r', 'g')
		}

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}
