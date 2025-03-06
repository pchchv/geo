# geo/tilecover [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/tilecover)

Package `tilecover` computes the covering set of tiles for `geo.Geometry`.

## Usage

```go
p := geo.Polygon{}
tiles, err := tilecover.Geometry(p, zoom)
if err != nil {
	// indicates a non-closed ring
}

for t := range tiles {
    // do something with tile
}

// to merge up to as much as possible to a specific zoom
tiles = tilecover.MergeUp(tiles, 0)
```

<div align="right">

##### based on nodejs library [tile-cover](https://github.com/mapbox/tile-cover)

</div>