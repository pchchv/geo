# geo/simplify [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/simplify)

Package *simplify* implements several reducing/simplifing function for `geo.Geometry` types.   
Currently implemented:
 - [Douglas-Peucker](#dp)
 - [Visvalingam](#vis)
 - [Radial](#radial)

**Note:** The geometry object can be modified, use `Clone()` if a copy is required.
