

# encoding
`import "github.com/coldog/sqlkit/encoding"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package encoding provides marshalling values to and from SQL. Gracefully
handles null values.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func Marshal(obj interface{}) ([]string, []interface{}, error)](#Marshal)
* [func Unmarshal(dest interface{}, rows *sql.Rows) error](#Unmarshal)
* [type Encoder](#Encoder)
  * [func NewEncoder() Encoder](#NewEncoder)
  * [func (e Encoder) Decode(dest interface{}, rows *sql.Rows) error](#Encoder.Decode)
  * [func (e Encoder) Encode(obj interface{}) ([]string, []interface{}, error)](#Encoder.Encode)
  * [func (e Encoder) Unsafe() Encoder](#Encoder.Unsafe)
  * [func (e Encoder) WithMapper(m *reflectx.Mapper) Encoder](#Encoder.WithMapper)

#### <a name="pkg-examples">Examples</a>
* [Package](#example_)

#### <a name="pkg-files">Package files</a>
[encoding.go](/src/github.com/coldog/sqlkit/encoding/encoding.go) [marshal.go](/src/github.com/coldog/sqlkit/encoding/marshal.go) [unmarshal.go](/src/github.com/coldog/sqlkit/encoding/unmarshal.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    // ErrRequiresPtr is returned if a non pointer value is passed to Decode.
    ErrRequiresPtr = errors.New("sqlkit/encoding: non pointer passed to Decode")
    // ErrMustNotBeNil is returned if a nil value is passed to Decode.
    ErrMustNotBeNil = errors.New("sqlkit/encoding: nil value passed to Decode")
    // ErrMissingDestination is returned if a destination value is missing and
    // unsafe is not configured.
    ErrMissingDestination = errors.New("sqlkit/encoding: missing destination")
    // ErrTooManyColumns is returned when too many columns are present to scan
    // into a value.
    ErrTooManyColumns = errors.New("sqlkit/encoding: too many columns to scan")
    // ErrNoRows is mirrored from the database/sql package.
    ErrNoRows = sql.ErrNoRows
)
```
``` go
var DefaultMapper = reflectx.NewMapperFunc("db", strings.ToLower)
```
DefaultMapper is the default reflectx mapper used. This uses strings.ToLower
to map field names.

``` go
var (
    // ErrDuplicateValue is returned when duplicate values exist in the struct.
    ErrDuplicateValue = errors.New("sqlkit/marshal: duplicate values")
)
```


## <a name="Marshal">func</a> [Marshal](/src/target/marshal.go?s=1104:1166#L52)
``` go
func Marshal(obj interface{}) ([]string, []interface{}, error)
```
Marshal runs the default encoder.



## <a name="Unmarshal">func</a> [Unmarshal](/src/target/unmarshal.go?s=2825:2879#L94)
``` go
func Unmarshal(dest interface{}, rows *sql.Rows) error
```
Unmarshal will run Decode with the default Decoder configuration.




## <a name="Encoder">type</a> [Encoder](/src/target/encoding.go?s=326:387#L11)
``` go
type Encoder struct {
    // contains filtered or unexported fields
}
```
Encoder manages options for encoding.







### <a name="NewEncoder">func</a> [NewEncoder](/src/target/encoding.go?s=237:262#L8)
``` go
func NewEncoder() Encoder
```
NewEncoder returns an Encoder with the default settings which are blank.





### <a name="Encoder.Decode">func</a> (Encoder) [Decode](/src/target/unmarshal.go?s=3644:3707#L111)
``` go
func (e Encoder) Decode(dest interface{}, rows *sql.Rows) error
```
Decode does the work of decoding an *sql.Rows into a struct, array or scalar
value. Depending on the value passed in the decoder will perform the
following actions:

* If an array is passed in the decoder will work through all rows


	initializing and scanning in a new instance of the value(s) for all rows.

* If the values inside the array are scalar, then the decoder will check for


	only a single column to scan and scan this in.

* If a single struct or scalar is passed in the decoder will loop once over


	the rows returning sql.ErrNoRows if this is improssible and scan the value.

The rows object is not closed after iteration is completed. The Decode
function is thread safe.




### <a name="Encoder.Encode">func</a> (Encoder) [Encode](/src/target/marshal.go?s=371:444#L16)
``` go
func (e Encoder) Encode(obj interface{}) ([]string, []interface{}, error)
```
Encode will encode to a set of fields and values using the Encoder's
settings. It will return and error if there are duplicate fields and unsafe
is not set.




### <a name="Encoder.Unsafe">func</a> (Encoder) [Unsafe](/src/target/encoding.go?s=578:611#L19)
``` go
func (e Encoder) Unsafe() Encoder
```
Unsafe configures and returns a new Encoder which uses unsafe options.
Specifically it will ignore duplicate fields on marshalling and will ignore
missing fields on unmarshalling.




### <a name="Encoder.WithMapper">func</a> (Encoder) [WithMapper](/src/target/encoding.go?s=801:856#L26)
``` go
func (e Encoder) WithMapper(m *reflectx.Mapper) Encoder
```
WithMapper configures the encoder with a reflectx.Mapper for configuring
different fields to be encoded. The DefaultMapper is used if this is not set.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
