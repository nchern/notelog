package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/stretchr/testify/assert"
)

const (
	homeDir = "/tmp/test_notes"
)

func TestAutoComplete(t *testing.T) {
	names := []string{
		"bar",
		"buzz",
		"foo"}

	withDirs(names, func() {
		var tests = []struct {
			name     string
			given    string
			expected string
		}{
			{"should complete names",
				"notelog ",
				text(names...)},
			{"should complete subcommands",
				"notelog -c ",
				text(commands...)},
			{"should complete flag",
				"notelog -",
				text("-c")},
			{"should complete flag 2",
				"notelog -c",
				text("-c")},
			{"should complete subcommands with common prefix only",
				"notelog -c li",
				text("list", "list-cmds")},
			{"should complete subcommands with common prefix only-2",
				"notelog -c p",
				text("path", "print", "print-home")},
			{"should complete names with common prefix only",
				"notelog b",
				text("bar", "buzz")},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				w := &bytes.Buffer{}
				pos := len(tt.given) - 1

				assert.NoError(t,
					autoComplete(note.List(homeDir), tt.given, pos, w))
				assert.Equal(t, tt.expected+"\n", w.String())
			})
		}
	})
}

func text(lines ...string) string { return strings.Join(lines, "\n") }

func withDirs(dirs []string, fn func()) {
	must(os.MkdirAll(homeDir, 0755))
	defer os.RemoveAll(homeDir)

	for _, name := range dirs {
		must(os.MkdirAll(filepath.Join(homeDir, name), 0755))
	}

	fn()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}