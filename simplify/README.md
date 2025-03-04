# geo/simplify [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/simplify)

Package *simplify* implements several reducing/simplifing function for `geo.Geometry` types.   
Currently implemented:
 - [Douglas-Peucker](#dp)
 - [Visvalingam](#vis)
 - [Radial](#radial)

**Note:**   
The algorithm is pass-through for 1d geometry such as Point and MultiPoint.  
The geometry object can be modified, use `Clone()` if a copy is required.

## <a name="dp"></a>[Douglas-Peucker-Algorithm](http://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm)

Probably the most popular simplification algorithm.

Usage:
```go
original := geo.LineString{}
reduced := simplify.DouglasPeucker(threshold).Simplify(original.Clone())
```

## <a name="vis"></a>[Visvalingam](https://en.wikipedia.org/wiki/Visvalingam%E2%80%93Whyatt_algorithm)

Usage:

```go
original := geo.Ring{}

// will remove all whose triangle is smaller than `threshold`
reduced := simplify.VisvalingamThreshold(threshold).Simplify(original)

// will remove points until there are only `toKeep` points left.
reduced := simplify.VisvalingamKeep(toKeep).Simplify(original)

// One can also combine the parameters.
// Will continue to remove points until:
//  - there are no more below the threshold,
//  - or the new path is of length `toKeep`
reduced := simplify.Visvalingam(threshold, toKeep).Simplify(original)
```