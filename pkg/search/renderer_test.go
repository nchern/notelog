package search

import (
	"bytes"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistentRendererShoudWriteLastSearchResults(t *testing.T) {
	files := m{
		"a": "foo bar buzz\nbbb foo aaa",
		"b": "foobar bar addd buzz",
		"c": "fuzz bar xx buzz",
	}
	testutil.WithNotes(files, func(notes note.List) {
		underTest := NewEngine(notes)

		res, err := underTest.Search("foo")
		require.NoError(t, err)

		buf := &bytes.Buffer{}
		r, err := NewPersistentRenderer(notes, &StreamRenderer{W: buf})

		assert.NoError(t, Render(r, res, false))

		expected := []string{
			"a:1:",
			"a:2:",
			"b:1:",
		}
		actual := mustReadLastResults(t, notes)
		assert.Equal(t, expected, toSortedLines(actual))
	})
}

func TestPersistentRendererShoudWriteLastSearchResultsInArchivedNotes(t *testing.T) {
	files := m{
		"a/main.org": "foo bar buzz\nbbb foo aaa",
		"b/main.org": "foobar bar addd buzz",
	}
	withArchivedNotes(files, func(notes note.List) {
		underTest := NewEngine(notes)

		res, err := underTest.Search("foo")
		require.NoError(t, err)

		buf := &bytes.Buffer{}
		r, err := NewPersistentRenderer(notes, &StreamRenderer{W: buf})

		assert.NoError(t, Render(r, res, false))

		expected := []string{
			"a:1:a",
			"a:2:a",
			"b:1:a",
		}
		actual := mustReadLastResults(t, notes)
		assert.Equal(t, expected, toSortedLines(actual))
	})
}

func TestSearchShoudWriteLastSearchResultsWithoutTermColor(t *testing.T) {
	testutil.WithNotes(mkTestFiles(), func(notes note.List) {
		underTest := NewEngine(notes)

		res, err := underTest.Search("foo bar")
		require.NoError(t, err)
		require.Equal(t, 1, len(res))

		buf := &bytes.Buffer{}
		r, err := NewPersistentRenderer(notes, &StreamRenderer{W: buf})

		assert.NoError(t, Render(r, res, false))

		expected := "a:1:\n"
		actual := mustReadLastResults(t, notes)
		assert.Equal(t, expected, string(actual))
	})
}
