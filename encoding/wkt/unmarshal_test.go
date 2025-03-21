package wkt

import (
	"encoding/json"
	"reflect"
	"testing"
)

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
