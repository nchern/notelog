package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSort(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"trivial case - no changes",
			text(
				"foo bar",
				"buzz fuzz",
				"* TODO a1",
				"content a1",
				"* TODO a2",
				"* non todo header",
				"** subheadr",
			),
			text(
				"foo bar",
				"buzz fuzz",
				"* TODO a1",
				"content a1",
				"* TODO a2",
				"* non todo header",
				"** subheadr",
			)},
		{"mixed todos and dones",
			text(
				"* TODO a2",
				"  content a2",
				"* TODO a4",
				"  content a4",
				"* DONE a1",
				"  content a1",
				"* DONE a3",
				"  content a3",
			),
			text(
				"* DONE a1",
				"  content a1",
				"* TODO a2",
				"  content a2",
				"* DONE a3",
				"  content a3",
				"* TODO a4",
				"  content a4",
			)},
		{"mixed todos, dones and non-todo headings",
			text(
				"* TODO a2",
				"  content a2.1",
				"  content a2.2",
				"* TODO a5",
				"  content a5",
				"* DONE a1",
				"  content a1.1",
				"  content a1.2",
				"  content a1.3",
				"* DONE a3",
				"* Non todo header",
				"** subheadr",
			),
			text(
				"* DONE a1",
				"  content a1.1",
				"  content a1.2",
				"  content a1.3",
				"* TODO a2",
				"  content a2.1",
				"  content a2.2",
				"* DONE a3",
				"* TODO a5",
				"  content a5",
				"* Non todo header",
				"** subheadr",
			)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := &bytes.Buffer{}
			in := bytes.NewBufferString(tt.given)

			assert.NoError(t, Sort(in, actual))
			assert.Equal(t, tt.expected+"\n", actual.String())
		})
	}
}

func text(lines ...string) string {
	return strings.Join(lines, "\n")
}
