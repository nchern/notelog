package note

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDir = "test-notes"
)

func withNotes(t *testing.T, fn func(notes List)) {
	// we can't use testutil.WithNotes as testutil depends on this pkg
	home := filepath.Join(t.TempDir(), testDir)
	must(os.MkdirAll(home, defaultDirPerms))

	fn(List(home))
}

func makeNote(t *testing.T, notes List, name string) *Note {
	nt, err := notes.GetOrCreate(name, Org)
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

		found, err := underTest.Exists()
		require.NoError(t, err)
		assert.False(t, found)

		bar, err := notes.Get("bar")
		require.NoError(t, err)
		found, err = bar.Exists()
		require.True(t, found)
	})
}

func TestArchive(t *testing.T) {
	const expectedFilename = "main.org"

	withNotes(t, func(notes List) {
		underTest := makeNote(t, notes, "foo")

		err := notes.Archive("foo")
		assert.NoError(t, err)

		found, _ := underTest.Exists()
		assert.False(t, found)

		_, err = os.Stat(filepath.Join(notes.HomeDir(), archiveNoteDir, underTest.name, expectedFilename))
		assert.NoError(t, err)
	})
}

func TestInit(t *testing.T) {

	withNotes(t, func(notes List) {
		assert.NoError(t, notes.Init())
		_, err := os.Stat(notes.metadataRoot())
		assert.NoError(t, err)

		_, err = os.Stat(filepath.Join(notes.HomeDir(), archiveNoteDir))
		assert.NoError(t, err)

		// any subsequent calls to init should also be successful
		for i := 0; i < 11; i++ {
			assert.NoError(t, notes.Init())
		}
	})
}

func TestGetOrCreate(t *testing.T) {
	const noteName = "new-one"

	withNotes(t, func(notes List) {
		underTest, err := notes.GetOrCreate(noteName, Org)
		assert.NoError(t, err)

		found, err := underTest.Exists()
		assert.NoError(t, err)
		assert.True(t, found)

		underTest2, err := notes.GetOrCreate(noteName, Org)
		assert.NoError(t, err)
		assert.Equal(t, underTest.FullPath(), underTest2.FullPath())
	})
}

func TestGet(t *testing.T) {
	// merge with GetOrCreate?
	const noteName = "new-one"

	withNotes(t, func(notes List) {
		missing, err := notes.Get("nonexistent")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotExist))
		assert.Nil(t, missing)

		created := makeNote(t, notes, noteName)

		underTest, err := notes.Get(noteName)
		assert.NoError(t, err)
		assert.Equal(t, created.FullPath(), underTest.FullPath())
	})
}

func TestAllShouldNotFailOnNonExistedNotes(t *testing.T) {
	withNotes(t, func(notes List) {
		res, err := notes.All()
		assert.Empty(t, res)
		assert.NoError(t, err)

		notePath := filepath.Join(notes.HomeDir(), "broken")
		must(os.MkdirAll(notePath, defaultDirPerms))

		notePath = filepath.Join(notes.HomeDir(), "file")
		_, err = os.Create(notePath)
		require.NoError(t, err)

		res, err = notes.All()
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestCopy(t *testing.T) {
	const (
		noteFoo = "foo"
		newName = "foo-copy"
		body    = "this is the note to copy"
	)
	actualBuf := &bytes.Buffer{}
	withNotes(t, func(underTest List) {
		nt := makeNote(t, underTest, noteFoo)
		w, err := nt.writer()
		require.NoError(t, err)
		_, err = io.WriteString(w, body)
		require.NoError(t, err)
		w.Close()

		err = underTest.Copy(noteFoo, newName)
		assert.NoError(t, err)

		newOne, err := underTest.Get(newName)
		require.NoError(t, err)
		assert.Equal(t, newName, newOne.name)

		assert.NoError(t, newOne.Dump(actualBuf))
		assert.Equal(t, body, actualBuf.String())

		err = underTest.Copy("nonexistent", newName)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotExist))
	})
}
