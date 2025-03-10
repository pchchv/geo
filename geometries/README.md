# geo/geometries [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/geometries)

Package `geometries` defines common 2d geometries.
Depending on what projection they are in,
such as lon/lat or flat on the plane,
area and distance calculations will be different.
Package `geometries` implements methods that assume lon/lat or WGS84 projection.

## Examples

Area of the [San Francisco Main Library](https://www.openstreetmap.org/way/24446086):

```go
poly := geo.Polygon{
	{
		{-122.4163816, 37.7792782},
		{-122.4162786, 37.7787626},
		{-122.4151027, 37.7789118},
		{-122.4152143, 37.7794274},
		{-122.4163816, 37.7792782},
	},
}
a := geometries.Area(poly)

fmt.Printf("%f m^2", a)
// Output:
// 6073.368008 m^2
```

Distance between two points:

```go
oakland := geo.Point{-122.270833, 37.804444}
sf := geo.Point{-122.416667, 37.783333}
d := geometries.Distance(oakland, sf)

fmt.Printf("%0.3f meters", d)
// Output:
// 13042.047 meters
```