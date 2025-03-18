package wkbcommon

import "errors"

var (
	ErrNotWKB              = errors.New("wkbcommon: invalid data")         // returned when unmarshalling WKB and the data is not valid
	ErrNotWKBHeader        = errors.New("wkbcommon: invalid header data")  // returned when unmarshalling first few bytes and there is an issue
	ErrIncorrectGeometry   = errors.New("wkbcommon: incorrect geometry")   // returned when unmarshalling WKB data into the wrong type (e. g. linestring into a point)
	ErrUnsupportedGeometry = errors.New("wkbcommon: unsupported geometry") // returned when geometry type is not supported by this package
)
