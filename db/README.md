

# db
`import "github.com/coldog/sqlkit/db"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package db provides a wrapper around the base database/sql package with a
convenient sql builder, transaction handling and encoding features.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func StdLogger(s SQL)](#StdLogger)
* [type DB](#DB)
  * [func New(opts ...Option) DB](#New)
  * [func Open(driverName, dataSourceName string, opts ...Option) (DB, error)](#Open)
* [type Dialect](#Dialect)
* [type InsertStmt](#InsertStmt)
  * [func Insert() InsertStmt](#Insert)
  * [func (i InsertStmt) Columns(cols ...string) InsertStmt](#InsertStmt.Columns)
  * [func (i InsertStmt) Into(table string) InsertStmt](#InsertStmt.Into)
  * [func (i InsertStmt) Record(obj interface{}) InsertStmt](#InsertStmt.Record)
  * [func (i InsertStmt) Row(cols []string, vals []interface{}) InsertStmt](#InsertStmt.Row)
  * [func (i InsertStmt) SQL() (string, []interface{}, error)](#InsertStmt.SQL)
  * [func (i InsertStmt) Value(column string, value interface{}) InsertStmt](#InsertStmt.Value)
  * [func (i InsertStmt) Values(vals ...interface{}) InsertStmt](#InsertStmt.Values)
* [type Option](#Option)
  * [func WithConn(conn *sql.DB) Option](#WithConn)
  * [func WithDialect(dialect Dialect) Option](#WithDialect)
  * [func WithEncoder(enc encoding.Encoder) Option](#WithEncoder)
  * [func WithLogger(logger func(SQL)) Option](#WithLogger)
* [type Result](#Result)
  * [func (r *Result) Decode(val interface{}) (err error)](#Result.Decode)
  * [func (r *Result) Err() error](#Result.Err)
* [type SQL](#SQL)
  * [func Raw(sql string, args ...interface{}) SQL](#Raw)
* [type SelectStmt](#SelectStmt)
  * [func Select(cols ...string) SelectStmt](#Select)
  * [func (q SelectStmt) From(table string) SelectStmt](#SelectStmt.From)
  * [func (q SelectStmt) GroupBy(groupBy ...string) SelectStmt](#SelectStmt.GroupBy)
  * [func (q SelectStmt) InnerJoin(table, on string) SelectStmt](#SelectStmt.InnerJoin)
  * [func (q SelectStmt) Join(kind, table, on string, values ...interface{}) SelectStmt](#SelectStmt.Join)
  * [func (q SelectStmt) LeftJoin(table, on string) SelectStmt](#SelectStmt.LeftJoin)
  * [func (q SelectStmt) Limit(limit int) SelectStmt](#SelectStmt.Limit)
  * [func (q SelectStmt) Offset(offset int) SelectStmt](#SelectStmt.Offset)
  * [func (q SelectStmt) OrderBy(orderBy ...string) SelectStmt](#SelectStmt.OrderBy)
  * [func (q SelectStmt) RightJoin(table, on string) SelectStmt](#SelectStmt.RightJoin)
  * [func (q SelectStmt) SQL() (string, []interface{}, error)](#SelectStmt.SQL)
  * [func (q SelectStmt) Select(cols ...string) SelectStmt](#SelectStmt.Select)
  * [func (q SelectStmt) Where(where string, values ...interface{}) SelectStmt](#SelectStmt.Where)
* [type UpdateStmt](#UpdateStmt)
  * [func Update(table string) UpdateStmt](#Update)
  * [func (i UpdateStmt) Columns(cols ...string) UpdateStmt](#UpdateStmt.Columns)
  * [func (i UpdateStmt) Record(obj interface{}) UpdateStmt](#UpdateStmt.Record)
  * [func (i UpdateStmt) SQL() (string, []interface{}, error)](#UpdateStmt.SQL)
  * [func (i UpdateStmt) Table(table string) UpdateStmt](#UpdateStmt.Table)
  * [func (i UpdateStmt) Value(name string, val interface{}) UpdateStmt](#UpdateStmt.Value)
  * [func (i UpdateStmt) Values(vals ...interface{}) UpdateStmt](#UpdateStmt.Values)
  * [func (i UpdateStmt) Where(where string, args ...interface{}) UpdateStmt](#UpdateStmt.Where)

#### <a name="pkg-examples">Examples</a>
* [Package](#example_)

#### <a name="pkg-files">Package files</a>
[db.go](/src/github.com/coldog/sqlkit/db/db.go) [dialect.go](/src/github.com/coldog/sqlkit/db/dialect.go) [dialect_generic.go](/src/github.com/coldog/sqlkit/db/dialect_generic.go) [insert.go](/src/github.com/coldog/sqlkit/db/insert.go) [select.go](/src/github.com/coldog/sqlkit/db/select.go) [update.go](/src/github.com/coldog/sqlkit/db/update.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    // ErrStatementInvalid is returned when a statement is invalid.
    ErrStatementInvalid = errors.New("sqlkit/db: statement invalid")
    // ErrNotAQuery is returned when Decode is called on an Exec.
    ErrNotAQuery = errors.New("sqlkit/db: query was not issued")
)
```


## <a name="StdLogger">func</a> [StdLogger](/src/target/db.go?s=614:635#L24)
``` go
func StdLogger(s SQL)
```
StdLogger is a basic logger that uses the "log" package to log sql queries.




## <a name="DB">type</a> [DB](/src/target/db.go?s=2797:2990#L118)
``` go
type DB interface {
    Query(context.Context, SQL) *Result
    Exec(context.Context, SQL) *Result
    Close() error

    Select(cols ...string) SelectStmt
    Insert() InsertStmt
    Update(string) UpdateStmt
}
```
DB is the interface for the DB object.







### <a name="New">func</a> [New](/src/target/db.go?s=1494:1521#L57)
``` go
func New(opts ...Option) DB
```
New initializes a new DB agnostic to the underlying SQL connection.


### <a name="Open">func</a> [Open](/src/target/db.go?s=1759:1831#L70)
``` go
func Open(driverName, dataSourceName string, opts ...Option) (DB, error)
```
Open will call database/sql Open under the hood and configure a database with
an appropriate dialect.





## <a name="Dialect">type</a> [Dialect](/src/target/dialect.go?s=56:72#L4)
``` go
type Dialect int
```
Dialect represents the SQL dialect type.


``` go
const (
    Generic Dialect = iota
    Postgres
    MySQL
)
```
Dialect selections.










## <a name="InsertStmt">type</a> [InsertStmt](/src/target/insert.go?s=190:333#L11)
``` go
type InsertStmt struct {
    // contains filtered or unexported fields
}
```
InsertStmt represents an INSERT in SQL.







### <a name="Insert">func</a> [Insert](/src/target/insert.go?s=97:121#L8)
``` go
func Insert() InsertStmt
```
Insert constructs an InsertStmt.





### <a name="InsertStmt.Columns">func</a> (InsertStmt) [Columns](/src/target/insert.go?s=487:541#L27)
``` go
func (i InsertStmt) Columns(cols ...string) InsertStmt
```
Columns configures the columns.




### <a name="InsertStmt.Into">func</a> (InsertStmt) [Into](/src/target/insert.go?s=370:419#L21)
``` go
func (i InsertStmt) Into(table string) InsertStmt
```
Into configures the table name.




### <a name="InsertStmt.Record">func</a> (InsertStmt) [Record](/src/target/insert.go?s=1000:1054#L46)
``` go
func (i InsertStmt) Record(obj interface{}) InsertStmt
```
Record will decode using the decoder into a list of fields and values.




### <a name="InsertStmt.Row">func</a> (InsertStmt) [Row](/src/target/insert.go?s=1316:1385#L57)
``` go
func (i InsertStmt) Row(cols []string, vals []interface{}) InsertStmt
```
Row configures a single row into the insert statement. If the columns don't
match previous insert statements then an error is forwarded.




### <a name="InsertStmt.SQL">func</a> (InsertStmt) [SQL](/src/target/insert.go?s=1934:1990#L85)
``` go
func (i InsertStmt) SQL() (string, []interface{}, error)
```
SQL implements the SQL interface.




### <a name="InsertStmt.Value">func</a> (InsertStmt) [Value](/src/target/insert.go?s=768:838#L39)
``` go
func (i InsertStmt) Value(column string, value interface{}) InsertStmt
```
Value configures a single value insert.




### <a name="InsertStmt.Values">func</a> (InsertStmt) [Values](/src/target/insert.go?s=620:678#L33)
``` go
func (i InsertStmt) Values(vals ...interface{}) InsertStmt
```
Values configures a single row of values.




## <a name="Option">type</a> [Option](/src/target/db.go?s=826:850#L34)
``` go
type Option func(db *db)
```
Option represents option configurations.







### <a name="WithConn">func</a> [WithConn](/src/target/db.go?s=1174:1208#L47)
``` go
func WithConn(conn *sql.DB) Option
```
WithConn configures a custom *sql.DB connection.


### <a name="WithDialect">func</a> [WithDialect](/src/target/db.go?s=1030:1070#L42)
``` go
func WithDialect(dialect Dialect) Option
```
WithDialect configures the SQL dialect.


### <a name="WithEncoder">func</a> [WithEncoder](/src/target/db.go?s=1330:1375#L52)
``` go
func WithEncoder(enc encoding.Encoder) Option
```
WithEncoder configures a custom encoder if a different mapper were needed.


### <a name="WithLogger">func</a> [WithLogger](/src/target/db.go?s=897:937#L37)
``` go
func WithLogger(logger func(SQL)) Option
```
WithLogger configures a logging function.





## <a name="Result">type</a> [Result](/src/target/db.go?s=3037:3161#L129)
``` go
type Result struct {
    *sql.Rows
    LastID       int64
    RowsAffected int64
    // contains filtered or unexported fields
}
```
Result wraps a database/sql query result.










### <a name="Result.Decode">func</a> (\*Result) [Decode](/src/target/db.go?s=3329:3381#L141)
``` go
func (r *Result) Decode(val interface{}) (err error)
```
Decode will decode the results into an interface.




### <a name="Result.Err">func</a> (\*Result) [Err](/src/target/db.go?s=3229:3257#L138)
``` go
func (r *Result) Err() error
```
Err forwards the error that may have come from the connection.




## <a name="SQL">type</a> [SQL](/src/target/db.go?s=2693:2753#L113)
``` go
type SQL interface {
    SQL() (string, []interface{}, error)
}
```
SQL is an interface for an SQL query that contains a string of SQL and
arguments.







### <a name="Raw">func</a> [Raw](/src/target/db.go?s=2384:2429#L98)
``` go
func Raw(sql string, args ...interface{}) SQL
```
Raw implements the SQL intorerface for providing SQL queries.





## <a name="SelectStmt">type</a> [SelectStmt](/src/target/select.go?s=190:451#L11)
``` go
type SelectStmt struct {
    // contains filtered or unexported fields
}
```
SelectStmt represents a SELECT in sql.







### <a name="Select">func</a> [Select](/src/target/select.go?s=71:109#L8)
``` go
func Select(cols ...string) SelectStmt
```
Select returns a new SelectStmt.





### <a name="SelectStmt.From">func</a> (SelectStmt) [From](/src/target/select.go?s=614:663#L33)
``` go
func (q SelectStmt) From(table string) SelectStmt
```
From configures the table.




### <a name="SelectStmt.GroupBy">func</a> (SelectStmt) [GroupBy](/src/target/select.go?s=955:1012#L47)
``` go
func (q SelectStmt) GroupBy(groupBy ...string) SelectStmt
```
GroupBy configures the GROUP BY clause.




### <a name="SelectStmt.InnerJoin">func</a> (SelectStmt) [InnerJoin](/src/target/select.go?s=1768:1826#L78)
``` go
func (q SelectStmt) InnerJoin(table, on string) SelectStmt
```
InnerJoin adds a join of type INNER.




### <a name="SelectStmt.Join">func</a> (SelectStmt) [Join](/src/target/select.go?s=1538:1620#L71)
``` go
func (q SelectStmt) Join(kind, table, on string, values ...interface{}) SelectStmt
```
Join adds a join statement of a specific kind.




### <a name="SelectStmt.LeftJoin">func</a> (SelectStmt) [LeftJoin](/src/target/select.go?s=1905:1962#L83)
``` go
func (q SelectStmt) LeftJoin(table, on string) SelectStmt
```
LeftJoin adds a join of type LEFT.




### <a name="SelectStmt.Limit">func</a> (SelectStmt) [Limit](/src/target/select.go?s=1378:1425#L65)
``` go
func (q SelectStmt) Limit(limit int) SelectStmt
```
Limit configures the LIMIT clause.




### <a name="SelectStmt.Offset">func</a> (SelectStmt) [Offset](/src/target/select.go?s=1226:1275#L59)
``` go
func (q SelectStmt) Offset(offset int) SelectStmt
```
Offset configures the OFFSET clause.




### <a name="SelectStmt.OrderBy">func</a> (SelectStmt) [OrderBy](/src/target/select.go?s=1092:1149#L53)
``` go
func (q SelectStmt) OrderBy(orderBy ...string) SelectStmt
```
OrderBy configures the ORDER BY clause.




### <a name="SelectStmt.RightJoin">func</a> (SelectStmt) [RightJoin](/src/target/select.go?s=2042:2100#L88)
``` go
func (q SelectStmt) RightJoin(table, on string) SelectStmt
```
RightJoin adds a join of type RIGHT.




### <a name="SelectStmt.SQL">func</a> (SelectStmt) [SQL](/src/target/select.go?s=2178:2234#L93)
``` go
func (q SelectStmt) SQL() (string, []interface{}, error)
```
SQL implements the SQL interface.




### <a name="SelectStmt.Select">func</a> (SelectStmt) [Select](/src/target/select.go?s=497:550#L27)
``` go
func (q SelectStmt) Select(cols ...string) SelectStmt
```
Select configures the columns to select.




### <a name="SelectStmt.Where">func</a> (SelectStmt) [Where](/src/target/select.go?s=734:807#L39)
``` go
func (q SelectStmt) Where(where string, values ...interface{}) SelectStmt
```
Where configures the WHERE clause.




## <a name="UpdateStmt">type</a> [UpdateStmt](/src/target/update.go?s=220:400#L11)
``` go
type UpdateStmt struct {
    // contains filtered or unexported fields
}
```
UpdateStmt represents an UPDATE in SQL.







### <a name="Update">func</a> [Update](/src/target/update.go?s=103:139#L8)
``` go
func Update(table string) UpdateStmt
```
Update returns a new Update statement.





### <a name="UpdateStmt.Columns">func</a> (UpdateStmt) [Columns](/src/target/update.go?s=573:627#L29)
``` go
func (i UpdateStmt) Columns(cols ...string) UpdateStmt
```
Columns sets the colums for the update.




### <a name="UpdateStmt.Record">func</a> (UpdateStmt) [Record](/src/target/update.go?s=1249:1303#L55)
``` go
func (i UpdateStmt) Record(obj interface{}) UpdateStmt
```
Record will encode the struct and append the columns and returned values.




### <a name="UpdateStmt.SQL">func</a> (UpdateStmt) [SQL](/src/target/update.go?s=1521:1577#L67)
``` go
func (i UpdateStmt) SQL() (string, []interface{}, error)
```
SQL implements the SQL interface.




### <a name="UpdateStmt.Table">func</a> (UpdateStmt) [Table](/src/target/update.go?s=447:497#L23)
``` go
func (i UpdateStmt) Table(table string) UpdateStmt
```
Table configures the table for the query.




### <a name="UpdateStmt.Value">func</a> (UpdateStmt) [Value](/src/target/update.go?s=1019:1085#L48)
``` go
func (i UpdateStmt) Value(name string, val interface{}) UpdateStmt
```
Value configures a single value for the query.




### <a name="UpdateStmt.Values">func</a> (UpdateStmt) [Values](/src/target/update.go?s=878:936#L42)
``` go
func (i UpdateStmt) Values(vals ...interface{}) UpdateStmt
```
Values sets the values for the update.




### <a name="UpdateStmt.Where">func</a> (UpdateStmt) [Where](/src/target/update.go?s=698:769#L35)
``` go
func (i UpdateStmt) Where(where string, args ...interface{}) UpdateStmt
```
Where configures the WHERE block.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
