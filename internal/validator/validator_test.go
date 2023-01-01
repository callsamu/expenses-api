package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInListOfValues(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		list     []string
		expected bool
	}{
		{
			name:     "Is in list",
			value:    "foo",
			list:     []string{"foo", "bar"},
			expected: true,
		},
		{
			name:     "Is not in list",
			value:    "foobar",
			list:     []string{"foo", "bar"},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualValues(t, tc.expected, In(tc.value, tc.list...))
		})
	}
}

func TestAllValuesInListUnique(t *testing.T) {
	cases := []struct {
		name     string
		list     []string
		expected bool
	}{
		{
			name:     "List has duplicates",
			list:     []string{"ha", "ha", "he"},
			expected: false,
		},
		{
			name:     "List has no duplicates",
			list:     []string{"ha", "he", "hi"},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualValues(t, tc.expected, Unique(tc.list...))
		})
	}
}
