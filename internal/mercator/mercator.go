package mercator

import "math"

// ToGeo projects world coordinates back to geo coordinates.
func ToGeo(x, y float64, level uint32) (lng, lat float64) {
	maxtiles := float64(uint64(1 << level))
	lng = 360.0 * (x/maxtiles - 0.5)
	lat = 2.0*math.Atan(math.Exp(math.Pi-(2*math.Pi)*(y/maxtiles)))*(180.0/math.Pi) - 90.0
	return lng, lat
}
