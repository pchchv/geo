package geo

import "math"

func roundPoints(ps []Point, f float64) {
	for i := range ps {
		ps[i][0] = math.Round(ps[i][0]*f) / f
		ps[i][1] = math.Round(ps[i][1]*f) / f
	}
}
