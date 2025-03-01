package mercator

import "math"

// ToGeo projects world coordinates back to geo coordinates.
func ToGeo(x, y float64, level uint32) (lng, lat float64) {
	maxtiles := float64(uint64(1 << level))
	lng = 360.0 * (x/maxtiles - 0.5)
	lat = 2.0*math.Atan(math.Exp(math.Pi-(2*math.Pi)*(y/maxtiles)))*(180.0/math.Pi) - 90.0
	return lng, lat
}

// ToPlanar converts the point to geo world coordinates at the given live.
func ToPlanar(lng, lat float64, level uint32) (x, y float64) {
	maxtiles := float64(uint64(1 << level))
	x = (lng/360.0 + 0.5) * maxtiles
	siny := math.Sin(lat * math.Pi / 180.0) // bound it because we have a top of the world problem
	if siny < -0.9999 {
		y = 0
	} else if siny > 0.9999 {
		y = maxtiles - 1
	} else {
		lat = 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		y = lat * maxtiles
	}

	return
}
