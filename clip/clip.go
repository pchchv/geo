package clip

import "github.com/pchchv/geo"

func push(out geo.MultiLineString, i int, p geo.Point) geo.MultiLineString {
	if i >= len(out) {
		out = append(out, geo.LineString{})
	}

	out[i] = append(out[i], p)
	return out
}

// bitCode returns the point position relative to the bbox:
//
//		     left  mid  right
//		top  1001  1000  1010
//		mid  0001  0000  0010
//	 bottom  0101  0100  0110
func bitCode(b geo.Bound, p geo.Point) (code int) {
	if p[0] < b.Min[0] {
		code |= 1
	} else if p[0] > b.Max[0] {
		code |= 2
	}

	if p[1] < b.Min[1] {
		code |= 4
	} else if p[1] > b.Max[1] {
		code |= 8
	}

	return
}

func bitCodeOpen(b geo.Bound, p geo.Point) (code int) {
	if p[0] <= b.Min[0] {
		code |= 1
	} else if p[0] >= b.Max[0] {
		code |= 2
	}

	if p[1] <= b.Min[1] {
		code |= 4
	} else if p[1] >= b.Max[1] {
		code |= 8
	}

	return
}

// intersect a segment against one of the 4 lines that make up the bbox
func intersect(box geo.Bound, edge int, a, b geo.Point) geo.Point {
	if edge&8 != 0 { // top
		return geo.Point{a[0] + (b[0]-a[0])*(box.Max[1]-a[1])/(b[1]-a[1]), box.Max[1]}
	} else if edge&4 != 0 { // bottom
		return geo.Point{a[0] + (b[0]-a[0])*(box.Min[1]-a[1])/(b[1]-a[1]), box.Min[1]}
	} else if edge&2 != 0 { // right
		return geo.Point{box.Max[0], a[1] + (b[1]-a[1])*(box.Max[0]-a[0])/(b[0]-a[0])}
	} else if edge&1 != 0 { // left
		return geo.Point{box.Min[0], a[1] + (b[1]-a[1])*(box.Min[0]-a[0])/(b[0]-a[0])}
	} else {
		panic("no edge??")
	}
}
