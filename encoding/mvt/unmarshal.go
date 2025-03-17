package mvt

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/pbr"
)

type decoder struct {
	keys     []string
	values   []interface{}
	features [][]byte
	valMsg   *pbr.Message
	tags     *pbr.Iterator
	geom     *pbr.Iterator
}

// geomDecoder holds state for geometry decoding.
type geomDecoder struct {
	iter  *pbr.Iterator
	count int
	used  int
	prev  geo.Point
}

func (gd *geomDecoder) NextPoint() (geo.Point, error) {
	gd.used += 2
	v, err := gd.iter.Uint32()
	if err != nil {
		return geo.Point{}, err
	}
	gd.prev[0] += unzigzag(v)

	v, err = gd.iter.Uint32()
	if err != nil {
		return geo.Point{}, err
	}
	gd.prev[1] += unzigzag(v)

	return gd.prev, nil
}
