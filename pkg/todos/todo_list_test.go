package todos

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTODORegexNotDone(t *testing.T) {
	var tests = []struct {
		name     string
		expected bool
		given    string
	}{
		{"undone", true, "   - [] foo"},
		{"undone with whitespace", true, "   - [ ] foo bar"},
		{"done item", false, "   - [x] foo bar"},
		{"done item, capital X", false, "   - [X] foo bar"},
		{"not an item", false, " [] foo bar"},
		{"not an item", false, " - ] foo bar"},
		{"not an item", false, " - [ foo bar"},
		{"not an item", false, " - [] "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, todoItemUndoneRx.MatchString(tt.given))
		})
	}
}

func TestTODORegexDone(t *testing.T) {
	var tests = []struct {
		name     string
		expected bool
		given    string
	}{
		{"done", true, "   - [x] foo"},
		{"done, capital X", true, "   - [X] foo"},
		{"undone item", false, " [] foo bar"},
		{"not an item", false, " - x] foo bar"},
		{"not an item", false, " - [x foo bar"},
		{"not an item", false, " - [x] "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, todoItemDoneRx.MatchString(tt.given))
		})
	}
}

func TestShouldSort(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"signle list",
			`
			- [] foo
			- [] bazz
			- [x] bar` + "\n",
			`
			- [] foo
			- [x] bar
			- [] bazz`},
		{"signle list with surrounding content",
			`some random text
			- [] foo
			- [] bazz
			- [x] bar
			more random text` + "\n",
			`some random text
			- [] foo
			- [x] bar
			- [] bazz
			more random text`},
		{"single list with single empty lines",
			`
			- [] bar

			- [] foobar

			- [] foobazz
			- [x] foo
			- [x] buzz` + "\n",
			`
			- [x] foo
			- [] bar

			- [x] buzz
			- [] foobar

			- [] foobazz`},
		{"two lists separated by random text",
			`
			- [] bar
			- [x] foo
			random text
			- [] buzz
			- [] foobazz
			- [x] foobar` + "\n",
			`
			- [x] foo
			- [] bar
			random text
			- [] buzz
			- [x] foobar
			- [] foobazz`},
		{"two lists separated by more than one new line",
			`
			- [] bar
			- [x] foo


			- [] buzz
			- [] foobazz
			- [x] foobar` + "\n",
			`
			- [x] foo
			- [] bar


			- [] buzz
			- [x] foobar
			- [] foobazz`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var actualBuf bytes.Buffer
			assert.NoError(t, Sort(bytes.NewBufferString(tt.given), &actualBuf))
			assert.Equal(t, tt.expected, actualBuf.String())
		})
	}
}
