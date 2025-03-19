# encoding/wkb [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/encoding/wkb)

Package **wkb** provides encoding and decoding of [WKB](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry#Well-known_binary) data.   
The interface is defined as:

```go
func Marshal(geom geo.Geometry, byteOrder ...binary.ByteOrder) ([]byte, error)
func MarshalToHex(geom geo.Geometry, byteOrder ...binary.ByteOrder) (string, error)
func MustMarshal(geom geo.Geometry, byteOrder ...binary.ByteOrder) []byte
func MustMarshalToHex(geom geo.Geometry, byteOrder ...binary.ByteOrder) string

func NewEncoder(w io.Writer) *Encoder
func (e *Encoder) SetByteOrder(bo binary.ByteOrder)
func (e *Encoder) Encode(geom geo.Geometry) error

func Unmarshal(b []byte) (geo.Geometry, error)

func NewDecoder(r io.Reader) *Decoder
func (d *Decoder) Decode() (geo.Geometry, error)
```

## Reading and Writing to a SQL database

Package provides wrappers for `geo.Geometry` types that implement `sql.Scanner` and `driver.Value`.   
For example:

```go
row := db.QueryRow("SELECT ST_AsBinary(point_column) FROM postgis_table")

var p geo.Point
err := row.Scan(wkb.Scanner(&p))

db.Exec("INSERT INTO table (point_column) VALUES (?)", wkb.Value(p))
```

The column can also be wrapped in `ST_AsEWKB`. The SRID will be ignored.   
If the type of the geometry is unknown:

```go
s := wkb.Scanner(nil)
err := row.Scan(&s)

switch g := s.Geometry.(type) {
case geo.Point:
case geo.LineString:
}
```

Scanning directly from MySQL columns is supported. By default, MySQL returns geometry data as WKB, but prefixed with a 4-byte SRID. To support this, if the data is not a WKB, the code will strip the first 4 bytes, the SRID and retry.