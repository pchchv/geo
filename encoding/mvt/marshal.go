package mvt

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"sort"
	"strconv"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/mvt/vectortile"
	"github.com/pchchv/geo/geojson"
	"google.golang.org/protobuf/proto"
)

// Marshal encodes a set of layers into a Mapbox Vector Tile format.
// Features that have a nil geometry,
// for some reason, will be skipped and not included.
func Marshal(layers Layers) ([]byte, error) {
	vt := &vectortile.Tile{
		Layers: make([]*vectortile.Tile_Layer, 0, len(layers)),
	}

	for _, l := range layers {
		v, e := l.Version, l.Extent
		kve := newKeyValueEncoder()
		layer := &vectortile.Tile_Layer{
			Name:     &l.Name,
			Version:  &v,
			Extent:   &e,
			Features: make([]*vectortile.Tile_Feature, 0, len(l.Features)),
		}

		for _, f := range l.Features {
			if err := addFeature(layer, kve, f); err != nil {
				return nil, err
			}
		}

		layer.Keys = kve.Keys
		layer.Values = kve.Values
		vt.Layers = append(vt.Layers, layer)
	}

	return proto.Marshal(vt)
}

// MarshalGzipped marshal the layers into Mapbox Vector Tile format and gzip the result.
// Often MVT data is already gzipped to a file, such as an mbtiles file.
func MarshalGzipped(layers Layers) ([]byte, error) {
	data, err := Marshal(layers)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	gzwriter := gzip.NewWriter(buf)
	if _, err = gzwriter.Write(data); err != nil {
		return nil, err
	}

	if err = gzwriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func convertIntID(i int) *uint64 {
	if i < 0 {
		return nil
	}

	v := uint64(i)
	return &v
}

func convertID(id interface{}) *uint64 {
	switch id := id.(type) {
	case nil:
		return nil
	case int:
		return convertIntID(id)
	case int8:
		return convertIntID(int(id))
	case int16:
		return convertIntID(int(id))
	case int32:
		return convertIntID(int(id))
	case int64:
		return convertIntID(int(id))
	case uint:
		v := uint64(id)
		return &v
	case uint8:
		v := uint64(id)
		return &v
	case uint16:
		v := uint64(id)
		return &v
	case uint32:
		v := uint64(id)
		return &v
	case uint64:
		v := uint64(id)
		return &v
	case float32:
		return convertIntID(int(id))
	case float64:
		return convertIntID(int(id))
	case string:
		if i, err := strconv.Atoi(id); err == nil {
			return convertIntID(i)
		}
	}

	return nil
}

func encodeProperties(kve *keyValueEncoder, properties geojson.Properties) ([]uint32, error) {
	tags := make([]uint32, 0, 2*len(properties))
	kve.keySortBuffer = kve.keySortBuffer[:0]
	for k := range properties {
		kve.keySortBuffer = append(kve.keySortBuffer, k)
	}
	sort.Strings(kve.keySortBuffer)

	for _, k := range kve.keySortBuffer {
		ki := kve.Key(k)
		if vi, err := kve.Value(properties[k]); err != nil {
			return nil, fmt.Errorf("property %s: %v", k, err)
		} else {
			tags = append(tags, ki, vi)
		}
	}

	return tags, nil
}

func addFeature(layer *vectortile.Tile_Layer, kve *keyValueEncoder, f *geojson.Feature) error {
	if f.Geometry == nil {
		return nil
	}

	if f.Geometry.GeoJSONType() == "GeometryCollection" {
		for _, g := range f.Geometry.(geo.Collection) {
			return addSingleGeometryFeature(layer, kve, g, f.Properties, f.ID)
		}
	}

	return addSingleGeometryFeature(layer, kve, f.Geometry, f.Properties, f.ID)
}

func addSingleGeometryFeature(layer *vectortile.Tile_Layer, kve *keyValueEncoder, g geo.Geometry, p geojson.Properties, id interface{}) error {
	geomType, encodedGeometry, err := encodeGeometry(g)
	if err != nil {
		return fmt.Errorf("error encoding geometry: %v : %s", g, err.Error())
	}

	tags, err := encodeProperties(kve, p)
	if err != nil {
		return fmt.Errorf("error encoding geometry: %v : %s", g, err.Error())
	}

	layer.Features = append(layer.Features, &vectortile.Tile_Feature{
		Id:       convertID(id),
		Tags:     tags,
		Type:     &geomType,
		Geometry: encodedGeometry,
	})

	return nil
}
