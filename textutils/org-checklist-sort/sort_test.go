package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexNotDone(t *testing.T) {
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
			assert.Equal(t, tt.expected, todoItemNotDoneRx.MatchString(tt.given))
		})
	}
}

func TestRegexDone(t *testing.T) {
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
			- [x] bar`,
			`
			- [] foo
			- [x] bar
			- [] bazz`},
		{"signle list with surrounding content",
			`some random text
			- [] foo
			- [] bazz
			- [x] bar
			more random text`,
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
			- [x] buzz`,
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
			- [x] foobar`,
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
			- [x] foobar`,
			`
			- [x] foo
			- [] bar


			- [] buzz
			- [x] foobar
			- [] foobazz`},

		{"hierarchical todos",
			`foo bar
				- [] bar
				- [] buzz
				- [] fuzzbuzz
					- [] fuzzbuzz - 1
					- [] fuzzbuzz - 2
				- [] foobar
				- [x] foo
					- [x] sub foo
					- [ ] sub bar
				- [x] foobazz`,
			`foo bar
				- [] bar
				- [x] foo
					- [x] sub foo
					- [ ] sub bar
				- [] buzz
				- [x] foobazz
				- [] fuzzbuzz
					- [] fuzzbuzz - 1
					- [] fuzzbuzz - 2
				- [] foobar`},
		{"hierarchical todos with non-todo sub items",
			`foo bar
				- [] bar
				- [] buzz
				- [] foobar
				- [x] foo
                  - foo context
					more context
				- [x] foobazz`,
			`foo bar
				- [] bar
				- [x] foo
                  - foo context
					more context
				- [] buzz
				- [x] foobazz
				- [] foobar`},
		{"empty line after list",
			`
				- [ ] foo
				- [x] bar
				- [x] buzzbar

				text`,
			`
				- [x] bar
				- [x] buzzbar
				- [ ] foo

				text`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var actualBuf bytes.Buffer
			assert.NoError(t, Sort(bytes.NewBufferString(tt.given), &actualBuf))
			assert.Equal(t, tt.expected+"\n", actualBuf.String())
		})
	}
}
