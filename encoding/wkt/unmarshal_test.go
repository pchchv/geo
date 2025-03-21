package wkt

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalPoint_errors(t *testing.T) {
	cases := []struct {
		name string
		s    string
		err  error
	}{
		{
			name: "just name",
			s:    "POINT",
			err:  ErrNotWKT,
		},
		{
			name: "too many points",
			s:    "POINT(1.34 2.35 3.36)",
			err:  ErrNotWKT,
		},
		{
			name: "not a point",
			s:    "MULTIPOINT((1.34 2.35))",
			err:  ErrIncorrectGeometry,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := UnmarshalPoint(tc.s); err != tc.err {
				t.Fatalf("incorrect error: %e != %e", err, tc.err)
			}
		})
	}
}

func TestUnmarshalMultiPoint_errors(t *testing.T) {
	cases := []struct {
		name string
		s    string
		err  error
	}{
		{
			name: "just name",
			s:    "MULTIPOINT",
			err:  ErrNotWKT,
		},
		{
			name: "too many points",
			s:    "MULTIPOINT((1 2),(3 4 5))",
			err:  ErrNotWKT,
		},
		{
			name: "not a multipoint",
			s:    "POINT(1 2)",
			err:  ErrIncorrectGeometry,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := UnmarshalMultiPoint(tc.s); err != tc.err {
				t.Fatalf("incorrect error: %e != %e", err, tc.err)
			}
		})
	}
}

func TestUnmarshalLineString_errors(t *testing.T) {
	cases := []struct {
		name string
		s    string
		err  error
	}{
		{
			name: "just name",
			s:    "LINESTRING",
			err:  ErrNotWKT,
		},
		{
			name: "too many points",
			s:    "LINESTRING(1 2,3 4 5)",
			err:  ErrNotWKT,
		},
		{
			name: "not a multipoint",
			s:    "POINT(1 2)",
			err:  ErrIncorrectGeometry,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := UnmarshalLineString(tc.s); err != tc.err {
				t.Fatalf("incorrect error: %e != %e", err, tc.err)
			}
		})
	}
}

func TestUnmarshalMultiLineString_errors(t *testing.T) {
	cases := []struct {
		name string
		s    string
		err  error
	}{
		{
			name: "just name",
			s:    "MULTILINESTRING",
			err:  ErrNotWKT,
		},
		{
			name: "too many points",
			s:    "MULTILINESTRING((1 2,3 4 5))",
			err:  ErrNotWKT,
		},
		{
			name: "not a multi linestring",
			s:    "POINT(1 2)",
			err:  ErrIncorrectGeometry,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := UnmarshalMultiLineString(tc.s); err != tc.err {
				t.Fatalf("incorrect error: %e != %e", err, tc.err)
			}
		})
	}
}

func TestTrimSpaceBrackets(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		expected string
	}{
		{
			name:     "empty string",
			s:        "",
			expected: "",
		},
		{
			name:     "blank string",
			s:        "   ",
			expected: "",
		},
		{
			name:     "single point",
			s:        "(1 2)",
			expected: "1 2",
		},
		{
			name:     "double brackets",
			s:        "((1 2),(0.5 1.5))",
			expected: "(1 2),(0.5 1.5)",
		},
		{
			name:     "multiple values",
			s:        "(1 2,0.5 1.5)",
			expected: "1 2,0.5 1.5",
		},
		{
			name:     "multiple points",
			s:        "((1 2,3 4),(5 6,7 8))",
			expected: "(1 2,3 4),(5 6,7 8)",
		},
		{
			name:     "triple brackets",
			s:        "(((1 2,3 4)),((5 6,7 8),(1 2,5 4)))",
			expected: "((1 2,3 4)),((5 6,7 8),(1 2,5 4))",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if v, err := trimSpaceBrackets(tc.s); err != nil {
				t.Fatalf("unexpected error: %e", err)
			} else if v != tc.expected {
				t.Log(trimSpaceBrackets(tc.s))
				t.Log(tc.expected)
				t.Errorf("trim space and brackets error")
			}
		})
	}
}

func TestTrimSpaceBrackets_errors(t *testing.T) {
	cases := []struct {
		name string
		s    string
		err  error
	}{
		{
			name: "no brackets",
			s:    "1 2",
			err:  ErrNotWKT,
		},
		{
			name: "no start bracket",
			s:    "1 2)",
			err:  ErrNotWKT,
		},
		{
			name: "no end bracket",
			s:    "(1 2",
			err:  ErrNotWKT,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := trimSpaceBrackets(tc.s); err != tc.err {
				t.Fatalf("wrong error: %e != %e", err, tc.err)
			}
		})
	}
}

func TestSplitOnComma(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "comma",
			input:    "0 1,3 0,4 3,0 4,0 1",
			expected: []string{"0 1", "3 0", "4 3", "0 4", "0 1"},
		},
		{
			name:     "comma spaces",
			input:    "0 1 ,3 0, 4 3 , 0 4  ,   0 1",
			expected: []string{"0 1", "3 0", "4 3", "0 4", "0 1"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var results []string
			if err := splitOnComma(tc.input, func(s string) error {
				results = append(results, s)
				return nil
			}); err != nil {
				t.Fatalf("impossible error: %e", err)
			}

			if !reflect.DeepEqual(tc.expected, results) {
				t.Log(tc.input)

				data, _ := json.Marshal(results)
				t.Log(string(data))

				t.Log(tc.expected)
				t.Errorf("incorrect results")
			}

		})
	}
}
