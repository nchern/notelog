package search

import (
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldGrep(t *testing.T) {
	files := m{
		"a": "foo bar buz",
		"b": "foobar bar addd buzz\nabc dfgh",
		"c": "fuzz Bar xx buzz",
		"d": "buzz bar xx",
	}
	testutil.WithNotes(files, func(notes note.List) {
		var tests = []struct {
			name          string
			caseSensitive bool
			expected      []*Result
			given         string
		}{
			{"trivial case",
				false,
				[]*Result{
					{name: "a", lineNum: 1, text: "foo bar buz", matches: []string{"foo"}},
					{name: "b", lineNum: 1, text: "foobar bar addd buzz", matches: []string{"foo"}},
				},
				"foo"},
			{"regexp - beginning of the line",
				false,
				[]*Result{
					{name: "d", lineNum: 1, text: "buzz bar xx", matches: []string{"buz"}},
				},
				"^buz"},

			{"regexp - or",
				false,
				[]*Result{
					{name: "b", lineNum: 2, text: "abc dfgh", matches: []string{"abc"}},
					{name: "c", lineNum: 1, text: "fuzz Bar xx buzz", matches: []string{"fuz"}},
				},
				"(abc|fuz)"},
			{"no results",
				false,
				[]*Result{},
				"aaa"},
			{"regexp - case sensitive",
				true,
				[]*Result{
					{name: "c", lineNum: 1, text: "fuzz Bar xx buzz", matches: []string{"Bar xx buzz"}},
				},
				"Ba.+"},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				underTest := NewGrepEngine(notes)
				underTest.CaseSensitive = tt.caseSensitive

				actual, err := underTest.Search(tt.given)

				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}

func TestGrepShouldFailOnBrokenRegexp(t *testing.T) {
	testutil.WithNotes(mkTestFiles(), func(notes note.List) {
		underTest := NewGrepEngine(notes)
		actual, err := underTest.Search("(a")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}
