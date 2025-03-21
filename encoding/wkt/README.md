# encoding/wkt [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/encoding/wkt)

Package **encoding/wkt** provides encoding and decoding of [WKT](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry) data.   
The interface is defined as:

```go
func MarshalString(geo.Geometry) string

func Unmarshal(string) (geo.Geometry, error)
func UnmarshalPoint(string) (geo.Point, err error)
func UnmarshalMultiPoint(string) (geo.MultiPoint, err error)
func UnmarshalLineString(string) (geo.LineString, err error)
func UnmarshalMultiLineString(string) (geo.MultiLineString, err error)
func UnmarshalPolygon(string) (geo.Polygon, err error)
func UnmarshalMultiPolygon(string) (geo.MultiPolygon, err error)
func UnmarshalCollection(string) (geo.Collection, err error)
```
