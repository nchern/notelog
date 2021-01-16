package searcher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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
	files = m{
		"a/main.org": "foo bar buzz",
		"b/main.org": "foobar bar addd buzz",
		"c/main.org": "fuzz bar xx buzz",
	}
)

type m map[string]string

func init() {
	must(os.Setenv("NOTELOG_GREP", defaultGrep)) // make sure we always use defaultGrep in tests
}

func TestShoudSearch(t *testing.T) {
	withFiles(files, func() {
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
	withFiles(files, func() {
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
	withFiles(files, func() {
		actual := &bytes.Buffer{}

		underTest := NewSearcher(note.List(homeDir), actual)
		underTest.SaveResults = true
		underTest.grepCmd = "grep -E --colour=always"

		require.NoError(t, underTest.Search("foo bar"))

		resultsFilename := filepath.Join(homeDir, note.DotNotelogDir, lastResultsFile)
		body, err := ioutil.ReadFile(resultsFilename)

		expected := "/tmp/test_notes/a/main.org:1:foo bar buzz\n"
		require.NoError(t, err) // file must exist
		assert.Equal(t, expected, string(body))
	})
}

func TestSearchShoudReturnExitErrorOneIfFoundNothing(t *testing.T) {
	withFiles(files, func() {
		actual := &bytes.Buffer{}
		underTest := NewSearcher(note.List(homeDir), actual)

		err := underTest.Search("you will not find me")
		require.NotNil(t, err)

		assert.Equal(t, ErrNoResults, err)
	})
}

func TestSearchShouldCorrectlyHandleCommandOverride(t *testing.T) {
	withFiles(files, func() { // search requires NOTELOG_HOME to exist
		n := note.List(homeDir)
		actual := &bytes.Buffer{}

		underTest := NewSearcher(n, actual)
		underTest.grepCmd = "echo --bar" // use echo to get args as output

		err := underTest.Search("foo")
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf("--bar -rni (foo) %s\n", n.HomeDir()), actual.String())
	})
}

func TestSearchShouldNotGetResultsFromLastResutsFile(t *testing.T) {
	withFiles(files, func() {
		// search 2 times so that last_results will be filled
		for i := 0; i < 2; i++ {
			out := &bytes.Buffer{}
			underTest := NewSearcher(note.List(homeDir), out)
			underTest.SaveResults = true

			err := underTest.Search("foo")
			require.NoError(t, err)

			expected := []string{
				"/tmp/test_notes/a/main.org:1:foo bar buzz",
				"/tmp/test_notes/b/main.org:1:foobar bar addd buzz",
			}
			actual := toLines(out.String())
			sort.Strings(actual)

			assert.Equal(t, expected, actual)
		}
	})
}

func TestSearcShouldSearchNamesOnlyIfSet(t *testing.T) {
	notes := map[string]string{
		"a/main.org": "abc\nfoo\nfoo bar",
		"b/main.org": "fuzz",
		"c/main.org": "bar foo",
	}
	withFiles(notes, func() {
		actual := &bytes.Buffer{}
		underTest := NewSearcher(note.List(homeDir), actual)
		underTest.OnlyNames = true

		err := underTest.Search("foo")
		require.NoError(t, err)

		expected := []string{
			"/tmp/test_notes/a/main.org:1:",
			"/tmp/test_notes/c/main.org:1:",
		}

		actualLines := toLines(actual.String())
		sort.Strings(actualLines)

		assert.Equal(t, expected, actualLines)
	})
}

func TestSearchShouldSearchInNoteNames(t *testing.T) {
	notes := map[string]string{
		"foo/main.org":     "bar",
		"findme/main.org":  "abc",
		"findme2/main.org": "dfg",
		"buzz/main.org":    "findme",
	}
	withFiles(notes, func() {
		var tests = []struct {
			name     string
			expected []string
			given    []string
		}{
			{"simple query",
				[]string{
					"/tmp/test_notes/buzz/main.org:1:findme",
					"/tmp/test_notes/findme/main.org:1",
					"/tmp/test_notes/findme2/main.org:1",
				},
				[]string{"findme"}},
			{"two terms",
				[]string{
					"/tmp/test_notes/findme2/main.org:1",
					"/tmp/test_notes/foo/main.org:1",
				},
				[]string{"findme2", "fo"}},
			{"with terms and excluded terms",
				[]string{
					"/tmp/test_notes/buzz/main.org:1:findme",
					"/tmp/test_notes/findme/main.org:1",
				},
				[]string{"find", "-findme2"}},
			{"terms and exclude terms are case insensitive",
				[]string{
					"/tmp/test_notes/buzz/main.org:1:findme",
					"/tmp/test_notes/findme/main.org:1",
				},
				[]string{"finD", "-FindmE2"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				out := &bytes.Buffer{}
				underTest := NewSearcher(note.List(homeDir), out)

				err := underTest.Search(tt.given...)
				require.NoError(t, err)

				actual := toLines(out.String())
				sort.Strings(actual)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}

func TestSearchNoteNamesOnlyShouldEnsureNoTermColorsInOutput(t *testing.T) {
	// this test requires sift command to present in the system
	if err := exec.Command("sift", "--version").Run(); errors.Is(err, exec.ErrNotFound) {
		t.Skip("can not run this test: sift is not installed")
	}

	notes := map[string]string{
		"a/main.org":   "foo",
		"b/main.org":   "fuzz",
		"foo/main.org": "bar\nbuzz",
	}
	withFiles(notes, func() {
		out := &bytes.Buffer{}
		underTest := NewSearcher(note.List(homeDir), out)
		underTest.OnlyNames = true
		underTest.grepCmd = "sift --color"

		err := underTest.Search("foo", "buzz")
		require.NoError(t, err)

		expected := []string{
			"/tmp/test_notes/a/main.org:1:",
			"/tmp/test_notes/foo/main.org:1:",
		}

		actual := toLines(out.String())
		sort.Strings(actual)
		assert.Equal(t, expected, actual)
	})
}

func withFiles(files m, fn func()) {
	must(os.MkdirAll(homeDir, 0755))
	defer os.RemoveAll(homeDir)

	must(os.MkdirAll(filepath.Join(homeDir, note.DotNotelogDir), 0755))

	for name, body := range files {
		fullName := filepath.Join(homeDir, name)
		dir, _ := filepath.Split(fullName)
		must(os.MkdirAll(dir, 0755))
		must(ioutil.WriteFile(fullName, []byte(body), mode))
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
