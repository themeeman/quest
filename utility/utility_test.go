package utility

import (
	"testing"
)

type testContainsResult struct {
	found bool
	index int
}

type testContainsData struct {
	slice    interface{}
	value    interface{}
	expected testContainsResult
}

func TestContains(t *testing.T) {
	testCases := []testContainsData{
		{
			slice: []int{1, 5, 3, 2, 1},
			value: 5,
			expected: testContainsResult{
				found: true,
				index: 1,
			},
		},
		{
			slice: []int{1, 5, 3, 2, 1},
			value: 1,
			expected: testContainsResult{
				found: true,
				index: 0,
			},
		},
		{
			slice: []int{1, 5, 3, 2, 1},
			value: 255,
			expected: testContainsResult{
				found: false,
				index: -1,
			},
		},
		{
			slice: []int{1, 5, 3, 2, 1},
			value: "5",
			expected: testContainsResult{
				found: false,
				index: -1,
			},
		},
		{
			slice: []int{1, 5, 3, 2, 1},
			value: float64(2.5),
			expected: testContainsResult{
				found: false,
				index: -1,
			},
		},
	}

	for _, test := range testCases {
		found, index := Contains(test.slice, test.value)
		result := testContainsResult{
			found: found,
			index: index,
		}

		if result != test.expected {
			t.Errorf("Test failed with slice=%v and value=%#v; expected %v, got %+v",
				test.slice, test.value, test.expected, result)
		} else {
			t.Logf("Test succeeded with slice=%v and value=%#v; got %+v",
				test.slice, test.value, test.expected)
		}
	}
}
