package search

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

var ()

type m map[string]string

func mkTestFiles() m {
	return m{
		"a/main.org": "foo bar buzz",
		"b/main.org": "foobar bar addd buzz",
		"c/main.org": "fuzz bar xx buzz",
	}
}

func TestShoudSearch(t *testing.T) {
	files := m{
		"a/main.org": "foo bar buzz",
		"b/main.org": "foobar bar addd buzz",
		"c/main.org": "fuzz bar xx buzz",
		"d/main.org": "xx* yyy) abc",
	}
	withNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			expected []*Result
			given    []string
		}{
			{"one term",
				[]*Result{
					{name: "a", lineNum: 1, text: "foo bar buzz"},
					{name: "b", lineNum: 1, text: "foobar bar addd buzz"},
				},
				[]string{"foo"}},
			{"with excluded terms",
				[]*Result{
					{name: "a", lineNum: 1, text: "foo bar buzz"},
				},
				[]string{"bar", "-fuzz", "-foobar"},
			},
			{"with special regexp characters - star",
				[]*Result{
					{name: "d", lineNum: 1, text: "xx* yyy) abc"},
				},
				[]string{"xx*"},
			},
			{"with special regexp characters - parenthesis",
				[]*Result{
					{name: "d", lineNum: 1, text: "xx* yyy) abc"},
				},
				[]string{"yy", ")"},
			},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				underTest := NewEngine(notes)
				actual, err := underTest.Search(tt.given...)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}

func TestSearchShoudReturnZeroResultsIfFoundNothing(t *testing.T) {
	withNotes(mkTestFiles(), func(notes note.List) {
		underTest := NewEngine(notes)

		actual, err := underTest.Search("you will not find me")
		require.NoError(t, err)

		assert.Equal(t, 0, len(actual))
	})
}

func TestSearchShouldNotSearchInLastResutsFile(t *testing.T) {
	withNotes(mkTestFiles(), func(notes note.List) {
		// search 2 times so that last_results will be filled
		for i := 0; i < 2; i++ {
			underTest := NewEngine(notes)

			actual, err := underTest.Search("foo")
			require.NoError(t, err)

			expected := []*Result{
				{
					name:    "a",
					lineNum: 1,
					text:    "foo bar buzz",
				},
				{
					name:    "b",
					lineNum: 1,
					text:    "foobar bar addd buzz",
				},
			}
			assert.Equal(t, len(expected), len(actual))
			assert.Equal(t, expected, actual)

			out := &bytes.Buffer{}
			r, err := NewPersistentRenderer(notes, &StreamRenderer{W: out})
			assert.NoError(t, Render(r, actual, false))
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
		prepare := func() *Engine {
			s := NewEngine(notes)
			s.OnlyNames = true
			return s
		}
		t.Run("with simple search", func(t *testing.T) {
			expected := []*Result{
				{name: "b", lineNum: 1, text: "fuzz"},
			}

			underTest := prepare()
			actual, err := underTest.Search("fuzz")
			require.NoError(t, err)

			assert.Equal(t, expected, actual)
		})
		t.Run("saved results should have line numbers of first occurrence", func(t *testing.T) {
			expected := []*Result{
				{name: "a", lineNum: 2, text: "foo"},
				{name: "c", lineNum: 1, text: "bar foo"},
			}

			underTest := prepare()

			actual, err := underTest.Search("foo")
			require.NoError(t, err)
			require.Equal(t, len(expected), len(actual))

			assert.Equal(t, expected, actual)
		})
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
			expected []*Result
			given    []string
		}{
			{"simple query",
				[]*Result{
					{name: "buzz", lineNum: 1, text: "findme"},
					{name: "findme", lineNum: 1, text: " "},
					{name: "findme2", lineNum: 1, text: " "},
				},
				[]string{"findme"}},
			{"two terms",
				[]*Result{
					{name: "findme2", lineNum: 1, text: " "},
					{name: "foo", lineNum: 1, text: " "},
				},
				[]string{"findme2", "fo"}},
			{"with terms and excluded terms",
				[]*Result{
					{name: "buzz", lineNum: 1, text: "findme"},
					{name: "findme", lineNum: 1, text: " "},
				},
				[]string{"find", "-findme2"}},
			{"terms and exclude terms are case insensitive",
				[]*Result{
					{name: "buzz", lineNum: 1, text: "findme"},
					{name: "findme", lineNum: 1, text: " "},
				},
				[]string{"finD", "-FindmE2"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				underTest := NewEngine(notes)

				actual, err := underTest.Search(tt.given...)
				require.NoError(t, err)

				assert.Equal(t, len(tt.expected), len(actual))
				assert.Equal(t, tt.expected, actual)
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
		underTest := NewEngine(notes)

		actual, err := underTest.Search("abc")
		require.NoError(t, err)

		expected := []string{
			fmt.Sprintf("%s/.archive/andme/main.org:1:abc d", notes.HomeDir()),
			fmt.Sprintf("%s/findme/main.org:1:abc", notes.HomeDir()),
		}
		out := &bytes.Buffer{}
		assert.Equal(t, len(expected), len(actual))
		assert.Equal(t, expected, toSortedLines(out.String()))
	})
}

func TestSearchSearchInNotesOfDifferentTypes(t *testing.T) {

	files := map[string]string{
		"a/main.org": "foo",
		"b/main.org": "fuzz",
		"c/main.md":  "bar\nfoo",
	}
	withNotes(files, func(notes note.List) {
		expected := []*Result{
			&Result{name: "a", lineNum: 1, text: "foo"},
			&Result{name: "c", lineNum: 2, text: "foo"},
		}
		underTest := NewEngine(notes)

		actual, err := underTest.Search("foo")
		require.NoError(t, err)

		assert.Equal(t, len(expected), len(actual))
		assert.Equal(t, expected, actual)
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
			expected []*Result
			given    []string
		}{
			{"simple query",
				[]*Result{
					{name: "a", lineNum: 2, text: "foo"},
				},
				[]string{"foo"}},
			{"simple query-2",
				[]*Result{
					{name: "a", lineNum: 3, text: "fOo bar"},
					{name: "c", lineNum: 1, text: "bar FOO"},
				},
				[]string{"FOO", "fOo"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				underTest := NewEngine(notes)
				underTest.CaseSensitive = true

				actual, err := underTest.Search(tt.given...)
				require.NoError(t, err)

				assert.Equal(t, len(tt.expected), len(actual))
				assert.Equal(t, tt.expected, actual)
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

func mustReadLastResults(t *testing.T, notes note.List) string {
	resultsFilename := filepath.Join(notes.HomeDir(), note.DotNotelogDir, lastResultsFile)
	body, err := ioutil.ReadFile(resultsFilename)
	require.NoError(t, err) // file must exist

	return string(body)
}
