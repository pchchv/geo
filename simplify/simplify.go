package simplify

import "github.com/pchchv/geo"

type simplifier interface {
	simplify(l geo.LineString, area bool, withIndexMap bool) (geo.LineString, []int)
}
