package clip

import "github.com/pchchv/geo"

func push(out geo.MultiLineString, i int, p geo.Point) geo.MultiLineString {
	if i >= len(out) {
		out = append(out, geo.LineString{})
	}

	out[i] = append(out[i], p)
	return out
}
