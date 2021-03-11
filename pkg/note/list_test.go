package note

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRoot = "/tmp"
	testDir  = "test-notes"
)

func withNotes(t *testing.T, fn func(notes List)) {
	home, err := ioutil.TempDir(testRoot, testDir)
	require.NoError(t, err)

	fn(List(home))
}

func makeNote(t *testing.T, notes List, name string) *Note {
	nt, err := notes.Add(name)
	require.NoError(t, err)
	found, _ := nt.Exists()
	require.True(t, found)
	return nt
}

func TestRemove(t *testing.T) {
	withNotes(t, func(notes List) {
		underTest := makeNote(t, notes, "foo")

		err := notes.Remove("foo")
		assert.NoError(t, err)

		found, _ := underTest.Exists()
		assert.False(t, found)
	})
}

func TestRename(t *testing.T) {
	withNotes(t, func(notes List) {
		underTest := makeNote(t, notes, "foo")

		err := notes.Rename("foo", "bar")
		assert.NoError(t, err)

		found, _ := underTest.Exists()
		assert.False(t, found)

		found, _ = notes.Note("bar").Exists()
		assert.True(t, found)
	})
}

func TestArchive(t *testing.T) {
	withNotes(t, func(notes List) {
		underTest := makeNote(t, notes, "foo")

		err := notes.Archive("foo")
		assert.NoError(t, err)

		found, _ := underTest.Exists()
		assert.False(t, found)

		_, err = os.Stat(filepath.Join(notes.HomeDir(), archiveNoteDir, underTest.name, defaultFilename))
		assert.NoError(t, err)
	})
}
