package mvt

import "github.com/pchchv/geo"

// Simplify will run all the geometry of all the
// layers through the provided simplifer.
func (ls Layers) Simplify(s geo.Simplifier) {
	for _, l := range ls {
		l.Simplify(s)
	}
}

// Simplify will run the layer geometries through the simplifier.
func (l *Layer) Simplify(s geo.Simplifier) {
	var count int
	for _, f := range l.Features {
		g := s.Simplify(f.Geometry)
		if g == nil {
			continue
		}

		f.Geometry = g
		l.Features[count] = f
		count++
	}

	l.Features = l.Features[:count]
}
