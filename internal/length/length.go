package length

import "github.com/pchchv/geo"

func lineStringLength(ls geo.LineString, df geo.DistanceFunc) (sum float64) {
	for i := 1; i < len(ls); i++ {
		sum += df(ls[i], ls[i-1])
	}

	return
}

func polygonLength(p geo.Polygon, df geo.DistanceFunc) (sum float64) {
	for _, r := range p {
		sum += lineStringLength(geo.LineString(r), df)
	}

	return
}
