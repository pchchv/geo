package tilecover

import "github.com/pchchv/geo/maptile"

// MergeUp will merge tiles in a given set to a given min zoom.
// Tiles will only be merged if the set has all 4 siblings.
// It is assumed that all tiles in the input set will have the same zoom,
// e.g. the outputs of the Geometry function.
func MergeUp(set maptile.Set, min maptile.Zoom) maptile.Set {
	max := maptile.Zoom(1)
	for t, v := range set {
		if v {
			max = t.Z
			break
		}
	}

	if min == max {
		return set
	}

	merged := make(maptile.Set)
	for z := max; z > min; z-- {
		parentSet := make(maptile.Set)
		for t, v := range set {
			if !v {
				continue
			}

			sibs := t.Siblings()
			s0 := set[sibs[0]]
			s1 := set[sibs[1]]
			s2 := set[sibs[2]]
			s3 := set[sibs[3]]
			if s0 && s1 && s2 && s3 {
				set[sibs[0]] = false
				set[sibs[1]] = false
				set[sibs[2]] = false
				set[sibs[3]] = false
				parent := t.Parent()
				if z-1 == min {
					merged[parent] = true
				} else {
					parentSet[parent] = true
				}
			} else {
				if s0 {
					merged[sibs[0]] = true
					set[sibs[0]] = false
				}

				if s1 {
					merged[sibs[1]] = true
					set[sibs[1]] = false
				}

				if s2 {
					merged[sibs[2]] = true
					set[sibs[2]] = false
				}

				if s3 {
					merged[sibs[3]] = true
					set[sibs[3]] = false
				}
			}
		}

		set = parentSet
		if len(set) < 4 {
			for t := range set {
				merged[t] = true
			}
			break
		}
	}

	return merged
}
