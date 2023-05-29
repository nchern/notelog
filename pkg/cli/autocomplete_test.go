package cli

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mkFiles(names ...string) map[string]string {
	files := map[string]string{}
	for _, name := range names {
		files[name] = ""
	}
	return files
}

func TestAutoCompleteShould(t *testing.T) {
	names := []string{"bar", "buzz", "drum", "foo"}

	allCommands := bytes.Buffer{}
	require.NoError(t, runCommandWithEnv("/tmp", &allCommands, "list-cmds"))

	files := mkFiles(names...)
	testutil.WithNotes(files, func(notes note.List) {
		var tests = []struct {
			name     string
			given    string
			expected string
		}{
			{"complete subcommands",
				"notelog ",
				allCommands.String()},
			{"complete subcommands with common prefix only",
				"notelog li",
				text("list", "list-cmds")},
			{"complete subcommands with common prefix only-2",
				"notelog p",
				text("path", "print", "print-home")},
			{"complete names after edit command",
				"notelog edit b",
				text("bar", "buzz")},
			{"complete names after subcommands and already given names",
				"notelog edit foo b",
				text("bar", "buzz")},
			{"complete note names subcommand edit",
				"notelog edit ",
				text("bar", "buzz", "drum", "foo")},
			{"complete help command correctly",
				"notelog hel",
				text("help")},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				w := &bytes.Buffer{}

				// HACK: setting env variables mean global state modification. Bug-prone
				require.NoError(t, os.Setenv("COMP_LINE", tt.given))
				require.NoError(t, os.Setenv("COMP_POINT", strconv.Itoa(len(tt.given))))
				require.NoError(t, runCommandWithEnv(notes.HomeDir(), w, "autocomplete"))

				assert.Equal(t, tt.expected, w.String())
			})
		}
	})
}

func runCommandWithEnv(homeDir string, w io.Writer, args ...string) error {
	rootCmd.SetOut(w)
	rootCmd.SetErr(w)
	rootCmd.SetArgs(args)
	if err := os.Setenv("NOTELOG_HOME", homeDir); err != nil {
		return err
	}
	return rootCmd.Execute()
}

func TestAutoCompleteWithArchivedNotesShould(t *testing.T) {
	files := mkFiles("bar", "buzz", "drum", "foo", "foobar")
	testutil.WithNotes(files, func(notes note.List) {
		require.NoError(t, notes.Archive("bar"))
		require.NoError(t, notes.Archive("foobar"))

		var tests = []struct {
			name     string
			given    string
			expected string
		}{
			{"complete non-archived names",
				"notelog edit ",
				text("buzz", "drum", "foo")},
			{"complete archived names for arch-open",
				"notelog arch-open ",
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
