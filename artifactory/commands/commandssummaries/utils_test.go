package commandssummaries

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceProtocol(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Replace https with http",
			input:    "https://example.com",
			expected: "http://example.com",
		},
		{
			name:     "Input already has http",
			input:    "http://example.com",
			expected: "http://example.com",
		},
		{
			name:     "Input is not a URL",
			input:    "Just a string",
			expected: "Just a string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := replaceProtocol(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
