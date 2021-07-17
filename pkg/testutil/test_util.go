package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/stretchr/testify/require"
)

const (
	testRoot = "/tmp"
	testDir  = "test-notes"
)

// WithNotes - test helper function
func WithNotes(t *testing.T, fn func(notes note.List)) {
	home, err := ioutil.TempDir(testRoot, testDir)
	require.NoError(t, err)

	defer os.RemoveAll(home)

	fn(note.List(home))
}
