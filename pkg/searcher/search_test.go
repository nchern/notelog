package searcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	homeDir = "/tmp/test_notes"

	mode = 0644
)

var (
	files = map[string]string{
		"a.txt": "foo bar buzz",
		"b.txt": "foobar bar addd buzz",
		"c.txt": "fuzz bar xx buzz",
	}
)

func init() {
	must(os.Setenv("NOTELOG_GREP", defaultGrep)) // make sure we always use defaultGrep in tests
}

func TestShoudSearch(t *testing.T) {
	withFiles(func() {
		var tests = []struct {
			name     string
			expected int
			given    []string
		}{
			{"one term",
				2, []string{"foo"}},
			{"with excluded terms",
				1, []string{"bar", "-fuzz", "-foobar"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual := &bytes.Buffer{}
				underTest := NewSearcher(note.List(homeDir), actual)

				require.NoError(t, underTest.Search(tt.given...))
				assert.Equal(t, tt.expected, len(toLines(actual.String())))
			})
		}
	})
}

func TestSearchShoudWriteLastSearchResults(t *testing.T) {
	withFiles(func() {
		actual := &bytes.Buffer{}

		underTest := NewSearcher(note.List(homeDir), actual)
		underTest.SaveResults = true

		require.NoError(t, underTest.Search("foobar"))

		resultsFilename := filepath.Join(homeDir, note.DotNotelogDir, lastResultsFile)
		body, err := ioutil.ReadFile(resultsFilename)

		require.NoError(t, err) // file must exist
		assert.Equal(t, actual.Bytes(), body)
	})
}

func TestSearchShoudWriteLastSearchResultsWithoutTermColor(t *testing.T) {
	withFiles(func() {
		actual := &bytes.Buffer{}

		underTest := NewSearcher(note.List(homeDir), actual)
		underTest.SaveResults = true
		underTest.grepCmd = "grep -E --colour=always"

		require.NoError(t, underTest.Search("foo bar"))

		resultsFilename := filepath.Join(homeDir, note.DotNotelogDir, lastResultsFile)
		body, err := ioutil.ReadFile(resultsFilename)

		expected := "/tmp/test_notes/a.txt:1:foo bar buzz\n"
		require.NoError(t, err) // file must exist
		assert.Equal(t, expected, string(body))
	})
}

func TestSearchShoudReturnExitErrorOneIfFoundNothing(t *testing.T) {
	withFiles(func() {
		actual := &bytes.Buffer{}
		underTest := NewSearcher(note.List(homeDir), actual)

		err := underTest.Search("you will not find me")

		assert.NotNil(t, err)
		assert.Equal(t, 1, (err.(*exec.ExitError)).ExitCode())
	})
}

func TestSearchShouldCorrectlyHandleCommandOverride(t *testing.T) {
	n := note.List(homeDir)
	actual := &bytes.Buffer{}
	underTest := NewSearcher(n, actual)
	underTest.grepCmd = "echo --bar" // use echo to get args as output

	err := underTest.Search("foo")

	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("--bar -rni (foo) %s\n", n.HomeDir()), actual.String())
}

func TestSearchShouldNotGetResultsFromLastResutsFile(t *testing.T) {
	// withFiles(func() {
	// 	// search 2 times so that last_results will be filled
	// 	for i := 0; i < 2; i++ {
	// 		actual := &bytes.Buffer{}
	// 		underTest := NewSearcher(note.List(homeDir), actual)
	// 		underTest.SaveResults = true

	// 		expected := "/tmp/test_notes/b.txt:1:foobar bar addd buzz\n/tmp/test_notes/a.txt:1:foo bar buzz\n"

	// 		require.NoError(t, underTest.Search("foo"))
	// 		assert.Equal(t, expected, actual.String())
	// 	}
	// })
}

func withFiles(fn func()) {
	must(os.MkdirAll(homeDir, 0755))
	defer os.RemoveAll(homeDir)

	must(os.MkdirAll(filepath.Join(homeDir, note.DotNotelogDir), 0755))

	for name, body := range files {
		must(ioutil.WriteFile(filepath.Join(homeDir, name), []byte(body), mode))
	}

	fn()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func toLines(s string) []string {
	return strings.Split(strings.Trim(s, "\n"), "\n")
}
