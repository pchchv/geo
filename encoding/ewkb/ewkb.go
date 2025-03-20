package ewkb

import (
	"encoding/binary"
	"io"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

var (
	// DefaultByteOrder is the order used for marshalling or encoding is none is specified.
	DefaultByteOrder binary.ByteOrder = binary.LittleEndian
	// DefaultSRID is a common SRID representing spatial data using
	// longitude and latitude coordinates on the
	// Earth's surface as defined in the WGS84 standard,
	// which is also used for the Global Positioning System (GPS).
	// This value will be used by the encoder if it is not specified.
	DefaultSRID int = 4326
)

// Encoder encodes a geometry as EWKB to the writer given at creation time.
type Encoder struct {
	srid int
	e    *wkbcommon.Encoder
}

// NewEncoder creates a new Encoder for the given writer.
func NewEncoder(w io.Writer) *Encoder {
	e := wkbcommon.NewEncoder(w)
	e.SetByteOrder(DefaultByteOrder)
	return &Encoder{e: e, srid: DefaultSRID}
}

// SetByteOrder overrides the default byte order set when the encoder was created.
func (e *Encoder) SetByteOrder(bo binary.ByteOrder) *Encoder {
	e.e.SetByteOrder(bo)
	return e
}

// SetSRID overrides the default srid.
func (e *Encoder) SetSRID(srid int) *Encoder {
	e.srid = srid
	return e
}

// Encode writes the geometry encoded as EWKB to the given writer.
func (e *Encoder) Encode(geom geo.Geometry, srid ...int) error {
	s := e.srid
	if len(srid) > 0 {
		s = srid[0]
	}

	return e.e.Encode(geom, s)
}
