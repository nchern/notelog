package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	homeDir = "/tmp/test_notes"
)

func mkFiles(names ...string) map[string]string {
	files := map[string]string{}
	for _, name := range names {
		files[name] = ""
	}
	return files
}

func TestAutoComplete(t *testing.T) {
	names := []string{"bar", "buzz", "drum", "foo"}

	allCommands := bytes.Buffer{}
	listCommands(&allCommands)

	files := mkFiles(names...)
	testutil.WithNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			given    string
			expected string
		}{
			{"should complete names and do subcommand",
				"notelog ",
				text(append([]string{"do"}, names...)...)},
			{"should complete subcommands",
				"notelog do ",
				allCommands.String()},
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
			{"should complete note names subcommand edit",
				"notelog do edit ",
				text("bar", "buzz", "drum", "foo")},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				w := &bytes.Buffer{}
				pos := len(tt.given) - 1

				assert.NoError(t,
					autoComplete(notes, tt.given, pos, w))
				assert.Equal(t, tt.expected, w.String())
			})
		}
	})
}

func TestAutoCompleteWithArchivedNotes(t *testing.T) {
	files := mkFiles("bar", "buzz", "drum", "foo", "foobar")
	testutil.WithNotes(files, func(notes note.List) {
		require.NoError(t, notes.Archive("bar"))
		require.NoError(t, notes.Archive("foobar"))

		var tests = []struct {
			name     string
			given    string
			expected string
		}{
			{"should complete non-archived names",
				"notelog ",
				text("do", "buzz", "drum", "foo")},
			{"should complete archived names for arch-open",
				"notelog do arch-open ",
				text("bar", "foobar")},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				pos := len(tt.given) - 1

				w := &bytes.Buffer{}
				assert.NoError(t,
					autoComplete(notes, tt.given, pos, w))
				assert.Equal(t, tt.expected, w.String())
			})
		}
	})
}

func text(lines ...string) string { return strings.Join(lines, "\n") + "\n" }
