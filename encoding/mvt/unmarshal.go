package mvt

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/mvt/vectortile"
	"github.com/pchchv/geo/geojson"
	"github.com/pchchv/pbr"
)

// Unmarshal takes Mapbox Vector Tile (MVT)
// data and converts it into a set of layers,
// but does not project coordinates.
func Unmarshal(data []byte) (Layers, error) {
	layers, err := unmarshalTile(data)
	if err != nil && dataIsGZipped(data) {
		return nil, errors.New("failed to unmarshal, data possibly gzipped")
	}

	return layers, err
}

// UnmarshalGzipped takes gzipped Mapbox Vector Tile (MVT)
// data and unzips it before decoding it into a set of layers,
// with no coordinates projected.
func UnmarshalGzipped(data []byte) (Layers, error) {
	gzreader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzreader: %e", err)
	}

	decoded, err := io.ReadAll(gzreader)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip: %e", err)
	}

	return Unmarshal(decoded)
}

type decoder struct {
	keys     []string
	values   []interface{}
	features [][]byte
	valMsg   *pbr.Message
	tags     *pbr.Iterator
	geom     *pbr.Iterator
}

func (d *decoder) Geometry(geomType vectortile.Tile_GeomType) (geo.Geometry, error) {
	gd := &geomDecoder{iter: d.geom, count: d.geom.Count(pbr.WireTypeVarint)}
	if gd.count < 2 {
		return nil, fmt.Errorf("geom is not long enough: %v", gd.count)
	}

	switch geomType {
	case vectortile.Tile_POINT:
		return gd.decodePoint()
	case vectortile.Tile_LINESTRING:
		return gd.decodeLineString()
	case vectortile.Tile_POLYGON:
		return gd.decodePolygon()
	default:
		return nil, fmt.Errorf("unknown geometry type: %v", geomType)
	}
}

func (d *decoder) Reset() {
	d.keys = d.keys[:0]
	d.values = d.values[:0]
	d.features = d.features[:0]
}

func (d *decoder) Feature(msg *pbr.Message) (feature *geojson.Feature, err error) {
	var geomType vectortile.Tile_GeomType
	feature = &geojson.Feature{Type: "Feature"}
	for msg.Next() {
		switch msg.FieldNumber() {
		case 1: // id
			if id, err := msg.Uint64(); err != nil {
				return nil, err
			} else {
				feature.ID = float64(id)
			}
		case 2: //tags, repeated packed
			d.tags, err = msg.Iterator(d.tags)
			if err != nil {
				return nil, err
			}

			count := d.tags.Count(pbr.WireTypeVarint)
			feature.Properties = make(geojson.Properties, count/2)
			for d.tags.HasNext() {
				k, err := d.tags.Uint32()
				if err != nil {
					return nil, err
				}

				v, err := d.tags.Uint32()
				if err != nil {
					return nil, err
				}

				if len(d.keys) <= int(k) || len(d.values) <= int(v) {
					continue
				}
				feature.Properties[d.keys[k]] = d.values[v]
			}
		case 3: // geomtype
			if t, err := msg.Int32(); err != nil {
				return nil, err
			} else {
				geomType = vectortile.Tile_GeomType(t)
			}
		case 4: // geometry
			if d.geom, err = msg.Iterator(d.geom); err != nil {
				return nil, err
			}
		default:
			msg.Skip()
		}
	}

	if msg.Error() != nil {
		return nil, msg.Error()
	}

	if geo, err := d.Geometry(geomType); err != nil {
		return nil, err
	} else {
		feature.Geometry = geo
	}

	return feature, nil
}

func (d *decoder) Layer(msg *pbr.Message) (layer *Layer, err error) {
	d.Reset()
	layer = &Layer{
		Version: vectortile.Default_Tile_Layer_Version,
		Extent:  vectortile.Default_Tile_Layer_Extent,
	}

	for msg.Next() {
		switch msg.FieldNumber() {
		case 15: // version
			if v, err := msg.Uint32(); err != nil {
				return nil, err
			} else {
				layer.Version = v
			}
		case 1: // name
			if s, err := msg.String(); err != nil {
				return nil, err
			} else {
				layer.Name = s
			}
		case 2: // feature
			if data, err := msg.MessageData(); err != nil {
				return nil, err
			} else {
				d.features = append(d.features, data)
			}
		case 3: // keys
			if s, err := msg.String(); err != nil {
				return nil, err
			} else {
				d.keys = append(d.keys, s)
			}
		case 4: // values
			d.valMsg, err = msg.Message(d.valMsg)
			if err != nil {
				return nil, err
			}

			if v, err := decodeValueMsg(d.valMsg); err != nil {
				return nil, err
			} else {
				d.values = append(d.values, v)
			}
		case 5: // extent
			if e, err := msg.Uint32(); err != nil {
				return nil, err
			} else {
				layer.Extent = e
			}
		default:
			msg.Skip()
		}
	}

	if msg.Error() != nil {
		return nil, msg.Error()
	}

	layer.Features = make([]*geojson.Feature, len(d.features))
	for i, data := range d.features {
		msg.Reset(data)
		if f, err := d.Feature(msg); err != nil {
			return nil, err
		} else {
			layer.Features[i] = f
		}
	}

	return layer, nil
}

// geomDecoder holds state for geometry decoding.
type geomDecoder struct {
	iter  *pbr.Iterator
	count int
	used  int
	prev  geo.Point
}

func (gd *geomDecoder) NextPoint() (geo.Point, error) {
	gd.used += 2
	v, err := gd.iter.Uint32()
	if err != nil {
		return geo.Point{}, err
	}
	gd.prev[0] += unzigzag(v)

	v, err = gd.iter.Uint32()
	if err != nil {
		return geo.Point{}, err
	}
	gd.prev[1] += unzigzag(v)

	return gd.prev, nil
}

func (gd *geomDecoder) done() bool {
	return !gd.iter.HasNext()
}

func (gd *geomDecoder) cmdAndCount() (uint32, uint32, error) {
	if gd.done() {
		return 0, 0, errors.New("no more data")
	}

	v, err := gd.iter.Uint32()
	if err != nil {
		return 0, 0, err
	}
	gd.used++

	cmd, count := v&0x07, v>>3
	if cmd != closePath {
		if v := gd.used + int(2*count); gd.count < v {
			return 0, 0, fmt.Errorf("data cut short: needed %d, have %d", v, gd.count)
		}
	}

	return cmd, count, nil
}

func (gd *geomDecoder) decodeLine() (geo.LineString, error) {
	cmd, count, err := gd.cmdAndCount()
	if err != nil {
		return nil, err
	}

	if cmd != moveTo || count != 1 {
		return nil, errors.New("first command not one moveTo")
	}

	first, err := gd.NextPoint()
	if err != nil {
		return nil, err
	}

	cmd, count, err = gd.cmdAndCount()
	if err != nil {
		return nil, err
	} else if cmd != lineTo {
		return nil, errors.New("second command not a lineTo")
	}

	ls := make(geo.LineString, 0, count+1)
	ls = append(ls, first)
	for i := uint32(0); i < count; i++ {
		p, err := gd.NextPoint()
		if err != nil {
			return nil, err
		}
		ls = append(ls, p)
	}

	return ls, nil
}

func (gd *geomDecoder) decodeLineString() (geo.Geometry, error) {
	var mls geo.MultiLineString
	for !gd.done() {
		ls, err := gd.decodeLine()
		if err != nil {
			return nil, err
		}

		if gd.done() && len(mls) == 0 {
			return ls, nil
		}

		mls = append(mls, ls)
	}

	return mls, nil
}

func (gd *geomDecoder) decodePoint() (geo.Geometry, error) {
	_, count, err := gd.cmdAndCount()
	if err != nil {
		return nil, err
	}

	if count == 1 {
		return gd.NextPoint()
	}

	mp := make(geo.MultiPoint, 0, count)
	for i := uint32(0); i < count; i++ {
		p, err := gd.NextPoint()
		if err != nil {
			return nil, err
		}
		mp = append(mp, p)
	}

	return mp, nil
}

func (gd *geomDecoder) decodePolygon() (geo.Geometry, error) {
	var mp geo.MultiPolygon
	var p geo.Polygon
	for !gd.done() {
		ls, err := gd.decodeLine()
		if err != nil {
			return nil, err
		}

		r := geo.Ring(ls)
		cmd, _, err := gd.cmdAndCount()
		if err != nil {
			return nil, err
		} else if cmd == closePath && !r.Closed() {
			r = append(r, r[0])
		}

		// figure out if new polygon
		if len(mp) == 0 && len(p) == 0 {
			p = append(p, r)
		} else {
			if r.Orientation() == geo.CCW {
				mp = append(mp, p)
				p = geo.Polygon{r}
			} else {
				p = append(p, r)
			}
		}
	}

	if len(mp) == 0 {
		return p, nil
	}

	return append(mp, p), nil
}

func decodeValueMsg(msg *pbr.Message) (interface{}, error) {
	for msg.Next() {
		switch msg.FieldNumber() {
		case 1:
			return msg.String()
		case 2:
			v, err := msg.Float()
			return float64(v), err
		case 3:
			return msg.Double()
		case 4:
			v, err := msg.Int64()
			return float64(v), err
		case 5:
			v, err := msg.Uint64()
			return float64(v), err
		case 6:
			v, err := msg.Sint64()
			return float64(v), err
		case 7:
			return msg.Bool()
		default:
			msg.Skip()
		}
	}

	return nil, msg.Error()
}

func unmarshalTile(data []byte) (layers Layers, err error) {
	var m *pbr.Message
	msg := pbr.New(data)
	d := &decoder{}
	for msg.Next() {
		switch msg.FieldNumber() {
		case 3:
			m, err = msg.Message(m)
			if err != nil {
				return nil, err
			}

			layer, err := d.Layer(m)
			if err != nil {
				return nil, err
			}

			layers = append(layers, layer)
		default:
			msg.Skip()
		}
	}

	if msg.Error() != nil {
		return nil, msg.Error()
	}

	return layers, nil
}

// Check if data is GZipped by reading the "magic bytes"
// Rarely this method can result in false positives
func dataIsGZipped(data []byte) bool {
	return (data[0] == 0x1F && data[1] == 0x8B)
}
