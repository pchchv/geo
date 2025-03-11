package mvt

import "github.com/pchchv/geo"

type projection struct {
	ToTile  geo.Projection
	ToWGS84 geo.Projection
}
