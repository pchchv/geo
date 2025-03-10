# geo [![CI](https://github.com/pchchv/geo/workflows/CI/badge.svg)](https://github.com/pchchv/geo/actions?query=workflow%3ACI+event%3Apush) [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/geo)](https://goreportcard.com/report/github.com/pchchv/geo) [![Go Reference](https://pkg.go.dev/badge/github.com/pchchv/geo.svg)](https://pkg.go.dev/github.com/pchchv/geo)

Package *geo* defines a set of types for working with 2D geo and planar/projected geometric data in Golang.

## Features
- **GeoJSON** - support as part of the [`geojson`](geojson) sub-package.
- **Simple types** - allow for natural operations using the `make`, `append`, `len`, `[s:e]` builtins.
- **DB integration** - dDirect to type from DB query results by scanning WKB data directly into types.

## Types

```go
type Point [2]float64
type MultiPoint []Point

type LineString []Point
type MultiLineString []LineString

type Ring LineString
type Polygon []Ring
type MultiPolygon []Polygon

type Collection []Geometry

type Bound struct { Min, Max Point }
```

Defining types as slices allows them to be accessed idiomatically using Go's built-in functions such as `make`, `append`, `len`.  
And also by using slice designators such as `[s:e]` e.g.:

```go
ls := make(geo.LineString, 0, 100)
ls = append(ls, geo.Point{1, 1})
point := ls[0]
```

### Shared `Geometry` interface

All of the base types implement the `geo.Geometry` interface defined as:

```go
type Geometry interface {
    GeoJSONType() string
    Dimensions() int // e.g. 0d, 1d, 2d
    Bound() Bound
}
```

This interface is accepted by functions in the sub-packages which then act on the base types correctly, e. g.:

```go
l := clip.Geometry(bound, geom)
```

will use the appropriate clipping algorithm depending on if the input is 1d or 2d, e.g. a `geo.LineString` or a `geo.Polygon`.

Only a few methods are defined directly on these type, for example `Clone`, `Equal`, `GeoJSONType`.  
Other operation that depend on geo vs. planar contexts are defined in the respective sub-package.  
For example:

- Computing the geo distance between two point:

  ```go
  p1 := geo.Point{-72.796408, -45.407131}
  p2 := geo.Point{-72.688541, -45.384987}

  geo.Distance(p1, p2)
  ```

- Compute the planar area and centroid of a polygon:

  ```go
  poly := geo.Polygon{...}
  centroid, area := planar.CentroidArea(poly)
  ```