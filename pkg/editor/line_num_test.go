package editor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineNumShouldHandle(t *testing.T) {
	var tests = []struct {
		name        string
		expected    int64
		expectedErr error
		given       string
	}{
		{"empty string", 0, nil, ""},
		{"correct number", 42, nil, "42"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			underTest := LineNumber(tt.given)

			actual, err := underTest.ToInt()
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLineNumShouldHandleIncorrectNumbers(t *testing.T) {
	underTest := LineNumber("abc")

	actual, err := underTest.ToInt()
	assert.Error(t, err)
	assert.Equal(t, int64(0), actual)
}
