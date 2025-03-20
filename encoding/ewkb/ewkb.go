package ewkb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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

// Marshal encodes the geometry with the given byte order.
// An SRID of 0 will not be included in the encoding and the
// result will be a wkb encoding of the geometry.
func Marshal(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, wkbcommon.GeomLength(geom, srid != 0)))
	e := NewEncoder(buf)
	e.SetSRID(srid)
	if len(byteOrder) > 0 {
		e.SetByteOrder(byteOrder[0])
	}

	if err := e.Encode(geom); err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, nil
	}

	return buf.Bytes(), nil
}

// MarshalToHex encodes the geometry into a hex string representation of the binary ewkb.
func MarshalToHex(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) (string, error) {
	data, err := Marshal(geom, srid, byteOrder...)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}

// MustMarshal encodes the geometry and panic on error.
func MustMarshal(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) []byte {
	if d, err := Marshal(geom, srid, byteOrder...); err != nil {
		panic(err)
	} else {
		return d
	}
}

// MustMarshalToHex encodes the geometry and panic on error.
func MustMarshalToHex(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) string {
	if d, err := MarshalToHex(geom, srid, byteOrder...); err != nil {
		panic(err)
	} else {
		return d
	}
}
