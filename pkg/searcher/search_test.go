package searcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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
				assert.Equal(t, tt.expected, len(toSortedLines(actual.String())))
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

			assert.Equal(t, expected, toSortedLines(out.String()))
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
			"/tmp/test_notes/a/main.org:1: ",
			"/tmp/test_notes/c/main.org:1: ",
		}

		assert.Equal(t, expected, toSortedLines(actual.String()))
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
					"/tmp/test_notes/findme/main.org:1: ",
					"/tmp/test_notes/findme2/main.org:1: ",
				},
				[]string{"findme"}},
			{"two terms",
				[]string{
					"/tmp/test_notes/findme2/main.org:1: ",
					"/tmp/test_notes/foo/main.org:1: ",
				},
				[]string{"findme2", "fo"}},
			{"with terms and excluded terms",
				[]string{
					"/tmp/test_notes/buzz/main.org:1:findme",
					"/tmp/test_notes/findme/main.org:1: ",
				},
				[]string{"find", "-findme2"}},
			{"terms and exclude terms are case insensitive",
				[]string{
					"/tmp/test_notes/buzz/main.org:1:findme",
					"/tmp/test_notes/findme/main.org:1: ",
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

				assert.Equal(t, tt.expected, toSortedLines(out.String()))
			})
		}
	})
}

func disTestSearchShouldLookInArchive(t *testing.T) {
	// TODO: enable
	notes := map[string]string{
		".archive/andme/main.org": "abc d",
		"findme/main.org":         "abc",
		"foo/main.org":            "bar",
	}
	withFiles(notes, func() {
		out := &bytes.Buffer{}
		underTest := NewSearcher(note.List(homeDir), out)

		err := underTest.Search("abc")
		require.NoError(t, err)

		expected := []string{
			fmt.Sprintf("%s/.archive/andme/main.org:1:abc d", homeDir),
			fmt.Sprintf("%s/findme/main.org:1:abc", homeDir),
		}
		assert.Equal(t, expected, toSortedLines(out.String()))
	})
}

func TestSearchNoteNamesOnlyShouldEnsureNoTermColorsInOutput(t *testing.T) {
	// TODO: remove if no term colors will be used

	notes := map[string]string{
		"a/main.org":   "foo",
		"b/main.org":   "fuzz",
		"foo/main.org": "bar\nbuzz",
	}
	withFiles(notes, func() {
		out := &bytes.Buffer{}
		underTest := NewSearcher(note.List(homeDir), out)
		underTest.OnlyNames = true

		err := underTest.Search("foo", "buzz")
		require.NoError(t, err)

		expected := []string{
			"/tmp/test_notes/a/main.org:1: ",
			"/tmp/test_notes/foo/main.org:1: ",
		}

		assert.Equal(t, expected, toSortedLines(out.String()))
	})
}

func TestSearcShouldSearchCaseSensitiveIfSet(t *testing.T) {
	notes := map[string]string{
		"a/main.org": "abc\nfoo\nfOo bar",
		"b/main.org": "fuzz",
		"c/main.org": "bar FOO",
	}
	withFiles(notes, func() {
		var tests = []struct {
			name     string
			expected []string
			given    []string
		}{
			{"simple query",
				[]string{
					"/tmp/test_notes/a/main.org:2:foo",
				},
				[]string{"foo"}},
			{"simple query-2",
				[]string{
					"/tmp/test_notes/a/main.org:3:fOo bar",
					"/tmp/test_notes/c/main.org:1:bar FOO",
				},
				[]string{"FOO", "fOo"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual := &bytes.Buffer{}
				underTest := NewSearcher(note.List(homeDir), actual)
				underTest.CaseSensitive = true

				err := underTest.Search(tt.given...)
				require.NoError(t, err)

				assert.Equal(t, tt.expected, toSortedLines(actual.String()))
			})
		}
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

func toSortedLines(s string) []string {
	lines := strings.Split(strings.Trim(s, "\n"), "\n")
	sort.Strings(lines)
	return lines
}
