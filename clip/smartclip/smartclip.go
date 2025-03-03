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

// Less sorts the points around the bound.
// First comparing what side it's on and then the
// actual point to determine the order.
// If two points are the same,
// sort by the edge attached to the point so lines that
// are "above" are shorted first.
func (e *sortableEndpoints) Less(i, j int) bool {
	if e.eps[i].Side != e.eps[j].Side {
		return e.eps[i].Side < e.eps[j].Side
	}

	switch e.eps[i].Side {
	case 1:
		if e.eps[i].Point[1] != e.eps[j].Point[1] {
			return e.eps[i].Point[1] >= e.eps[j].Point[1]
		}
		return e.eps[i].Before(e.mls)[1] >= e.eps[j].Before(e.mls)[1]
	case 2:
		if e.eps[i].Point[0] != e.eps[j].Point[0] {
			return e.eps[i].Point[0] < e.eps[j].Point[0]
		}
		return e.eps[i].Before(e.mls)[0] < e.eps[j].Before(e.mls)[0]
	case 3:
		if e.eps[i].Point[1] != e.eps[j].Point[1] {
			return e.eps[i].Point[1] < e.eps[j].Point[1]
		}
		return e.eps[i].Before(e.mls)[1] < e.eps[j].Before(e.mls)[1]
	case 4:
		if e.eps[i].Point[0] != e.eps[j].Point[0] {
			return e.eps[i].Point[0] >= e.eps[j].Point[0]
		}
		return e.eps[i].Before(e.mls)[0] >= e.eps[j].Before(e.mls)[0]
	default:
		panic("unreachable")
	}
}
