package geojson_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/geojson"
)

func ExampleFeatureCollection_foreignMembers() {
	rawJSON := []byte(`
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      }
	    ],
	    "title": "Title as Foreign Member"
	  }`)

	fc := geojson.NewFeatureCollection()
	if err := json.Unmarshal(rawJSON, &fc); err != nil {
		log.Fatalf("invalid json: %e", err)
	}

	fmt.Println(fc.Features[0].Geometry)
	fmt.Println(fc.ExtraMembers["title"])

	data, _ := json.Marshal(fc)
	fmt.Println(string(data))

	// Output:
	// [102 0.5]
	// Title as Foreign Member
	// {"features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[102,0.5]},"properties":{"prop0":"value0"}}],"title":"Title as Foreign Member","type":"FeatureCollection"}
}

// MyFeatureCollection is a depricated/no longer supported way to
// extract foreign/extra members from a feature collection.
// Now an UnmarshalJSON method, like below, is required for it to work.
type MyFeatureCollection struct {
	geojson.FeatureCollection
	Title string `json:"title"`
}

// UnmarshalJSON implemented as below is now required for the
// extra members to be decoded directly into the type.
func (fc *MyFeatureCollection) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &fc.FeatureCollection); err != nil {
		return err
	}

	fc.Title = fc.ExtraMembers.MustString("title", "")
	return nil
}

func ExampleFeatureCollection_foreignMembersCustom() {
	// this approach to handling foreign/extra members requires
	// implementing an `UnmarshalJSON` method on the new type
	rawJSON := []byte(`
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      }
	    ],
	    "title": "Title as Foreign Member"
	  }`)

	fc := &MyFeatureCollection{}
	if err := json.Unmarshal(rawJSON, &fc); err != nil {
		log.Fatalf("invalid json: %e", err)
	}

	fmt.Println(fc.FeatureCollection.Features[0].Geometry)
	fmt.Println(fc.Features[0].Geometry)
	fmt.Println(fc.Title)
	// Output:
	// [102 0.5]
	// [102 0.5]
	// Title as Foreign Member
}

func ExampleUnmarshalFeatureCollection() {
	rawJSON := []byte(`
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      }
	    ]
	  }`)

	fc, _ := geojson.UnmarshalFeatureCollection(rawJSON)

	// Geometry will be unmarshalled into the correct geo.Geometry type.
	point := fc.Features[0].Geometry.(geo.Point)
	fmt.Println(point)

	// Output:
	// [102 0.5]
}

func Example_unmarshal() {
	rawJSON := []byte(`
	  { "type": "FeatureCollection",
	    "features": [
	      { "type": "Feature",
	        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	        "properties": {"prop0": "value0"}
	      }
	    ]
	  }`)

	fc := geojson.NewFeatureCollection()
	if err := json.Unmarshal(rawJSON, &fc); err != nil {
		log.Fatalf("invalid json: %e", err)
	}

	// Geometry will be unmarshalled into the correct geo.Geometry type.
	point := fc.Features[0].Geometry.(geo.Point)
	fmt.Println(point)

	// Output:
	// [102 0.5]
}

func ExampleFeatureCollection_MarshalJSON() {
	fc := geojson.NewFeatureCollection()
	fc.Append(geojson.NewFeature(geo.Point{1, 2}))

	if _, err := fc.MarshalJSON(); err != nil {
		log.Fatalf("marshal error: %e", err)
	}

	// standard lib encoding/json package will also work
	if data, err := json.MarshalIndent(fc, "", " "); err != nil {
		log.Fatalf("marshal error: %e", err)
	} else {
		fmt.Println(string(data))
	}

	// Output:
	// {
	//  "features": [
	//   {
	//    "type": "Feature",
	//    "geometry": {
	//     "type": "Point",
	//     "coordinates": [
	//      1,
	//      2
	//     ]
	//    },
	//    "properties": null
	//   }
	//  ],
	//  "type": "FeatureCollection"
	// }
}
