package server

import "testing"

func TestContainsStr(t *testing.T) {
	testCases := []struct {
		number     int32
		slice    []string
		str      string
		expected bool
	}{
		{number: 1, slice: []string{"a", "b"}, str: "a", expected: true},
		{number: 2, slice: []string{"a", "b"}, str: "d", expected: false},
		{number: 3, slice: nil, str: "a", expected: false},
		{number: 4, slice: []string{}, str: "a", expected: false},
		{number: 5, slice: []string{"a", "b"}, str: "", expected: false},
	}
	for _, testCase := range testCases {
		if ContainsStr(testCase.slice, testCase.str) != testCase.expected {
			t.Errorf("Test %d error", testCase.number)
		}
	}
}
