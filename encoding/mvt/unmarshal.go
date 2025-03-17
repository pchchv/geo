package mvt

import (
	"errors"
	"fmt"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/mvt/vectortile"
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

func (d *decoder) Geometry(geomType vectortile.Tile_GeomType) (geo.Geometry, error) {
	gd := &geomDecoder{iter: d.geom, count: d.geom.Count(pbr.WireTypeVarint)}
	if gd.count < 2 {
		return nil, fmt.Errorf("geom is not long enough: %v", gd.count)
	}

	switch geomType {
	case vectortile.Tile_POINT:
		return gd.decodePoint()
	case vectortile.Tile_LINESTRING:
		return gd.decodeLineString()
	case vectortile.Tile_POLYGON:
		return gd.decodePolygon()
	default:
		return nil, fmt.Errorf("unknown geometry type: %v", geomType)
	}
}

func (d *decoder) Reset() {
	d.keys = d.keys[:0]
	d.values = d.values[:0]
	d.features = d.features[:0]
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

func (gd *geomDecoder) done() bool {
	return !gd.iter.HasNext()
}

func (gd *geomDecoder) cmdAndCount() (uint32, uint32, error) {
	if gd.done() {
		return 0, 0, errors.New("no more data")
	}

	v, err := gd.iter.Uint32()
	if err != nil {
		return 0, 0, err
	}
	gd.used++

	cmd, count := v&0x07, v>>3
	if cmd != closePath {
		if v := gd.used + int(2*count); gd.count < v {
			return 0, 0, fmt.Errorf("data cut short: needed %d, have %d", v, gd.count)
		}
	}

	return cmd, count, nil
}

func (gd *geomDecoder) decodeLine() (geo.LineString, error) {
	cmd, count, err := gd.cmdAndCount()
	if err != nil {
		return nil, err
	}

	if cmd != moveTo || count != 1 {
		return nil, errors.New("first command not one moveTo")
	}

	first, err := gd.NextPoint()
	if err != nil {
		return nil, err
	}

	cmd, count, err = gd.cmdAndCount()
	if err != nil {
		return nil, err
	} else if cmd != lineTo {
		return nil, errors.New("second command not a lineTo")
	}

	ls := make(geo.LineString, 0, count+1)
	ls = append(ls, first)
	for i := uint32(0); i < count; i++ {
		p, err := gd.NextPoint()
		if err != nil {
			return nil, err
		}
		ls = append(ls, p)
	}

	return ls, nil
}

func (gd *geomDecoder) decodeLineString() (geo.Geometry, error) {
	var mls geo.MultiLineString
	for !gd.done() {
		ls, err := gd.decodeLine()
		if err != nil {
			return nil, err
		}

		if gd.done() && len(mls) == 0 {
			return ls, nil
		}

		mls = append(mls, ls)
	}

	return mls, nil
}

func (gd *geomDecoder) decodePoint() (geo.Geometry, error) {
	_, count, err := gd.cmdAndCount()
	if err != nil {
		return nil, err
	}

	if count == 1 {
		return gd.NextPoint()
	}

	mp := make(geo.MultiPoint, 0, count)
	for i := uint32(0); i < count; i++ {
		p, err := gd.NextPoint()
		if err != nil {
			return nil, err
		}
		mp = append(mp, p)
	}

	return mp, nil
}

func (gd *geomDecoder) decodePolygon() (geo.Geometry, error) {
	var mp geo.MultiPolygon
	var p geo.Polygon
	for !gd.done() {
		ls, err := gd.decodeLine()
		if err != nil {
			return nil, err
		}

		r := geo.Ring(ls)
		cmd, _, err := gd.cmdAndCount()
		if err != nil {
			return nil, err
		} else if cmd == closePath && !r.Closed() {
			r = append(r, r[0])
		}

		// figure out if new polygon
		if len(mp) == 0 && len(p) == 0 {
			p = append(p, r)
		} else {
			if r.Orientation() == geo.CCW {
				mp = append(mp, p)
				p = geo.Polygon{r}
			} else {
				p = append(p, r)
			}
		}
	}

	if len(mp) == 0 {
		return p, nil
	}

	return append(mp, p), nil
}
