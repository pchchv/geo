package geojson

import "encoding/json"

// CustomJSONUnmarshaler can be set to have the code use a
// different json unmarshaler than the default in the standard library.
// One use case in enabling `github.com/json-iterator/go`
// with something like this:
//
//	import (
//	  jsoniter "github.com/json-iterator/go"
//	  "github.com/pchchv/geo"
//	)
//
//	var c = jsoniter.Config{
//	  EscapeHTML:              true,
//	  SortMapKeys:             false,
//	  MarshalFloatWith6Digits: true,
//	}.Froze()
//
//	geo.CustomJSONMarshaler = c
//	geo.CustomJSONUnmarshaler = c
//
// Note that any errors encountered during unmarshaling will be different.
var CustomJSONUnmarshaler interface {
	Unmarshal(data []byte, v interface{}) error
} = nil

func unmarshalJSON(data []byte, v interface{}) error {
	if CustomJSONUnmarshaler == nil {
		return json.Unmarshal(data, v)
	}

	return CustomJSONUnmarshaler.Unmarshal(data, v)
}
