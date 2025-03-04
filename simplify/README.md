# geo/simplify [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/simplify)

Package *simplify* implements several reducing/simplifing function for `geo.Geometry` types.   
Currently implemented:
 - [Douglas-Peucker](#dp)
 - [Visvalingam](#vis)
 - [Radial](#radial)

**Note:** The geometry object can be modified, use `Clone()` if a copy is required.

## <a name="dp"></a>[Douglas-Peucker-Algorithm](http://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm)

Probably the most popular simplification algorithm.

The algorithm is pass-through for 1d geometry such as Point and MultiPoint.
Algorithms can modify the original geometry, use `Clone()` if a copy is required.

Usage:
```go
original := geo.LineString{}
reduced := simplify.DouglasPeucker(threshold).Simplify(original.Clone())
```
