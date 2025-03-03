package smartclip

import "github.com/pchchv/geo"

type endpoint struct {
	Point    geo.Point
	Start    bool
	Used     bool
	Side     uint8
	Index    int
	OtherEnd int
}

func (e *endpoint) Before(mls []geo.LineString) geo.Point {
	ls := mls[e.Index]
	if e.Start {
		return ls[0]
	}

	return ls[len(ls)-2]
}

type sortableEndpoints struct {
	mls []geo.LineString
	eps []*endpoint
}
