# encoding/ewkb [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/encoding/ewkb)

Package **ewkb** provides encoding and decoding of [extended WKB](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry#Format_variations) data. This format includes [SRID](https://en.wikipedia.org/wiki/Spatial_reference_system) in the data. If SRID is not needed, use the [wkb](../wkb) package for a simpler interface.   
The interface is defined as:

```go
func Marshal(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) ([]byte, error)
func MarshalToHex(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) (string, error)
func MustMarshal(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) []byte
func MustMarshalToHex(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) string

func NewEncoder(w io.Writer) *Encoder
func (e *Encoder) SetByteOrder(bo binary.ByteOrder) *Encoder
func (e *Encoder) SetSRID(srid int) *Encoder
func (e *Encoder) Encode(geom geo.Geometry) error

func Unmarshal(b []byte) (geo.Geometry, int, error)

func NewDecoder(r io.Reader) *Decoder
func (d *Decoder) Decode() (geo.Geometry, int, error)
```

## Inserting geometry into a database

Depending on the database different formats and functions are supported:

### PostgreSQL and PostGIS

PostGIS stores geometry as EWKB internally.
As a result it can be inserted without a wrapper function.

```go
db.Exec("INSERT INTO geodata(geom) VALUES (ST_GeomFromEWKB($1))", ewkb.Value(coord, 4326))

db.Exec("INSERT INTO geodata(geom) VALUES ($1)", ewkb.Value(coord, 4326))
```

### MySQL/MariaDB

MySQL and MariaDB [store geometry](https://dev.mysql.com/doc/refman/5.7/en/gis-data-formats.html) data in WKB format with a 4 byte SRID prefix.

```go
coord := geo.Point{1, 2}

// as WKB in hex format
data := wkb.MustMarshalToHex(coord)
db.Exec("INSERT INTO geodata(geom) VALUES (ST_GeomFromWKB(UNHEX(?), 4326))", data)

// relying on the raw encoding
db.Exec("INSERT INTO geodata(geom) VALUES (?)", ewkb.ValuePrefixSRID(coord, 4326))
```