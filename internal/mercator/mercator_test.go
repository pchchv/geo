package mercator

import (
	"math"
	"testing"
)

var (
	Epsilon = 1e-6
	Cities  = [][2]float64{
		{57.09700, 9.85000}, {49.03000, -122.32000}, {39.23500, -76.17490},
		{57.20000, -2.20000}, {16.75000, -99.76700}, {5.60000, -0.16700},
		{51.66700, -176.46700}, {9.00000, 38.73330}, {-34.7666, 138.53670},
		{12.80000, 45.00000}, {42.70000, -110.86700}, {13.48167, 144.79330},
		{33.53300, -81.71700}, {42.53300, -99.85000}, {26.01670, 50.55000},
		{35.75000, -84.00000}, {51.11933, -1.15543}, {82.52000, -62.28000},
		{32.91700, -85.91700}, {31.19000, 29.95000}, {36.70000, 3.21700},
		{34.14000, -118.10700}, {32.50370, -116.45100}, {47.83400, 10.86800},
		{28.25000, 129.70000}, {16.75000, -22.95000}, {31.95000, 35.95000},
		{52.35000, 4.86660}, {13.58670, 144.93670}, {6.90000, 134.15000},
		{40.03000, 32.90000}, {33.65000, -85.78300}, {49.33000, 10.59700},
		{17.13330, -61.78330}, {-23.4333, -70.60000}, {51.21670, 4.40000},
		{29.60000, 35.01000}, {38.58330, -121.48300}, {34.16700, -97.13300},
		{45.60000, 9.15000}, {-18.3500, -70.33330}, {-7.88000, -14.42000},
		{15.28330, 38.90000}, {-25.2333, -57.51670}, {23.96500, 32.82000},
		{-36.8832, 174.75000}, {-38.0333, 144.46670}, {46.03300, 12.60000},
		{41.66700, -72.83300}, {35.45000, 139.45000}}
)

func TestScalarMercator(t *testing.T) {
	x, y := ToPlanar(0, 0, 31)
	lat, lng := ToGeo(x, y, 31)
	if lat != 0.0 {
		t.Errorf("Scalar Mercator, latitude should be 0: %f", lat)
	}

	if lng != 0.0 {
		t.Errorf("Scalar Mercator, longitude should be 0: %f", lng)
	}

	// specific case
	if x, y := ToPlanar(-87.65005229999997, 41.850033, 20); math.Floor(x) != 268988 || math.Floor(y) != 389836 {
		t.Errorf("Scalar Mercator, projection incorrect, got %v %v", x, y)
	}

	if x, y := ToPlanar(-87.65005229999997, 41.850033, 28); math.Floor(x) != 68861112 || math.Floor(y) != 99798110 {
		t.Errorf("Scalar Mercator, projection incorrect, got %v %v", x, y)
	}

	// testing level > 32 to verify correct type conversion
	for _, city := range Cities {
		x, y := ToPlanar(city[1], city[0], 35)
		lng, lat = ToGeo(x, y, 35)
		if math.IsNaN(lng) {
			t.Error("Scalar Mercator, lng is NaN")
		}

		if math.IsNaN(lat) {
			t.Error("Scalar Mercator, lat is NaN")
		}

		if math.Abs(lat-city[0]) > Epsilon {
			t.Errorf("Scalar Mercator, latitude miss match: %f != %f", lat, city[0])
		}

		if math.Abs(lng-city[1]) > Epsilon {
			t.Errorf("Scalar Mercator, longitude miss match: %f != %f", lng, city[1])
		}
	}

	// test polar regions
	if _, y := ToPlanar(0, 89.9, 32); y != (1<<32)-1 {
		t.Errorf("Scalar Mercator, top of the world error, got %v", y)
	}

	if _, y := ToPlanar(0, -89.9, 32); y != 0 {
		t.Errorf("Scalar Mercator, bottom of the world error, got %v", y)
	}
}
