# geo/resample [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/resample)

Package `resample` contains a couple functions to resampling line geometry into more or less evenly spaces points.

```go
func Resample(ls geo.LineString, df geo.DistanceFunc, totalPoints int) geo.LineString

func ToInterval(ls geo.LineString, df geo.DistanceFunc, dist float64) geo.LineString
```

i.e., resampling a line string so the points are 1 planar unit apart:

```go
ls := resample.ToInterval(ls, planar.Distance, 1.0)
```