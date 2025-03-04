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
