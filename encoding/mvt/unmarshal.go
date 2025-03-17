package mvt

import (
	"errors"
	"fmt"

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
