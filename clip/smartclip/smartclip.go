package smartclip

import "github.com/pchchv/geo"

const notOnSide = 0xFF

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

func (e *sortableEndpoints) Len() int {
	return len(e.eps)
}

func (e *sortableEndpoints) Swap(i, j int) {
	e.eps[e.eps[i].OtherEnd].OtherEnd, e.eps[e.eps[j].OtherEnd].OtherEnd = j, i
	e.eps[i], e.eps[j] = e.eps[j], e.eps[i]
}

//	 4
//	+-+
//
// 1 | | 3
//
//	+-+
//	 2
func pointSide(b geo.Bound, p geo.Point) uint8 {
	if p[1] == b.Max[1] {
		return 4
	} else if p[1] == b.Min[1] {
		return 2
	} else if p[0] == b.Max[0] {
		return 3
	} else if p[0] == b.Min[0] {
		return 1
	} else {
		return notOnSide
	}
}

func polygonContains(outer geo.Ring, r geo.Ring) bool {
	for _, p := range r {
		var inside bool
		x, y := p[0], p[1]
		i, j := 0, len(outer)-1
		for i < len(outer) {
			xi, yi := outer[i][0], outer[i][1]
			xj, yj := outer[j][0], outer[j][1]
			if ((yi > y) != (yj > y)) &&
				(x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
				inside = !inside
			}

			j = i
			i++
		}

		if inside {
			return true
		}
	}

	return false
}

// addToMultiPolygon does a lookup to see which polygon the ring intersects.
// This should work fine if the input is well formed.
func addToMultiPolygon(mp geo.MultiPolygon, ring geo.Ring) geo.MultiPolygon {
	for i := range mp {
		if polygonContains(mp[i][0], ring) {
			mp[i] = append(mp[i], ring)
			return mp
		}
	}

	return mp
}
