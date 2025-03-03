# geo/clip/smartclip [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/clip/smartclip)

Package smartclip extends clip functionality to handle partial 2d geometries. The input polygon rings need to only intersect the bound. The algorithm will use this, plus orientation, to wrap/close the rings around the edge of the bound.  
As an example, an [OSM multipolyon relations](https://wiki.openstreetmap.org/wiki/Relation#Multipolygon) is when a ring (inner or outer) contains multiple ways, but only one of them is in the current viewport. Only paths that intersect the viewport and their orientation allow the correct shape to be drawn.

## Example

```go
bound := geo.Bound{Min: geo.Point{1, 1}, Max: geo.Point{10, 10}}
// a partial ring cutting the bound down the middle
ring := geo.Ring{{0, 0}, {11, 11}}
clipped := smartclip.Ring(bound, ring, geo.CCW)

// clipped is a multipolyon with one ring that wraps counter-clockwise
// around the top triangle of the box
// [[[[1 1] [10 10] [5.5 10] [1 10] [1 5.5] [1 1]]]]
```


