package search

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/testutil"
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
		"a": "foo bar buzz",
		"b": "foobar bar addd buzz",
		"c": "fuzz bar xx buzz",
	}
}

func TestShoudSearch(t *testing.T) {
	files := m{
		"a": "foo bar buzz",
		"b": "foobar bar addd buzz",
		"c": "fuzz bar xx buzz",
		"d": "xx* yyy) abc",
	}
	testutil.WithNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			expected []*Result
			given    []string
		}{
			{"one term",
				[]*Result{
					{name: "a", lineNum: 1, text: "foo bar buzz", matches: []string{"foo"}},
					{name: "b", lineNum: 1, text: "foobar bar addd buzz", matches: []string{"foo"}},
				},
				[]string{"foo"}},
			{"with excluded terms",
				[]*Result{
					{name: "a", lineNum: 1, text: "foo bar buzz", matches: []string{"bar"}},
				},
				[]string{"bar", "-fuzz", "-foobar"},
			},
			{"with special regexp characters - star",
				[]*Result{
					{name: "d", lineNum: 1, text: "xx* yyy) abc", matches: []string{"xx*"}},
				},
				[]string{"xx*"},
			},
			{"with special regexp characters - parenthesis",
				[]*Result{
					{name: "d", lineNum: 1, text: "xx* yyy) abc", matches: []string{"yy", ")"}},
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
	testutil.WithNotes(mkTestFiles(), func(notes note.List) {
		underTest := NewEngine(notes)

		actual, err := underTest.Search("you will not find me")
		require.NoError(t, err)

		assert.Equal(t, 0, len(actual))
	})
}

func TestSearchShouldNotSearchInLastResutsFile(t *testing.T) {
	testutil.WithNotes(mkTestFiles(), func(notes note.List) {
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
					matches: []string{"foo"},
				},
				{
					name:    "b",
					lineNum: 1,
					text:    "foobar bar addd buzz",
					matches: []string{"foo"},
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
		"a": "abc\nfoo\nfoo bar",
		"b": "fuzz",
		"c": "bar foo",
	}
	testutil.WithNotes(files, func(notes note.List) {
		prepare := func() *Engine {
			s := NewEngine(notes)
			s.OnlyNames = true
			return s
		}
		t.Run("with simple search", func(t *testing.T) {
			expected := []*Result{
				{name: "b", lineNum: 1, text: "fuzz", matches: []string{"fuzz"}},
			}

			underTest := prepare()
			actual, err := underTest.Search("fuzz")
			require.NoError(t, err)

			assert.Equal(t, expected, actual)
		})
		t.Run("saved results should have line numbers of first occurrence", func(t *testing.T) {
			expected := []*Result{
				{name: "a", lineNum: 2, text: "foo", matches: []string{"foo"}},
				{name: "c", lineNum: 1, text: "bar foo", matches: []string{"foo"}},
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
		"foo":     "bar",
		"findme":  "abc",
		"findme2": "dfg",
		"buzz":    "findme",
	}
	testutil.WithNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			expected []*Result
			given    []string
		}{
			{"simple query",
				[]*Result{
					{name: "buzz", lineNum: 1, text: "findme", matches: []string{"findme"}},
					{name: "findme", lineNum: 1, text: " ", matches: []string{"findme"}},
					{name: "findme2", lineNum: 1, text: " ", matches: []string{"findme"}},
				},
				[]string{"findme"}},

			{"two terms",
				[]*Result{
					{name: "findme2", lineNum: 1, text: " ", matches: []string{"findme2"}},
					{name: "foo", lineNum: 1, text: " ", matches: []string{"fo"}},
				},
				[]string{"findme2", "fo"}},
			{"with terms and excluded terms",
				[]*Result{
					{name: "buzz", lineNum: 1, text: "findme", matches: []string{"find"}},
					{name: "findme", lineNum: 1, text: " ", matches: []string{"find"}},
				},
				[]string{"find", "-findme2"}},
			{"terms and exclude terms are case insensitive",
				[]*Result{
					{name: "buzz", lineNum: 1, text: "findme", matches: []string{"find"}},
					{name: "findme", lineNum: 1, text: " ", matches: []string{"find"}},
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
	testutil.WithNotes(files, func(notes note.List) {
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
		"a": "foo",
		"b": "fuzz",
		"c": "bar\nfoo",
	}
	testutil.WithNotes(files, func(notes note.List) {
		expected := []*Result{
			&Result{name: "a", lineNum: 1, text: "foo", matches: []string{"foo"}},
			&Result{name: "c", lineNum: 2, text: "foo", matches: []string{"foo"}},
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
		"a": "abc\nfoo\nfOo bar",
		"b": "fuzz",
		"c": "bar FOO",
	}
	testutil.WithNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			expected []*Result
			given    []string
		}{
			{"simple query",
				[]*Result{
					{name: "a", lineNum: 2, text: "foo", matches: []string{"foo"}},
				},
				[]string{"foo"}},
			{"simple query-2",
				[]*Result{
					{name: "a", lineNum: 3, text: "fOo bar", matches: []string{"fOo"}},
					{name: "c", lineNum: 1, text: "bar FOO", matches: []string{"FOO"}},
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

func toSortedLines(s string) []string {
	lines := strings.Split(strings.Trim(s, "\n"), "\n")
	sort.Strings(lines)
	return lines
}

func mustReadLastResults(t *testing.T, notes note.List) string {
	resultsFilename := filepath.Join(notes.HomeDir(), note.DotNotelogDir, lastResultsFile)
	body, err := os.ReadFile(resultsFilename)
	require.NoError(t, err) // file must exist

	return string(body)
}
