package wkbcommon

import (
	"encoding/binary"
	"io"
)

// Encoder encodes a geometry as (E)WKB for the
// writer specified at creation.
type Encoder struct {
	buf   []byte
	w     io.Writer
	order binary.ByteOrder
}
