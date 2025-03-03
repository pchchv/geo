package smartclip

import (
	"sort"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/clip"
)

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

// clipRings will take a set of rings and clip them to the boundary.
// It returns the open lineStrings with endpoints on the boundary and the closed interior rings.
func clipRings(box geo.Bound, rings []geo.Ring) (open []geo.LineString, closed []geo.Ring) {
	var result []geo.LineString
	for _, r := range rings {
		if !r.Closed() && (box.Contains(r[0]) || box.Contains(r[len(r)-1])) {
			r = append(r, r[0])
		}

		out := clip.LineString(box, geo.LineString(r), clip.OpenBound(true))
		if len(out) == 0 {
			continue // outside of bound
		}

		if r.Closed() {
			// if input is a closed ring whose endpoints are inside the bound, then join the sections
			// this operation is O(n^2), but n is the number of segments, not edges
			for i := 0; i < len(out); i++ {
				end := out[i][len(out[i])-1]
				if end[0] == box.Min[0] || box.Max[0] == end[0] || end[1] == box.Min[1] || box.Max[1] == end[1] {
					// endpoint must be within the bound to try join
					continue
				}

				for j := 0; j < len(out); j++ {
					if i == j {
						continue
					}

					if out[j][0] == end {
						out[i] = append(out[i], out[j][1:]...)
						i--
						out[j] = out[len(out)-1]
						out = out[:len(out)-1]
					}
				}
			}
		}

		result = append(result, out...)
	}

	var at int
	for _, ls := range result {
		// closed ring, so completely inside bound
		// unless it touches a boundary
		if ls[0] == ls[len(ls)-1] && pointSide(box, ls[0]) == notOnSide {
			closed = append(closed, geo.Ring(ls))
		} else {
			result[at] = ls
			at++
		}
	}

	return result[:at], closed
}

// smartWrap takes the open lineStrings with endpoints on the boundary and connects them correctly.
func smartWrap(box geo.Bound, input []geo.LineString, o geo.Orientation) (result geo.MultiPolygon) {
	points := make([]*endpoint, 0, 2*len(input)+2)
	for i, r := range input {
		// start
		points = append(points, &endpoint{
			Point:    r[0],
			Start:    true,
			Side:     pointSide(box, r[0]),
			Index:    i,
			OtherEnd: 2*i + 1,
		})

		// end
		points = append(points, &endpoint{
			Point:    r[len(r)-1],
			Start:    false,
			Side:     pointSide(box, r[len(r)-1]),
			Index:    i,
			OtherEnd: 2 * i,
		})
	}

	if o == geo.CCW {
		sort.Sort(&sortableEndpoints{
			mls: input,
			eps: points,
		})
	} else {
		sort.Sort(sort.Reverse(&sortableEndpoints{
			mls: input,
			eps: points,
		}))
	}

	var current geo.Ring
	// this operation is O(n^2)
	// it is technically possible to use a linked list and remove points instead of marking them as “used”
	// however, since n is 2 times the number of segments, lets not do this
	for i := 0; i < 2*len(points); i++ {
		ep := points[i%len(points)]
		if ep.Used {
			continue
		}

		if !ep.Start {
			if len(current) == 0 {
				current = geo.Ring(input[ep.Index])
				ep.Used = true
			}
			continue
		}

		if len(current) == 0 {
			continue
		}
		ep.Used = true

		// previous was end, connect to this start
		var r geo.Ring
		if ep.Point == current[len(current)-1] {
			r = geo.Ring{{}, {}}
		} else {
			r = aroundBound(box, geo.Ring{ep.Point, current[len(current)-1]}, o)
		}

		if ep.Point.Equal(current[0]) {
			// loop complete!!
			current = append(current, r[2:]...)
			result = append(result, geo.Polygon{current})
			current = nil
			i = -1 // start over looking for unused endpoints
		} else {
			if len(r) > 2 {
				current = append(current, r[2:len(r)-1]...)
			}

			current = append(current, input[ep.Index]...)
			points[ep.OtherEnd].Used = true
			i = ep.OtherEnd
		}
	}

	return
}

// Polygon will smart clip a polygon to the bound.
// Rings that are NOT closed AND have an endpoint in the bound will be implicitly closed.
func Polygon(box geo.Bound, p geo.Polygon, o geo.Orientation) geo.MultiPolygon {
	if len(p) == 0 {
		return nil
	}

	open, closed := clipRings(box, p)
	if len(open) == 0 {
		// nothing was clipped
		if len(closed) == 0 {
			return nil // everything outside bound
		}

		return geo.MultiPolygon{p} // everything inside bound
	}

	result := smartWrap(box, open, o)
	if len(result) == 1 {
		result[0] = append(result[0], closed...)
	} else {
		for _, i := range closed {
			result = addToMultiPolygon(result, i)
		}
	}

	return result
}

// MultiPolygon will smart clip a multipolygon to the bound.
// Rings that are NOT closed AND have an endpoint in the bound will be implicitly closed.
func MultiPolygon(box geo.Bound, mp geo.MultiPolygon, o geo.Orientation) geo.MultiPolygon {
	if len(mp) == 0 {
		return nil
	}

	// outer rings
	outerRings := make([]geo.Ring, 0, len(mp))
	for _, p := range mp {
		outerRings = append(outerRings, p[0])
	}

	outers, closedOuters := clipRings(box, outerRings)
	if len(outers) == 0 {
		// nothing was clipped
		if len(closedOuters) == 0 {
			return nil // everything outside bound
		}

		return mp // everything inside bound
	}

	// inner rings
	var innerRings []geo.Ring
	for _, p := range mp {
		for _, r := range p[1:] {
			innerRings = append(innerRings, r)
		}
	}

	inners, closedInners := clipRings(box, innerRings)
	// smart wrap everything that touches the edges
	result := smartWrap(box, append(outers, inners...), o)
	for _, o := range closedOuters {
		result = append(result, geo.Polygon{o})
	}

	for _, i := range closedInners {
		result = addToMultiPolygon(result, i)
	}

	return result
}
