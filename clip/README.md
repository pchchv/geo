# geo/clip [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/clip)

Package geo/clip provide functions for clipping lines and polygons to a bounding box.
Uses:

- [Cohen-Sutherland algorithm](https://en.wikipedia.org/wiki/Cohen%E2%80%93Sutherland_algorithm) for line clipping
- [Sutherland-Hodgman algorithm](https://en.wikipedia.org/wiki/Sutherland%E2%80%93Hodgman_algorithm) for polygon clipping

## Example

```go
bound := geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{30, 30}}
ls := geo.LineString{
    {-10, 10}, {10, 10}, {10, -10}, {20, -10}, {20, 10},
    {40, 10}, {40, 20}, {20, 20}, {20, 40}, {10, 40},
    {10, 20}, {5, 20}, {-10, 20},
}

// works on and returns an geo.Geometry interface.
clipped = clip.Geometry(bound, ls)

// or clip the line string directly
clipped = clip.LineString(bound, ls)
```

<div align="right">

##### based on [lineclip](https://github.com/mapbox/lineclip)

</div>