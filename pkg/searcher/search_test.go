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
	withNotes(files, func(notes note.List) {
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
				underTest := NewSearcher(notes, actual)
				// FIXME
				_, err := underTest.Search(tt.given...)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, len(toSortedLines(actual.String())))
			})
		}
	})
}

func TestSearchShoudWriteLastSearchResults(t *testing.T) {
	withNotes(files, func(notes note.List) {
		actual := &bytes.Buffer{}

		underTest := NewSearcher(notes, actual)
		underTest.SaveResults = true

		_, err := underTest.Search("foobar")
		require.NoError(t, err)

		resultsFilename := filepath.Join(notes.HomeDir(), note.DotNotelogDir, lastResultsFile)
		body, err := ioutil.ReadFile(resultsFilename)

		expected := []string{
			notes.HomeDir() + "/b/main.org:1:foobar bar addd buzz",
		}
		require.NoError(t, err) // file must exist
		assert.Equal(t, expected, toSortedLines(string(body)))
	})
}

func TestSearchShoudWriteLastSearchResultsWithoutTermColor(t *testing.T) {
	withNotes(files, func(notes note.List) {
		actual := &bytes.Buffer{}

		underTest := NewSearcher(notes, actual)
		underTest.SaveResults = true

		n, err := underTest.Search("foo bar")
		require.NoError(t, err)
		require.Equal(t, 1, n)

		resultsFilename := filepath.Join(notes.HomeDir(), note.DotNotelogDir, lastResultsFile)
		body, err := ioutil.ReadFile(resultsFilename)

		expected := notes.HomeDir() + "/a/main.org:1:foo bar buzz\n"
		require.NoError(t, err) // file must exist
		assert.Equal(t, expected, string(body))
	})
}

func TestSearchShoudReturnZeroResultsIfFoundNothing(t *testing.T) {
	withNotes(files, func(notes note.List) {
		actual := &bytes.Buffer{}
		underTest := NewSearcher(notes, actual)

		n, err := underTest.Search("you will not find me")
		require.NoError(t, err)

		assert.Equal(t, 0, n)
	})
}

func TestSearchShouldNotSearchInLastResutsFile(t *testing.T) {
	withNotes(files, func(notes note.List) {
		// search 2 times so that last_results will be filled
		for i := 0; i < 2; i++ {
			out := &bytes.Buffer{}
			underTest := NewSearcher(notes, out)
			underTest.SaveResults = true

			n, err := underTest.Search("foo")
			require.NoError(t, err)

			expected := []string{
				"a:1:foo bar buzz",
				"b:1:foobar bar addd buzz",
			}
			assert.Equal(t, len(expected), n)
			assert.Equal(t, expected, toSortedLines(out.String()))
		}
	})
}

func TestSearcShouldSearchNamesOnlyIfSet(t *testing.T) {
	files := map[string]string{
		"a/main.org": "abc\nfoo\nfoo bar",
		"b/main.org": "fuzz",
		"c/main.org": "bar foo",
	}
	withNotes(files, func(notes note.List) {
		actual := &bytes.Buffer{}
		underTest := NewSearcher(notes, actual)
		underTest.OnlyNames = true

		n, err := underTest.Search("foo")
		require.NoError(t, err)

		expected := []string{
			"a:1: ",
			"c:1: ",
		}
		assert.Equal(t, len(expected), n)
		assert.Equal(t, expected, toSortedLines(actual.String()))
	})
}

func TestSearchShouldSearchInNoteNames(t *testing.T) {
	files := map[string]string{
		"foo/main.org":     "bar",
		"findme/main.org":  "abc",
		"findme2/main.org": "dfg",
		"buzz/main.org":    "findme",
	}
	withNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			expected []string
			given    []string
		}{
			{"simple query",
				[]string{
					"buzz:1:findme",
					"findme2:1: ",
					"findme:1: ",
				},
				[]string{"findme"}},
			{"two terms",
				[]string{
					"findme2:1: ",
					"foo:1: ",
				},
				[]string{"findme2", "fo"}},
			{"with terms and excluded terms",
				[]string{
					"buzz:1:findme",
					"findme:1: ",
				},
				[]string{"find", "-findme2"}},
			{"terms and exclude terms are case insensitive",
				[]string{
					"buzz:1:findme",
					"findme:1: ",
				},
				[]string{"finD", "-FindmE2"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				out := &bytes.Buffer{}
				underTest := NewSearcher(notes, out)

				n, err := underTest.Search(tt.given...)
				require.NoError(t, err)

				assert.Equal(t, len(tt.expected), n)
				assert.Equal(t, tt.expected, toSortedLines(out.String()))
			})
		}
	})
}

func disabledTestSearchShouldLookInArchive(t *testing.T) {
	// TODO: enable
	files := map[string]string{
		".archive/andme/main.org": "abc d",
		"findme/main.org":         "abc",
		"foo/main.org":            "bar",
	}
	withNotes(files, func(notes note.List) {
		out := &bytes.Buffer{}
		underTest := NewSearcher(notes, out)

		n, err := underTest.Search("abc")
		require.NoError(t, err)

		expected := []string{
			fmt.Sprintf("%s/.archive/andme/main.org:1:abc d", notes.HomeDir()),
			fmt.Sprintf("%s/findme/main.org:1:abc", notes.HomeDir()),
		}
		assert.Equal(t, len(expected), n)
		assert.Equal(t, expected, toSortedLines(out.String()))
	})
}

func TestSearchNoteNamesOnlyShouldEnsureNoTermColorsInOutput(t *testing.T) {
	// TODO: remove if no term colors will be used

	files := map[string]string{
		"a/main.org":   "foo",
		"b/main.org":   "fuzz",
		"foo/main.org": "bar\nbuzz",
	}
	withNotes(files, func(notes note.List) {
		out := &bytes.Buffer{}
		underTest := NewSearcher(notes, out)
		underTest.OnlyNames = true

		n, err := underTest.Search("foo", "buzz")
		require.NoError(t, err)

		expected := []string{
			"a:1: ",
			"foo:1: ",
		}

		assert.Equal(t, len(expected), n)
		assert.Equal(t, expected, toSortedLines(out.String()))
	})
}

func TestSearcShouldSearchCaseSensitiveIfSet(t *testing.T) {
	files := map[string]string{
		"a/main.org": "abc\nfoo\nfOo bar",
		"b/main.org": "fuzz",
		"c/main.org": "bar FOO",
	}
	withNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			expected []string
			given    []string
		}{
			{"simple query",
				[]string{
					"a:2:foo",
				},
				[]string{"foo"}},
			{"simple query-2",
				[]string{
					"a:3:fOo bar",
					"c:1:bar FOO",
				},
				[]string{"FOO", "fOo"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual := &bytes.Buffer{}
				underTest := NewSearcher(notes, actual)
				underTest.CaseSensitive = true

				n, err := underTest.Search(tt.given...)
				require.NoError(t, err)

				assert.Equal(t, len(tt.expected), n)
				assert.Equal(t, tt.expected, toSortedLines(actual.String()))
			})
		}
	})
}

func withNotes(files m, fn func(notes note.List)) {
	homeDir, err := ioutil.TempDir("", "test_notes")
	if err != nil {
		panic(err)
	}

	must(os.MkdirAll(homeDir, 0755))
	defer os.RemoveAll(homeDir)

	must(os.MkdirAll(filepath.Join(homeDir, note.DotNotelogDir), 0755))

	for name, body := range files {
		fullName := filepath.Join(homeDir, name)
		dir, _ := filepath.Split(fullName)
		must(os.MkdirAll(dir, 0755))
		must(ioutil.WriteFile(fullName, []byte(body), mode))
	}

	fn(note.List(homeDir))
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
