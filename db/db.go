package sql

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/coldog/sqlkit/marshal"
)

func StdLogger(s SQL) {
	sql, args, err := s.SQL()
	if err != nil {
		log.Printf("sql: error %v", err)
		return
	}
	log.Printf("sql: executing %s -- %v", sql, args)
}

type DBOption func(db *db)

func WithLogger(logger func(SQL)) DBOption {
	return func(db *db) {
		db.logger = logger
	}
}

func Open(driverName, dataSourceName string, opts ...DBOption) (DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = d.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	var dialect Dialect
	switch driverName {
	case "postgres":
		dialect = Postgres
	case "mysql":
		dialect = MySQL
	default:
		dialect = Generic
	}

	out := &db{
		dialect: dialect,
		DB:      d,
		cache:   map[string]*sql.Stmt{},
		logger:  func(SQL) {},
	}
	for _, o := range opts {
		o(out)
	}
	return out, nil
}

func Raw(sql string, args ...interface{}) SQL {
	return raw{sql: sql, args: args}
}

type raw struct {
	sql  string
	args []interface{}
}

func (q raw) SQL() (string, []interface{}, error) {
	return q.sql, q.args, nil
}

type SQL interface {
	SQL() (string, []interface{}, error)
}

type DB interface {
	Query(context.Context, SQL) *Result
	Exec(context.Context, SQL) *Result
	Close() error

	Select(cols ...string) SelectStmt
	Insert() InsertStmt
	Update(string) UpdateStmt
}

type Result struct {
	*sql.Rows
	LastID       int64
	RowsAffected int64
	err          error
}

func (r *Result) Err() error { return r.err }

func (r *Result) List(val interface{}) (err error) {
	return r.unmarshal(val, false)
}

func (r *Result) First(val interface{}) (err error) {
	return r.unmarshal(val, true)
}

func (r *Result) unmarshal(val interface{}, first bool) (err error) {
	if r.Err() != nil {
		return r.Err()
	}
	if r.Rows == nil {
		return ErrNotAQuery
	}
	defer r.Rows.Close()
	if first {
		if !r.Rows.Next() {
			return sql.ErrNoRows
		}
		return marshal.UnmarshalRow(val, r.Rows)
	}
	return marshal.UnmarshalRows(val, r.Rows)
}

type db struct {
	*sql.DB

	dialect Dialect
	lock    sync.RWMutex
	cache   map[string]*sql.Stmt
	logger  func(SQL)
}

func (d *db) stmt(ctx context.Context, str string) (*sql.Stmt, error) {
	d.lock.RLock()
	s, ok := d.cache[str]
	d.lock.RUnlock()
	if ok {
		return s, nil
	}

	d.lock.Lock()
	defer d.lock.Unlock()
	s, err := d.PrepareContext(ctx, str)
	if err != nil {
		return nil, err
	}
	d.cache[str] = s
	return s, nil
}

func (d *db) Query(ctx context.Context, q SQL) *Result {
	defer d.logger(q)

	sql, args, err := q.SQL()
	if err != nil {
		return &Result{err: err}
	}
	st, err := d.stmt(ctx, sql)
	if err != nil {
		return &Result{err: err}
	}
	rows, err := st.QueryContext(ctx, args...)
	return &Result{Rows: rows, err: err}
}

func (d *db) Exec(ctx context.Context, q SQL) *Result {
	defer d.logger(q)

	sql, args, err := q.SQL()
	if err != nil {
		return &Result{err: err}
	}
	st, err := d.stmt(ctx, sql)
	if err != nil {
		return &Result{err: err}
	}
	r, err := st.ExecContext(ctx, args...)
	var lastID int64
	var affected int64
	if r != nil {
		lastID, _ = r.LastInsertId()
		affected, _ = r.RowsAffected()
	}
	return &Result{
		LastID:       lastID,
		RowsAffected: affected,
		err:          err,
	}
}

func (d *db) Select(cols ...string) SelectStmt {
	return SelectStmt{columns: cols, dialect: d.dialect}
}

func (d *db) Insert() InsertStmt {
	return InsertStmt{dialect: d.dialect}
}

func (d *db) Update(table string) UpdateStmt {
	return UpdateStmt{dialect: d.dialect, table: table}
}

func (d *db) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, s := range d.cache {
		s.Close()
	}
	return d.DB.Close()
}
