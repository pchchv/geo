package geojson

import "github.com/pchchv/geo"

type geometry interface {
	Geometry() geo.Geometry
}
