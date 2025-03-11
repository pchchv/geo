package mvt

import (
	"math"
	"math/bits"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/internal/mercator"
	"github.com/pchchv/geo/maptile"
)

type projection struct {
	ToTile  geo.Projection
	ToWGS84 geo.Projection
}

func newProjection(tile maptile.Tile, extent uint32) *projection {
	if isPowerOfTwo(extent) {
		// powers of two extents allows for some more simplicity
		n := uint32(bits.TrailingZeros32(extent))
		z := uint32(tile.Z) + n
		minx, miny := float64(uint64(tile.X)<<n), float64(uint64(tile.Y)<<n)

		return &projection{
			ToTile: func(p geo.Point) geo.Point {
				x, y := mercator.ToPlanar(p[0], p[1], z)
				return geo.Point{
					math.Floor(x - minx),
					math.Floor(y - miny),
				}
			},
			ToWGS84: func(p geo.Point) geo.Point {
				lon, lat := mercator.ToGeo(p[0]+minx+0.5, p[1]+miny+0.5, z)
				return geo.Point{lon, lat}
			},
		}
	}

	return nonPowerOfTwoProjection(tile, extent)
}

func isPowerOfTwo(n uint32) bool {
	return (n & (n - 1)) == 0
}

func nonPowerOfTwoProjection(tile maptile.Tile, extent uint32) *projection {
	e, z := float64(extent), uint32(tile.Z)
	minx, miny := float64(tile.X), float64(tile.Y)

	return &projection{
		ToTile: func(p geo.Point) geo.Point {
			x, y := mercator.ToPlanar(p[0], p[1], z)
			return geo.Point{
				math.Floor((x - minx) * e),
				math.Floor((y - miny) * e),
			}
		},
		ToWGS84: func(p geo.Point) geo.Point {
			lon, lat := mercator.ToGeo((p[0]/e)+minx, (p[1]/e)+miny, z)
			return geo.Point{lon, lat}
		},
	}
}
