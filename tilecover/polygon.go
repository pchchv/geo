package tilecover

import (
	"errors"
	"sort"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/maptile"
)

func polygon(set maptile.Set, p geo.Polygon, zoom maptile.Zoom) error {
	intersections := make([][2]uint32, 0)
	for _, r := range p {
		ring := line(set, geo.LineString(r), zoom, make([][2]uint32, 0))
		pi := len(ring) - 2
		for i := range ring {
			pi = (pi + 1) % len(ring)
			ni := (i + 1) % len(ring)
			y := ring[i][1]
			// add interesction if it's not local extremum or duplicate
			if (ring[pi][1] < y || ring[ni][1] < y) && // not local minimum
				(y < ring[pi][1] || y < ring[ni][1]) && // not local maximum
				y != ring[ni][1] {
				intersections = append(intersections, ring[i])
			}
		}
	}

	if len(intersections)%2 != 0 {
		return errors.New("tilecover: uneven intersections, ring not closed?")
	}

	// sort by y, then x
	sort.Slice(intersections, func(i, j int) bool {
		it, jt := intersections[i], intersections[j]
		if it[1] != jt[1] {
			return it[1] < jt[1]
		}

		return it[0] < jt[0]
	})

	for i := 0; i < len(intersections); i += 2 {
		// fill tiles between pairs of intersections
		y := intersections[i][1]
		for x := intersections[i][0] + 1; x < intersections[i+1][0]; x++ {
			set[maptile.New(x, y, zoom)] = true
		}
	}

	return nil
}
