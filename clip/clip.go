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
