package geojson

import "testing"

func TestPropertiesClone(t *testing.T) {
	props := Properties{
		"one": 2,
	}

	clone := props.Clone()
	if clone["one"] != 2 {
		t.Errorf("should clone properties")
	}

	clone["one"] = 3
	if props["one"] != 2 {
		t.Errorf("should clone properties")
	}
}

func propertiesTestFeature() *Feature {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {
	      "type": "Point",
	      "coordinates": [102.0, 0.5]
	    },
	    "properties": {
	      "bool":true,
	      "falsebool":false,
	      "int": 1,
	      "float64": 1.2,
	      "string":"text"
	    }
	  }`

	f, _ := UnmarshalFeature([]byte(rawJSON))
	return f
}
