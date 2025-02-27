package quadtree

import "github.com/pchchv/geo"

// FilterFunc is a function that filters the points to search for.
type FilterFunc func(p geo.Pointer) bool
