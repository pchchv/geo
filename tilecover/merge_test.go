package tilecover

import "testing"

func TestMergeUp(t *testing.T) {
	f := loadFeature(t, "./testdata/line.geojson")
	tiles, _ := Geometry(f.Geometry, 15)
	c1 := len(MergeUpPartial(tiles, 1, 1))

	tiles, _ = Geometry(f.Geometry, 15)
	c2 := len(MergeUpPartial(tiles, 1, 2))

	tiles, _ = Geometry(f.Geometry, 15)
	c3 := len(MergeUpPartial(tiles, 1, 3))

	tiles, _ = Geometry(f.Geometry, 15)
	c4 := len(MergeUpPartial(tiles, 1, 4))

	tiles, _ = Geometry(f.Geometry, 15)
	c := len(MergeUp(tiles, 1))
	if c1 > c2 {
		t.Errorf("c1 should be bigger than c2: %v != %v", c1, c2)
	}

	if c2 > c3 {
		t.Errorf("c2 should be bigger than c3: %v != %v", c2, c3)
	}

	if c3 > c4 {
		t.Errorf("c3 should be bigger than c4: %v != %v", c3, c4)
	}

	if c4 != c {
		t.Errorf("count 4 should be same as mergeUp: %v != %v", c4, c)
	}
}
