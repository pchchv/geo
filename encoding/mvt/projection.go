package mvt

import "github.com/pchchv/geo"

type projection struct {
	ToTile  geo.Projection
	ToWGS84 geo.Projection
}

func isPowerOfTwo(n uint32) bool {
	return (n & (n - 1)) == 0
}
