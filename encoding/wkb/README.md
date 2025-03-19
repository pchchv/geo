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