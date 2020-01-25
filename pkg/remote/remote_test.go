package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryPush(t *testing.T) {
	var tests = []struct {
		name        string
		expected    []string
		expectedErr error
		given       string
	}{
		{"rsync", []string{"rsync", "-r", "src", "localhost:foo"}, nil, "rsync"},
		//		{"git", "git", "git"},
		{"unknown scheme", nil, errUnknownScheme, "bla"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{Scheme: tt.given, Addr: "localhost:foo"}
			args, err := e.Push("src")
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, args)
		})
	}
}
