// Package db provides a wrapper around the base database/sql package with a
// convenient sql builder, transaction handling and encoding features.
package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/coldog/sqlkit/encoding"
)

var (
	// ErrStatementInvalid is returned when a statement is invalid.
	ErrStatementInvalid = errors.New("sqlkit/db: statement invalid")
	// ErrNotAQuery is returned when Decode is called on an Exec.
	ErrNotAQuery = errors.New("sqlkit/db: query was not issued")
)

// StdLogger is a basic logger that uses the "log" package to log sql queries.
func StdLogger(s SQL) {
	sql, args, err := s.SQL()
	if err != nil {
		log.Printf("sql: error %v", err)
		return
	}
	log.Printf("sql: executing %s -- %v", sql, args)
}

// Option represents option configurations.
type Option func(db *db)

// WithLogger configures a logging function.
func WithLogger(logger func(SQL)) Option {
	return func(db *db) { db.logger = logger }
}

// WithDialect configures the SQL dialect.
func WithDialect(dialect Dialect) Option {
	return func(db *db) { db.dialect = dialect }
}

// WithConn configures a custom *sql.DB connection.
func WithConn(conn *sql.DB) Option {
	return func(db *db) { db.DB = conn }
}

// WithEncoder configures a custom encoder if a different mapper were needed.
func WithEncoder(enc encoding.Encoder) Option {
	return func(db *db) { db.encoder = enc }
}

// New initializes a new DB agnostic to the underlying SQL connection.
func New(opts ...Option) DB {
	out := &db{
		cache:  map[string]*sql.Stmt{},
		logger: func(SQL) {},
	}
	for _, o := range opts {
		o(out)
	}
	return out
}

// Open will call database/sql Open under the hood and configure a database with
// an appropriate dialect.
func Open(driverName, dataSourceName string, opts ...Option) (DB, error) {
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

	opts = append(opts, WithConn(d), WithDialect(dialect))
	out := New(opts...)
	return out, nil
}

// Raw implements the SQL intorerface for providing SQL queries.
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

// SQL is an interface for an SQL query that contains a string of SQL and
// arguments.
type SQL interface {
	SQL() (string, []interface{}, error)
}

// DB is the interface for the DB object.
type DB interface {
	Query(context.Context, SQL) *Result
	Exec(context.Context, SQL) *Result
	Close() error

	Select(cols ...string) SelectStmt
	Insert() InsertStmt
	Update(string) UpdateStmt
}

// Result wraps a database/sql query result.
type Result struct {
	*sql.Rows
	LastID       int64
	RowsAffected int64
	err          error
	encoder      encoding.Encoder
}

// Err forwards the error that may have come from the connection.
func (r *Result) Err() error { return r.err }

// Decode will decode the results into an interface.
func (r *Result) Decode(val interface{}) (err error) {
	if r.Err() != nil {
		return r.Err()
	}
	if r.Rows == nil {
		return ErrNotAQuery
	}
	defer func() {
		if rErr := r.Rows.Close(); rErr != nil {
			err = rErr
		}
	}()
	return r.encoder.Decode(val, r.Rows)
}

type db struct {
	*sql.DB

	encoder encoding.Encoder
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
	return &Result{Rows: rows, err: err, encoder: d.encoder}
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
		encoder:      d.encoder,
	}
}

func (d *db) Select(cols ...string) SelectStmt {
	return SelectStmt{columns: cols, dialect: d.dialect}
}

func (d *db) Insert() InsertStmt {
	return InsertStmt{dialect: d.dialect, encoder: d.encoder}
}

func (d *db) Update(table string) UpdateStmt {
	return UpdateStmt{dialect: d.dialect, table: table, encoder: d.encoder}
}

func (d *db) Close() (err error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, s := range d.cache {
		if cErr := s.Close(); cErr != nil {
			err = cErr
		}
	}
	if dErr := d.DB.Close(); dErr != nil {
		err = dErr
	}
	return
}
