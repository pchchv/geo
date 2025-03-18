package wkbcommon

import "errors"

var (
	ErrNotWKB       = errors.New("wkbcommon: invalid data")        // returned when unmarshalling WKB and the data is not valid
	ErrNotWKBHeader = errors.New("wkbcommon: invalid header data") // returned when unmarshalling first few bytes and there is an issue
)
