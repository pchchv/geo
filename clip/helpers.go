package clip

import "github.com/pchchv/geo"

// Ring clips the ring to the bounding box and returns another ring.
// This operation will modify the input by
// using as a scratch space so clone if necessary.
func Ring(b geo.Bound, r geo.Ring) geo.Ring {
	result := ring(b, r)
	if len(result) == 0 {
		return nil
	}

	return result
}
