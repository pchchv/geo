package mvt

import (
	"github.com/pchchv/geo"
	"github.com/pchchv/geo/clip"
)

// MapboxGLDefaultExtentBound holds the default mapbox vector tile bounds used by mapbox-gl.
var MapboxGLDefaultExtentBound = geo.Bound{
	Min: geo.Point{-1 * DefaultExtent, -1 * DefaultExtent},
	Max: geo.Point{2*DefaultExtent - 1, 2*DefaultExtent - 1},
}

// Clip will clip all geometries in this layer to the given bounds.
// Will remove features that clip to an empty geometry, modifies the
// layer.Features slice in place.
func (l *Layer) Clip(box geo.Bound) {
	var at int
	for _, f := range l.Features {
		g := clip.Geometry(box, f.Geometry)
		if g != nil {
			f.Geometry = g
			l.Features[at] = f
			at++
		}
	}

	l.Features = l.Features[:at]
}

// Clip will clip all geometries in all layers to the given bounds.
func (ls Layers) Clip(box geo.Bound) {
	for _, l := range ls {
		l.Clip(box)
	}
}
