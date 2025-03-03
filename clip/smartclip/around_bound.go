package smartclip

import "github.com/pchchv/geo"

//		     left  mid  right
//		top  1001  1000  1010
//		mid  0001  0000  0010
//	 bottom  0101  0100  0110
//
// on the boundary is outside
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

// pointFor returns a representative point for the side of the given bitCode.
func pointFor(b geo.Bound, code int) geo.Point {
	switch code {
	case 1:
		return geo.Point{b.Min[0], (b.Max[1] + b.Min[1]) / 2}
	case 2:
		return geo.Point{b.Max[0], (b.Max[1] + b.Min[1]) / 2}
	case 4:
		return geo.Point{(b.Max[0] + b.Min[0]) / 2, b.Min[1]}
	case 5:
		return geo.Point{b.Min[0], b.Min[1]}
	case 6:
		return geo.Point{b.Max[0], b.Min[1]}
	case 8:
		return geo.Point{(b.Max[0] + b.Min[0]) / 2, b.Max[1]}
	case 9:
		return geo.Point{b.Min[0], b.Max[1]}
	case 10:
		return geo.Point{b.Max[0], b.Max[1]}
	default:
		panic("invalid code")
	}
}
