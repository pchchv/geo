package simplify

import "github.com/pchchv/geo"

type simplifier interface {
	simplify(l geo.LineString, area bool, withIndexMap bool) (geo.LineString, []int)
}

func runSimplify(s simplifier, ls geo.LineString, area bool) geo.LineString {
	if len(ls) <= 2 {
		return ls
	}

	ls, _ = s.simplify(ls, area, false)
	return ls
}

func lineString(s simplifier, ls geo.LineString) geo.LineString {
	return runSimplify(s, ls, false)
}

func multiLineString(s simplifier, mls geo.MultiLineString) geo.MultiLineString {
	for i := range mls {
		mls[i] = runSimplify(s, mls[i], false)
	}

	return mls
}

func polygon(s simplifier, p geo.Polygon) geo.Polygon {
	var count int
	for i := range p {
		r := geo.Ring(runSimplify(s, geo.LineString(p[i]), true))
		if i != 0 && len(r) <= 2 {
			continue
		}

		p[count] = r
		count++
	}

	return p[:count]
}

func multiPolygon(s simplifier, mp geo.MultiPolygon) geo.MultiPolygon {
	var count int
	for i := range mp {
		p := polygon(s, mp[i])
		if len(p[0]) <= 2 {
			continue
		}

		mp[count] = p
		count++
	}

	return mp[:count]
}

func ring(s simplifier, r geo.Ring) geo.Ring {
	return geo.Ring(runSimplify(s, geo.LineString(r), true))
}

func collection(s simplifier, c geo.Collection) geo.Collection {
	for i := range c {
		c[i] = simplify(s, c[i])
	}

	return c
}
