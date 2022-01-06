package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldFormat(t *testing.T) {

	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"empty",
			text(""),
			text(""),
		},
		{"lines with spaces",
			text(""),
			text("   \t"),
		},
		{"header Wiki style",
			text(
				"### Header",
				"    content",
			),
			text(
				"### Header",
				" content",
			),
		},
		{"subtext in numbered list",
			text(
				"abc",
				"1. foo bar",
				"   fuzz buzz",
				"   subline",
				"22. hello world",
				"    - item 1",
				"    - item 2",
				"    subline2",
			),
			text(
				"abc",
				"1. foo bar",
				"fuzz buzz",
				"subline",
				"22. hello world",
				"- item 1",
				" - item 2",
				"  subline2",
			),
		},
		{"subtext with tabs in numbered list",
			text(
				"abc",
				"1. foo bar",
				"   fuzz buzz",
				"   subline",
				"22. hello world",
				"    - item 1",
				"    - item 2",
				"    subline2",
			),
			text(
				"abc",
				"1. foo bar",
				"fuzz buzz",
				"\t\t\tsubline",
				"22. hello world",
				"- item 1",
				"\t- item 2",
				"\t\t  subline2",
			),
		},
		{"subtext under org mode headers",
			text(
				"abc",
				"* Header 1",
				"  lineA",
				"  lineB",
				"** Header 2",
				"   - item 1",
				"   - item 2",
				"*** Header3",
				"    lineA",
				"    - item2",
			),
			text(
				"abc",
				"* Header 1",
				" lineA",
				"lineB",
				"** Header 2",
				"  - item 1",
				" - item 2",
				"*** Header3",
				" \t lineA",
				"\t- item2",
			),
		},
		// {"subtext under bullet list",
		// 	text(
		// 		"- foobar",
		// 		"  aaabbbcc",
		// 		"  aabb",
		// 		"- foo barr",
		// 		"  fuzzbuzz",
		// 		"  uzzbuzz",
		// 	),
		// 	text(
		// 		"- foobar",
		// 		"aaabbbcc",
		// 		"aabb",
		// 		"- foo barr",
		// 		"fuzzbuzz",
		// 		"    uzzbuzz",
		// 	),
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := &bytes.Buffer{}

			err := Format(bytes.NewBufferString(tt.given), actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual.String())
		})
	}
}

func text(lines ...string) string {
	return strings.Join(lines, "\n")
}
