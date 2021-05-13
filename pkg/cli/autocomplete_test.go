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
		"drum",
		"foo"}

	allCommands := bytes.Buffer{}
	listCommands(&allCommands)

	withDirs(names, func() {
		var tests = []struct {
			name     string
			given    string
			expected string
		}{
			{"should complete names",
				"notelog ",
				text(append([]string{"do"}, names...)...)},
			{"should complete subcommands",
				"notelog do ",
				strings.TrimSuffix(allCommands.String(), "\n")},
			{"should complete do command",
				"notelog d",
				text("do", "drum")},
			{"should complete exclude do comand",
				"notelog dr",
				text("drum")},
			{"should complete flag 2",
				"notelog do",
				text("do")},
			{"should complete subcommands with common prefix only",
				"notelog do li",
				text("list", "list-cmds")},
			{"should complete subcommands with common prefix only-2",
				"notelog do p",
				text("path", "print", "print-home")},
			{"should complete names with common prefix only",
				"notelog b",
				text("bar", "buzz")},
			{"should complete names after subcommands",
				"notelog do edit b",
				text("bar", "buzz")},
			{"should complete names after subcommands and already given names",
				"notelog do edit foo b",
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
