package searcher

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/stretchr/testify/assert"
)

func TestGetLastNthResult(t *testing.T) {
	withFiles(func() {
		n := note.List(homeDir)
		buf := &bytes.Buffer{}
		underTest := NewSearcher(n, buf)
		underTest.SaveResults = true

		// perform search to generate last results file
		err := underTest.Search("foo")
		assert.NoError(t, err)

		var tests = []struct {
			name     string
			expected string
			given    int
		}{
			{"should return first result",
				filepath.Join(homeDir, "b.txt") + ":1:foobar bar addd buzz", 1},
			{"should return second result",
				filepath.Join(homeDir, "a.txt") + ":1:foo bar buzz", 2},
			{"should return empty on out of bounds",
				"", 100500},
			{"should return empty on negative",
				"", -1},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual, err := GetLastNthResult(n, tt.given)
				fmt.Println(actual)
				assert.NoError(t, err)
				assert.Equal(t, actual, tt.expected)
			})
		}
	})
}

func TestGetLastNthResultShouldReturnEmptyStringIfNoResultsFound(t *testing.T) {
	withFiles(func() {
		n := note.List(homeDir)
		// make sure last results file does not exist
		_, err := os.Stat(n.MetadataFilename(lastResultsFile))
		assert.True(t, os.IsNotExist(err))

		actual, err := GetLastNthResult(n, 1)

		assert.NoError(t, err)
		assert.Empty(t, actual)
	})
}
