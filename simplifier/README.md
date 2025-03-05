# geo/simplifier [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/simplifier)

Package *simplifier* implements several reducing/simplifing function for `geo.Geometry` types.   
Currently implemented:
 - [Douglas-Peucker](#dp)
 - [Visvalingam-Whyatt](#vis)
 - [Radial distance](#radial)

**Note:**   
The algorithm is pass-through for 1d geometry such as Point and MultiPoint.  
The geometry object can be modified, use `Clone()` if a copy is required.

## <a name="dp"></a>[Douglas-Peucker algorithm](http://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm)

Probably the most popular simplification algorithm.

Usage:
```go
original := geo.LineString{}
reduced := simplifier.DouglasPeucker(threshold).Simplify(original.Clone())
```

## <a name="vis"></a>[Visvalingam algorithm](https://en.wikipedia.org/wiki/Visvalingam%E2%80%93Whyatt_algorithm)

Usage:

```go
original := geo.Ring{}

// will remove all whose triangle is smaller than `threshold`
reduced := simplifier.VisvalingamThreshold(threshold).Simplify(original)

// will remove points until there are only `toKeep` points left.
reduced := simplifier.VisvalingamKeep(toKeep).Simplify(original)

// One can also combine the parameters.
// Will continue to remove points until:
//  - there are no more below the threshold,
//  - or the new path is of length `toKeep`
reduced := simplifier.Visvalingam(threshold, toKeep).Simplify(original)
```

## <a name="radial"></a>[Radial distance](http://psimpl.sourceforge.net/radial-distance.html)

Radial reduces the path by removing points that are close together.

Usage:

```go
original := geo.Polygon{}

// this method uses a Euclidean distance measure.
reduced := simplifier.Radial(planar.Distance, threshold).Simplify(path)

// if the points are in the lng/lat space Radial Geo will
// compute the geo distance between the coordinates.
reduced := simplifier.Radial(geo.Distance, meters).Simplify(path)
```