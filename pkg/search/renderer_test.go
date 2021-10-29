package search

import (
	"bytes"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistentRendererShoudWriteLastSearchResults(t *testing.T) {
	files := m{
		"a/main.org": "foo bar buzz\nbbb foo aaa",
		"b/main.org": "foobar bar addd buzz",
		"c/main.org": "fuzz bar xx buzz",
	}
	withNotes(files, func(notes note.List) {
		underTest := NewEngine(notes)

		res, err := underTest.Search("foo")
		require.NoError(t, err)

		buf := &bytes.Buffer{}
		r, err := NewPersistentRenderer(notes, &StreamRenderer{W: buf})

		assert.NoError(t, Render(r, res, false))

		expected := []string{
			"a:1",
			"a:2",
			"b:1",
		}
		actual := mustReadLastResults(t, notes)
		assert.Equal(t, expected, toSortedLines(actual))
	})
}

func TestSearchShoudWriteLastSearchResultsWithoutTermColor(t *testing.T) {
	withNotes(mkTestFiles(), func(notes note.List) {
		underTest := NewEngine(notes)

		res, err := underTest.Search("foo bar")
		require.NoError(t, err)
		require.Equal(t, 1, len(res))

		buf := &bytes.Buffer{}
		r, err := NewPersistentRenderer(notes, &StreamRenderer{W: buf})

		assert.NoError(t, Render(r, res, false))

		expected := "a:1\n"
		actual := mustReadLastResults(t, notes)
		assert.Equal(t, expected, string(actual))
	})
}
