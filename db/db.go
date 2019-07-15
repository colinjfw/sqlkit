// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/colinjfw/sqlkit/encoding"
)

var (
	// ErrStatementInvalid is returned when a statement is invalid.
	ErrStatementInvalid = errors.New("sqlkit/db: statement invalid")
	// ErrNotAQuery is returned when Decode is called on an Exec.
	ErrNotAQuery = errors.New("sqlkit/db: query was not issued")
	// ErrNestedTransactionsNotAllowed is returned when a nested transaction
	// cannot be executed.
	ErrNestedTransactionsNotAllowed = errors.New("sqlkit/db: nested transactions not allowed")
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

// WithDisabledSavepoints will disable savepoints for a database.
func WithDisabledSavepoints(enc encoding.Encoder) Option {
	return func(db *db) { db.disableSavepoints = true }
}

// New initializes a new DB agnostic to the underlying SQL connection.
func New(opts ...Option) DB {
	out := &db{
		logger: func(SQL) {},
	}
	for _, o := range opts {
		o(out)
	}
	out.cache = getCache(out.DB)
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

// Raw implements the SQL interface for a simple string of SQL.
type Raw string

// SQL implements the SQL interface.
func (r Raw) SQL() (string, []interface{}, error) {
	return string(r), nil, nil
}

// RawWithValues returns an SQL interface that ties values to the string.
func RawWithValues(sql string, values ...interface{}) SQL {
	return sqlHolder{sql: sql, args: values}
}

type sqlHolder struct {
	sql  string
	args []interface{}
	err  error
}

func (q sqlHolder) SQL() (string, []interface{}, error) {
	return q.sql, q.args, q.err
}

// SQL is an interface for an SQL query that contains a string of SQL and
// arguments. This is the interface that must be implemented to exec or query on
// the database. Raw(...) can be used to transform a raw string into an SQL
// interface.
type SQL interface {
	SQL() (string, []interface{}, error)
}

// TX represents a transaction. It contains a context as well as commit and
// rollback calls. It implements the Context interface so it can be passed into
// the query and exec calls.
type TX interface {
	context.Context

	// Rollback will rollback the current transaction. If this is a savepoint
	// then the savepoint will be rolled back. If Rollback or Commit has already
	// been called then Rollback will take no action and return no error. If the
	// context is already closed then Rollback will return.
	Rollback() error
	// Commit will commit the current transaction. If this is a savepoint then
	// the savepoint will be released.
	Commit() error
	// WithContext will copy the transaction using a new context as the base.
	// This can be used to set a new context with a cancel or deadline.
	WithContext(ctx context.Context) TX
}

// DB is the interface for the DB object.
type DB interface {
	// Query will execute an SQL query returning a result object. If the context
	// is a transaction then this will be used to run the query.
	Query(context.Context, SQL) *Result
	// Exec will execute an SQL query returning a result object. If the context
	// is a transaction then this will be used to run the query.
	Exec(context.Context, SQL) *Result
	// Close will close the underlying DB connection.
	Close() error
	// Begin will create a new transaction. If the passed in context is a TX
	// then a savepoint will be used. If the passed in context is cancellable it
	// will monitor this context and rollback.
	Begin(context.Context) (TX, error)
	// TX provides a safe way to execute a transaction. It ensures that if an
	// error is raised Rollback() is called and if no error is raised Commit()
	// is called.
	TX(ctx context.Context, fn func(ctx context.Context) error) error
	// Select returns a SelectStmt for the dialect.
	Select(cols ...string) SelectStmt
	// Insert returns an InsertStmt for the dialect.
	Insert() InsertStmt
	// Update returns an UpdateStmt for the dialect.
	Update(string) UpdateStmt
	// Delete returns a DeleteStmt for the dialect.
	Delete() DeleteStmt
}

// Result wraps a database/sql query result. It returns the same result for both
// Exec and Query responses.
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

type preparer interface {
	PrepareContext(context.Context, string) (*sql.Stmt, error)
}

func getCache(prep preparer) *cache {
	return &cache{prep: prep, cache: map[string]*sql.Stmt{}}
}

type cache struct {
	prep  preparer
	lock  sync.RWMutex
	cache map[string]*sql.Stmt
}

func (d *cache) stmt(ctx context.Context, str string) (*sql.Stmt, error) {
	d.lock.RLock()
	s, ok := d.cache[str]
	d.lock.RUnlock()
	if ok {
		return s, nil
	}

	d.lock.Lock()
	defer d.lock.Unlock()
	s, err := d.prep.PrepareContext(ctx, str)
	if err != nil {
		return nil, err
	}
	d.cache[str] = s
	return s, nil
}

func (d *cache) close() (err error) {
	for _, st := range d.cache {
		if sErr := st.Close(); sErr != nil {
			err = sErr
		}
	}
	return
}

type db struct {
	*sql.DB

	encoder encoding.Encoder
	dialect Dialect
	lock    sync.RWMutex
	cache   *cache
	logger  func(SQL)

	disableSavepoints bool
}

type tx struct {
	context.Context
	*sql.Tx

	dialect   Dialect
	cache     *cache
	savepoint string
	logger    func(SQL)
	done      uint32
}

func (t *tx) WithContext(ctx context.Context) TX {
	return &tx{
		Context:   ctx,
		Tx:        t.Tx,
		dialect:   t.dialect,
		cache:     t.cache,
		savepoint: t.savepoint,
		logger:    t.logger,
		done:      t.done,
	}
}

// beginSavepoint will execute a savepoint for a given transaction. It will use
// the parent transaction to execute the savepoint command.
func (t *tx) beginSavepoint() (string, error) {
	name := fmt.Sprintf("s%d", time.Now().UnixNano())
	sql := dialects[t.dialect].beginSavepoint(name)
	_, err := t.Tx.Exec(sql)
	t.logger(Raw(sql))
	return name, err
}

// releaseSavepoint will execute the release savepoint command. This is used in
// committing a savepoint.
func (t *tx) releaseSavepoint() error {
	sql := dialects[t.dialect].releaseSavepoint(t.savepoint)
	_, err := t.Tx.Exec(sql)
	t.logger(Raw(sql))
	return err
}

// rollbackSavepoint will execute the rollback savepoint sql.
func (t *tx) rollbackSavepoint() error {
	sql := dialects[t.dialect].rollbackSavepoint(t.savepoint)
	_, err := t.Tx.Exec(sql)
	t.logger(Raw(sql))
	return err
}

func (t *tx) awaitCtx() {
	<-t.Context.Done()
	if err := t.Rollback(); err != nil {
		t.logger(sqlHolder{err: err})
	}
}

func (t *tx) Commit() error {
	select {
	default:
	case <-t.Context.Done():
		return t.Context.Err()
	}
	if !atomic.CompareAndSwapUint32(&t.done, 0, 1) {
		return nil
	}
	if t.savepoint != "" {
		return t.releaseSavepoint()
	}
	t.logger(Raw("COMMIT"))
	return t.Tx.Commit()
}

func (t *tx) Rollback() error {
	if !atomic.CompareAndSwapUint32(&t.done, 0, 1) {
		return nil
	}
	if t.savepoint != "" {
		return t.rollbackSavepoint()
	}
	t.logger(Raw("ROLLBACK"))
	return t.Tx.Rollback()
}

func (d *db) TX(ctx context.Context, fn func(context.Context) error) error {
	tx, err := d.Begin(ctx)
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("sqlkit/db: rollback error %v: %v", rerr, err)
		}
		return err
	}
	return tx.Commit()
}

func (d *db) Begin(ctx context.Context) (TX, error) {
	if t, ok := ctx.(*tx); ok {
		if d.disableSavepoints {
			return nil, ErrNestedTransactionsNotAllowed
		}
		name, err := t.beginSavepoint()
		if err != nil {
			return nil, err
		}
		inner := &tx{
			Context:   t.Context,
			Tx:        t.Tx,
			cache:     t.cache,
			savepoint: name,
			logger:    d.logger,
		}
		if inner.Context.Done() != nil {
			go inner.awaitCtx()
		}
		return inner, nil
	}

	d.logger(Raw("BEGIN"))
	stx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	t := &tx{
		Context: ctx,
		Tx:      stx,
		cache:   getCache(stx),
		logger:  d.logger,
	}
	return t, nil
}

func (d *db) Query(ctx context.Context, q SQL) *Result {
	defer d.logger(q)

	sql, args, err := q.SQL()
	if err != nil {
		return &Result{err: err}
	}

	var c *cache
	if t, ok := ctx.(*tx); ok {
		c = t.cache
	} else {
		c = d.cache
	}

	st, err := c.stmt(ctx, sql)
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
	var c *cache
	if t, ok := ctx.(*tx); ok {
		c = t.cache
	} else {
		c = d.cache
	}
	st, err := c.stmt(ctx, sql)
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

func (d *db) Delete() DeleteStmt {
	return DeleteStmt{dialect: d.dialect}
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
	if cErr := d.cache.close(); cErr != nil {
		err = cErr
	}
	if dErr := d.DB.Close(); dErr != nil {
		err = dErr
	}
	return
}
