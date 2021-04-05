package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNoteName(t *testing.T) {
	var tests = []struct {
		name     string
		given    string
		expected error
	}{
		{"should not be empty", "", errEmptyName},
		{"should not start with a dot", ".dotname", errNameStartsWithDot},
		{"should accept simple alpha-numeric name", "alpha-numeric_123", nil},
		{"should accept simple small and capital letters", "fooBar", nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := validateNoteName(tt.given)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestValidateNoteNameShouldNOTContainSymbols(t *testing.T) {
	forbiddenSymbols := ".,<>:/\\'; +?~!@#$%^&*[]()"
	for _, c := range forbiddenSymbols {
		t.Run("should not accept "+string(c), func(t *testing.T) {
			actual := validateNoteName("note" + string(c))
			assert.Equal(t, errNameRegexNoMatch, actual)
		})
	}
}
