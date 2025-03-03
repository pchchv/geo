package smartclip

import "github.com/pchchv/geo"

func deepEqualRing(r1, r2 geo.Ring) bool {
	if len(r1) != len(r2) || r1[0] != r1[len(r1)-1] || r2[0] != r2[len(r2)-1] {
		return false
	}

	// find match
	start := -1
	for i, p := range r2 {
		if p == r1[0] {
			start = i
			break
		}
	}

	if start == -1 {
		return false
	}

	for i := range r1 {
		var p2 geo.Point
		if i+start >= len(r2) {
			p2 = r2[(i+start)%len(r2)+1]
		} else {
			p2 = r2[i+start]
		}

		if r1[i] != p2 {
			return false
		}
	}

	return true
}

func deepEqualPolygon(p1, p2 geo.Polygon) bool {
	if len(p1) != len(p2) {
		return false
	}

	for i := range p1 {
		if !deepEqualRing(p1[i], p2[i]) {
			return false
		}
	}

	return true
}

func deepEqualMultiPolygon(mp1, mp2 geo.MultiPolygon) bool {
	if len(mp1) != len(mp2) {
		return false
	}

	for _, p1 := range mp1 {
		var found bool
		for _, p2 := range mp2 {
			if deepEqualPolygon(p1, p2) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
