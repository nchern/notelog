package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nchern/notelog/pkg/note"
)

func autoCompleteScript() string {
	name := os.Args[0]
	return fmt.Sprintf("# Bash autocompletion for %s. Completes notes\ncomplete -C \"`%s -c=autocomplete`\" %s",
		name, name, name)
}

func autoComplete(list note.List, line string, i int, w io.Writer) error {
	const cmdFlag = "-" + subCommand

	beforeCursor := line[0 : i+1]
	curTok := getCurrentCompletingToken(beforeCursor)
	prevToks := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(beforeCursor), curTok))

	if strings.HasPrefix(curTok, "-") {
		_, err := fmt.Fprintln(w, cmdFlag)
		return err
	}

	if strings.HasSuffix(prevToks, cmdFlag) {
		return printCommandsWithPrefix(curTok, w)
	}

	notes, err := list.All()
	if err != nil {
		return err
	}
	return printNotesWithPrefix(notes, curTok, w)
}

func printCommandsWithPrefix(prefix string, w io.Writer) error {
	for _, c := range commands {
		if !strings.HasPrefix(c, prefix) {
			continue
		}
		if _, err := fmt.Fprintln(w, c); err != nil {
			return err
		}
	}
	return nil
}

func printNotesWithPrefix(notes []*note.Note, prefix string, w io.Writer) error {
	for _, note := range notes {
		if !strings.HasPrefix(note.Name(), prefix) {
			continue
		}
		if _, err := fmt.Fprintln(w, note.Name()); err != nil {
			return err
		}
	}
	return nil
}

func getCurrentCompletingToken(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ' ' {
			return s[i+1:]
		}
	}
	return s
}
