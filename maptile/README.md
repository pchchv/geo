# geo/maptile [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/maptile)

Package `maptile` provides types and methods for working with
[web mercator map tiles](https://en.wikipedia.org/wiki/Tiled_web_map).
It defines a tile as:

```go
type Zoom uint32

type Tile struct {
    X, Y uint32
    Z    Zoom
}
```

Functions are provided to create tiles from lon/lat points as well as
[quadkeys](https://msdn.microsoft.com/en-us/library/bb259689.aspx).  
The tile defines helper methods such as `Parent()`, `Children()`, `Siblings()`, etc.