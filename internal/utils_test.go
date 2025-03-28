package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolToString(t *testing.T) {
	tests := []struct {
		input    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, test := range tests {
		t.Run("Testing BoolToString", func(t *testing.T) {
			result := BoolToString(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
