package mvt

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

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

// RemoveEmpty will remove line strings shorter/smaller than the limits.
func (l *Layer) RemoveEmpty(lineLimit, areaLimit float64) {
	var count int
	for i := 0; i < len(l.Features); i++ {
		f := l.Features[i]
		if f.Geometry == nil {
			continue
		}

		switch f.Geometry.Dimensions() {
		case 0: // point geometry
			l.Features[count] = f
			count++
		case 1: // line geometry
			if planar.Length(f.Geometry) >= lineLimit {
				l.Features[count] = f
				count++
			}
		case 2:
			if planar.Area(f.Geometry) >= areaLimit {
				l.Features[count] = f
				count++
			}
		}
	}

	l.Features = l.Features[:count]
}

// Simplify will run all the geometry of all the
// layers through the provided simplifer.
func (ls Layers) Simplify(s geo.Simplifier) {
	for _, l := range ls {
		l.Simplify(s)
	}
}

// RemoveEmpty will remove line strings shorter/smaller than the limits.
func (ls Layers) RemoveEmpty(lineLimit, areaLimit float64) {
	for _, l := range ls {
		l.RemoveEmpty(lineLimit, areaLimit)
	}
}
