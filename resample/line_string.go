package resample

import "github.com/pchchv/geo"

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
