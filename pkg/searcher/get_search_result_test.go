package searcher

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastNthResult(t *testing.T) {
	withFiles(files, func() {
		n := note.List(homeDir)
		buf := &bytes.Buffer{}

		underTest := NewSearcher(n, buf)
		underTest.SaveResults = true

		// perform search to generate last results file
		err := underTest.Search("foo")
		require.NoError(t, err)

		var tests = []struct {
			name     string
			expected string
			given    int
		}{
			{"should return first result",
				filepath.Join(homeDir, "a/main.org") + ":1:foo bar buzz", 1},
			{"should return second result",
				filepath.Join(homeDir, "b/main.org") + ":1:foobar bar addd buzz", 2},
			{"should return empty on out of bounds",
				"", 100500},
			{"should return empty on negative",
				"", -1},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual, err := GetLastNthResult(n, tt.given)
				assert.NoError(t, err)
				assert.Equal(t, actual, tt.expected)
			})
		}
	})
}

func TestGetLastNthResultShouldReturnEmptyStringIfNoResultsFound(t *testing.T) {
	withFiles(files, func() {
		n := note.List(homeDir)
		// make sure last results file does not exist
		_, err := os.Stat(n.MetadataFilename(lastResultsFile))
		assert.True(t, os.IsNotExist(err))

		actual, err := GetLastNthResult(n, 1)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})
}
