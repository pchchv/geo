package mvt

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/mvt/vectortile"
)

const (
	moveTo    = 1
	lineTo    = 2
	closePath = 7
)

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

func (kve *keyValueEncoder) Key(s string) uint32 {
	if i, ok := kve.keyMap[s]; ok {
		return i
	}

	i := uint32(len(kve.Keys))
	kve.Keys = append(kve.Keys, s)
	kve.keyMap[s] = i

	return i
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

func (ge *geomEncoder) MoveTo(points []geo.Point) {
	l := uint32(len(points))
	ge.Data = append(ge.Data, (l<<3)|moveTo)
	ge.addPoints(points)
}

func (ge *geomEncoder) LineTo(points []geo.Point) {
	l := uint32(len(points))
	ge.Data = append(ge.Data, (l<<3)|lineTo)
	ge.addPoints(points)
}

func (ge *geomEncoder) addPoints(points []geo.Point) {
	for i := range points {
		x := int32(points[i][0]) - ge.prevX
		y := int32(points[i][1]) - ge.prevY
		ge.prevX = int32(points[i][0])
		ge.prevY = int32(points[i][1])
		ge.Data = append(ge.Data,
			uint32((x<<1)^(x>>31)),
			uint32((y<<1)^(y>>31)),
		)
	}
}

func encodeValue(v interface{}) (*vectortile.Tile_Value, error) {
	tv := &vectortile.Tile_Value{}
	switch t := v.(type) {
	case string:
		tv.Value = &vectortile.Tile_Value_StringValue{StringValue: t}
	case fmt.Stringer:
		s := t.String()
		tv.Value = &vectortile.Tile_Value_StringValue{StringValue: s}
	case int:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case int8:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case int16:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case int32:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case int64:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case uint:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case uint8:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case uint16:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case uint32:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case uint64:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case float32:
		tv.Value = &vectortile.Tile_Value_FloatValue{FloatValue: float32(t)}
	case float64:
		tv.Value = &vectortile.Tile_Value_DoubleValue{DoubleValue: t}
	case bool:
		tv.Value = &vectortile.Tile_Value_BoolValue{BoolValue: t}
	default:
		return nil, fmt.Errorf("unable to encode value of type %T: %v", v, v)
	}

	return tv, nil
}

func elMLS(mls geo.MultiLineString) (c int) {
	for _, ls := range mls {
		c += 2 + 2*len(ls)
	}
	return
}

func elP(p geo.Polygon) (c int) {
	for _, r := range p {
		c += 3 + 2*len(r)
	}
	return
}

func elMP(mp geo.MultiPolygon) (c int) {
	for _, p := range mp {
		c += elP(p)
	}
	return
}

func encodeGeometry(g geo.Geometry) (vectortile.Tile_GeomType, []uint32, error) {
	switch g := g.(type) {
	case geo.Point:
		e := newGeomEncoder(3)
		e.MoveTo([]geo.Point{g})
		return vectortile.Tile_POINT, e.Data, nil
	case geo.MultiPoint:
		e := newGeomEncoder(1 + 2*len(g))
		e.MoveTo([]geo.Point(g))
		return vectortile.Tile_POINT, e.Data, nil
	case geo.LineString:
		e := newGeomEncoder(2 + 2*len(g))
		e.MoveTo([]geo.Point{g[0]})
		e.LineTo([]geo.Point(g[1:]))
		return vectortile.Tile_LINESTRING, e.Data, nil
	case geo.MultiLineString:
		e := newGeomEncoder(elMLS(g))
		for _, ls := range g {
			e.MoveTo([]geo.Point{ls[0]})
			e.LineTo([]geo.Point(ls[1:]))
		}
		return vectortile.Tile_LINESTRING, e.Data, nil
	case geo.Ring:
		e := newGeomEncoder(3 + 2*len(g))
		e.MoveTo([]geo.Point{g[0]})
		if g.Closed() {
			e.LineTo([]geo.Point(g[1 : len(g)-1]))
		} else {
			e.LineTo([]geo.Point(g[1:]))
		}
		e.ClosePath()
		return vectortile.Tile_POLYGON, e.Data, nil
	case geo.Polygon:
		e := newGeomEncoder(elP(g))
		for _, r := range g {
			e.MoveTo([]geo.Point{r[0]})
			if r.Closed() {
				e.LineTo([]geo.Point(r[1 : len(r)-1]))
			} else {
				e.LineTo([]geo.Point(r[1:]))
			}
			e.ClosePath()
		}
		return vectortile.Tile_POLYGON, e.Data, nil
	case geo.MultiPolygon:
		e := newGeomEncoder(elMP(g))
		for _, p := range g {
			for _, r := range p {
				e.MoveTo([]geo.Point{r[0]})
				if r.Closed() {
					e.LineTo([]geo.Point(r[1 : len(r)-1]))
				} else {
					e.LineTo([]geo.Point(r[1:]))
				}
				e.ClosePath()
			}
		}
		return vectortile.Tile_POLYGON, e.Data, nil
	case geo.Collection:
		return 0, nil, errors.New("geometry collections are not supported")
	case geo.Bound:
		return encodeGeometry(g.ToPolygon())
	default:
		panic(fmt.Sprintf("geometry type not supported: %T", g))
	}
}

func unzigzag(v uint32) float64 {
	return float64(int32(((v >> 1) & ((1 << 32) - 1)) ^ -(v & 1)))
}

func decodeValue(v *vectortile.Tile_Value) interface{} {
	if v == nil || v.Value == nil {
		return nil
	}

	switch value := v.Value.(type) {
	case *vectortile.Tile_Value_StringValue:
		return value.StringValue
	case *vectortile.Tile_Value_FloatValue:
		return float64(value.FloatValue)
	case *vectortile.Tile_Value_DoubleValue:
		return value.DoubleValue
	case *vectortile.Tile_Value_IntValue:
		return value.IntValue
	case *vectortile.Tile_Value_UintValue:
		return value.UintValue
	case *vectortile.Tile_Value_SintValue:
		return value.SintValue
	case *vectortile.Tile_Value_BoolValue:
		return value.BoolValue
	default:
		return nil
	}
}
