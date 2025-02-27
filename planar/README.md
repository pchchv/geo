# geo/planar [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/planar)

The geometries defined in the *geo* package are generic 2d geometries. Depending on what projection they are in, such as lon/lat or flat on the plane, area and distance calculations will be different. This package implements methods that assume a plane or Euclidean context.

## Examples

Area of 3-4-5 triangle:

```go
r := geo.Ring{{0, 0}, {3, 0}, {0, 4}, {0, 0}}
a := planar.Area(r)

fmt.Println(a)
// Output:
// 6
```

Distance between two points:

```go
d := planar.Distance(geo.Point{0, 0}, geo.Point{3, 4})

fmt.Println(d)
// Output:
// 5
```

Length/circumference of a 3-4-5 triangle:

```go
r := geo.Ring{{0, 0}, {3, 0}, {0, 4}, {0, 0}}
l := planar.Length(r)

fmt.Println(l)
// Output:
// 12
```