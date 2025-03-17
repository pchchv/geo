package mvt

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pchchv/geo/encoding/mvt/vectortile"
	"github.com/pchchv/pbr"
)

type stringer int

func (s stringer) String() string {
	return fmt.Sprintf("%d", s)
}

func TestKeyValueEncoder_JSON(t *testing.T) {
	kve := newKeyValueEncoder()
	t.Run("non comparable value", func(t *testing.T) {
		i, err := kve.Value([]int{1, 2, 3})
		if err != nil {
			t.Fatalf("failed to get value: %e", err)
		}

		value := decodeValue(kve.Values[i])
		if value != "[1,2,3]" {
			t.Errorf("should encode non standard types as json")
		}
	})

	t.Run("nil value", func(t *testing.T) {
		i, err := kve.Value(nil)
		if err != nil {
			t.Fatalf("failed to get value: %e", err)
		}

		value := decodeValue(kve.Values[i])
		if value != "null" {
			t.Errorf("should encode null as json")
		}
	})
}

func TestEncodeValue(t *testing.T) {
	cases := []struct {
		name   string
		input  interface{}
		output interface{}
	}{
		{
			name:   "string",
			input:  "abc",
			output: "abc",
		},
		{
			name:   "stringer",
			input:  stringer(10),
			output: "10",
		},
		{
			name:   "int",
			input:  int(1),
			output: float64(1),
		},
		{
			name:   "int8",
			input:  int8(2),
			output: float64(2),
		},
		{
			name:   "int16",
			input:  int16(3),
			output: float64(3),
		},
		{
			name:   "int32",
			input:  int32(4),
			output: float64(4),
		},
		{
			name:   "int64",
			input:  int64(5),
			output: float64(5),
		},
		{
			name:   "uint",
			input:  int(1),
			output: float64(1),
		},
		{
			name:   "uint8",
			input:  int8(2),
			output: float64(2),
		},
		{
			name:   "uint16",
			input:  int16(3),
			output: float64(3),
		},
		{
			name:   "uint32",
			input:  int32(4),
			output: float64(4),
		},
		{
			name:   "uint64",
			input:  int64(5),
			output: float64(5),
		},
		{
			name:   "float32",
			input:  float32(6),
			output: float64(6),
		},
		{
			name:   "float64",
			input:  float64(7),
			output: float64(7),
		},
		{
			name:   "bool",
			input:  true,
			output: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := encodeValue(tc.input)
			if err != nil {
				t.Fatalf("encode failure: %e", err)
			}

			result := decodeValue(val)
			if !reflect.DeepEqual(result, tc.output) {
				t.Errorf("incorrect value: %[1]T != %[2]T, %[1]v != %[2]v", result, tc.output)
			}
		})
	}

	// error if a weird type, but typical json decode result
	input := map[string]interface{}{
		"a": 1,
		"b": 2,
	}

	if _, err := encodeValue(input); err == nil {
		t.Errorf("expecting error: %e", err)
	}
}

func sliceToIterator(vals []uint32) *pbr.Iterator {
	feature := &vectortile.Tile_Feature{
		Geometry: vals,
	}

	data, err := feature.Marshal()
	if err != nil {
		panic(err)
	}

	msg := pbr.New(data)
	for msg.Next() {
		switch msg.FieldNumber() {
		case 4:
			iter, err := msg.Iterator(nil)
			if err != nil {
				panic(err)
			}

			return iter
		default:
			msg.Skip()
		}
	}

	panic("unreachable")
}
