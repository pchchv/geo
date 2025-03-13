package mvt

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pchchv/geo/encoding/mvt/vectortile"
)

const closePath = 7

type keyValueEncoder struct {
	Keys          []string
	keyMap        map[string]uint32
	Values        []*vectortile.Tile_Value
	valueMap      map[interface{}]uint32
	keySortBuffer []string
}

func newKeyValueEncoder() *keyValueEncoder {
	return &keyValueEncoder{
		keyMap:   make(map[string]uint32),
		valueMap: make(map[interface{}]uint32),
	}
}

func (kve *keyValueEncoder) Key(s string) (i uint32) {
	if i, ok := kve.keyMap[s]; !ok {
		i = uint32(len(kve.Keys))
		kve.Keys = append(kve.Keys, s)
		kve.keyMap[s] = i
	}

	return
}

func (kve *keyValueEncoder) Value(v interface{}) (uint32, error) {
	// If a type is not comparable can't figure out uniqueness in the hash,
	// we also can't encode it into a vectortile.Tile_Value.
	// So we encoded it as a json string, which is what other encoders also do.
	if v == nil || !reflect.TypeOf(v).Comparable() {
		if data, err := json.Marshal(v); err != nil {
			return 0, fmt.Errorf("uncomparable: %T", v)
		} else {
			v = string(data)
		}
	}

	if i, ok := kve.valueMap[v]; ok {
		return i, nil
	}

	tv, err := encodeValue(v)
	if err != nil {
		return 0, err
	}

	i := uint32(len(kve.Values))
	kve.Values = append(kve.Values, tv)
	kve.valueMap[v] = i

	return i, nil
}

type geomEncoder struct {
	prevX int32
	prevY int32
	Data  []uint32
}

func newGeomEncoder(l int) *geomEncoder {
	return &geomEncoder{
		Data: make([]uint32, 0, l),
	}
}

func (ge *geomEncoder) ClosePath() {
	ge.Data = append(ge.Data, (1<<3)|closePath)
}

func encodeValue(v interface{}) (*vectortile.Tile_Value, error) {
	tv := &vectortile.Tile_Value{}
	switch t := v.(type) {
	case string:
		tv.StringValue = &t
	case fmt.Stringer:
		s := t.String()
		tv.StringValue = &s
	case int:
		i := int64(t)
		tv.SintValue = &i
	case int8:
		i := int64(t)
		tv.SintValue = &i
	case int16:
		i := int64(t)
		tv.SintValue = &i
	case int32:
		i := int64(t)
		tv.SintValue = &i
	case int64:
		i := int64(t)
		tv.SintValue = &i
	case uint:
		i := uint64(t)
		tv.UintValue = &i
	case uint8:
		i := uint64(t)
		tv.UintValue = &i
	case uint16:
		i := uint64(t)
		tv.UintValue = &i
	case uint32:
		i := uint64(t)
		tv.UintValue = &i
	case uint64:
		i := uint64(t)
		tv.UintValue = &i
	case float32:
		tv.FloatValue = &t
	case float64:
		tv.DoubleValue = &t
	case bool:
		tv.BoolValue = &t
	default:
		return nil, fmt.Errorf("unable to encode value of type %T: %v", v, v)
	}

	return tv, nil
}
