package resample

import "github.com/pchchv/geo"

// Resample converts the line string into totalPoints-1 evenly spaced segments.
func Resample(ls geo.LineString, df geo.DistanceFunc, totalPoints int) geo.LineString {
	if totalPoints <= 0 {
		return nil
	}

	if ls, ret := resampleEdgeCases(ls, totalPoints); ret {
		return ls
	} else {
		// precomputes the total distance and intermediate distances
		total, dists := precomputeDistances(ls, df)
		return resample(ls, dists, total, totalPoints)
	}
}

func resample(ls geo.LineString, dists []float64, totalDistance float64, totalPoints int) geo.LineString {
	var dist float64
	step := 1
	points := make([]geo.Point, 1, totalPoints)
	points[0] = ls[0] // start stays the same
	currentDistance := totalDistance / float64(totalPoints-1)
	// declare here and update had nice performance benefits need to retest
	currentSeg := [2]geo.Point{}
	for i := 0; i < len(ls)-1; i++ {
		currentSeg[0] = ls[i]
		currentSeg[1] = ls[i+1]
		currentSegDistance := dists[i]
		nextDistance := dist + currentSegDistance
		for currentDistance <= nextDistance {
			// need to add a point
			percent := (currentDistance - dist) / currentSegDistance
			points = append(points, geo.Point{
				currentSeg[0][0] + percent*(currentSeg[1][0]-currentSeg[0][0]),
				currentSeg[0][1] + percent*(currentSeg[1][1]-currentSeg[0][1]),
			})

			// move to the next distance we want
			step++
			currentDistance = totalDistance * float64(step) / float64(totalPoints-1)
			if step == totalPoints-1 { // weird round off error on my machine
				currentDistance = totalDistance
			}
		}

		// past the current point in the original segment, so move to the next one
		dist = nextDistance
	}

	// end stays the same, to handle round off errors
	if totalPoints != 1 { // for 1, we want the first point
		points[totalPoints-1] = ls[len(ls)-1]
	}

	return geo.LineString(points)
}

// resampleEdgeCases is used to handle edge cases when resampling,
// i.e., not enough points and the line string is all the same point.
// Returns nil if there are no edge cases,
// true if one of these edge cases was found and handled.
func resampleEdgeCases(ls geo.LineString, totalPoints int) (geo.LineString, bool) {
	// degenerate case
	if len(ls) <= 1 {
		return ls, true
	}

	// if all the points are the same, treat as special case.
	equal := true
	for _, point := range ls {
		if !ls[0].Equal(point) {
			equal = false
			break
		}
	}

	if equal {
		if totalPoints > len(ls) {
			// extend to be requested length
			for len(ls) != totalPoints {
				ls = append(ls, ls[0])
			}

			return ls, true
		}

		// contract to be requested length
		ls = ls[:totalPoints]
		return ls, true
	}

	return ls, false
}

// precomputeDistances precomputes the total distance and intermediate distances.
func precomputeDistances(ls geo.LineString, df geo.DistanceFunc) (total float64, dists []float64) {
	dists = make([]float64, len(ls)-1)
	for i := 0; i < len(ls)-1; i++ {
		dists[i] = df(ls[i], ls[i+1])
		total += dists[i]
	}

	return
}
