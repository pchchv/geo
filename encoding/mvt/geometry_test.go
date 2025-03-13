package mvt

import (
	"fmt"
	"reflect"
	"testing"
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
			t.Fatalf("failed to get value: %v", err)
		}

		value := decodeValue(kve.Values[i])
		if value != "[1,2,3]" {
			t.Errorf("should encode non standard types as json")
		}
	})

	t.Run("nil value", func(t *testing.T) {
		i, err := kve.Value(nil)
		if err != nil {
			t.Fatalf("failed to get value: %v", err)
		}

		value := decodeValue(kve.Values[i])
		if value != "null" {
			t.Errorf("should encode null as json")
		}
	})
}
