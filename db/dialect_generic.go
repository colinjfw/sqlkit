package db

import (
	"strconv"
	"strings"
)

type genericMapper struct {
	bindType int
}

func (m genericMapper) query(q SelectStmt) string {
	var sql string
	sql += "SELECT " + strings.Join(q.columns, ",") + " "
	sql += "FROM " + q.table + " "
	for _, join := range q.join {
		sql += join[0] + " JOIN " + join[1] + " ON " + join[2] + " "
	}
	if q.where != "" {
		sql += "WHERE ( " + q.where + " ) "
	}
	if q.groupBy != nil {
		sql += "GROUP BY " + strings.Join(q.groupBy, ", ") + " "
	}
	if q.orderBy != nil {
		sql += "ORDER BY " + strings.Join(q.orderBy, ", ") + " "
	}
	if q.limit != "" {
		sql += "LIMIT " + q.limit + " "
	}
	if q.offset != "" {
		sql += "OFFSET " + q.offset + " "
	}
	return rebind(m.bindType, sql)
}

func (m genericMapper) insert(q InsertStmt) string {
	var sql string
	sql += "INSERT INTO " + q.table + " "
	sql += "(" + strings.Join(q.columns, ", ") + ") "
	sql += "VALUES "
	for i, row := range q.rows {
		sql += questions(len(row))
		if i != len(q.rows)-1 {
			sql += ", "
		}
	}
	return rebind(m.bindType, sql)
}

func (m genericMapper) update(q UpdateStmt) string {
	var sql string
	sql += "UPDATE " + q.table + " SET "
	for i := range q.columns {
		sql += q.columns[i] + "=?"
		if i == len(q.columns)-1 {
			sql += " "
		} else {
			sql += ", "
		}
	}
	if q.where != "" {
		sql += "WHERE " + q.where
	}
	return rebind(m.bindType, sql)
}

func questions(count int) string {
	qs := []string{}
	for i := 0; i < count; i++ {
		qs = append(qs, "?")
	}
	return "(" + strings.Join(qs, ", ") + ")"
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
