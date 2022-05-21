package search

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastNthResult(t *testing.T) {
	testutil.WithNotes(mkTestFiles(), func(notes note.List) {
		underTest := NewEngine(notes)

		// perform search to generate last results file
		actual, err := underTest.Search("foo")
		require.NoError(t, err)
		require.Equal(t, 2, len(actual))

		buf := &bytes.Buffer{}
		r, err := NewPersistentRenderer(
			notes,
			&StreamRenderer{W: buf})
		assert.NoError(t, Render(r, actual, false))

		b, err := ioutil.ReadFile(notes.MetadataFilename(lastResultsFile))
		require.NoError(t, err)
		persistedResults := strings.Split(string(b), "\n")

		var tests = []struct {
			name     string
			expected string
			given    int
		}{
			{"should return first result",
				persistedResults[0], 1},
			{"should return second result",
				persistedResults[1], 2},
			{"should return empty on out of bounds",
				"", 100500},
			{"should return empty on negative",
				"", -1},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual, err := GetLastNthResult(notes, tt.given)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}

func TestGetLastNthResultShouldReturnEmptyStringIfNoResultsFound(t *testing.T) {
	testutil.WithNotes(mkTestFiles(), func(notes note.List) {
		// make sure last results file does not exist
		_, err := os.Stat(notes.MetadataFilename(lastResultsFile))
		assert.True(t, os.IsNotExist(err))

		actual, err := GetLastNthResult(notes, 1)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})
}
