package cli

import (
	"flag"
	"os"

	"github.com/nchern/notelog/pkg/note"
)

func printNote() error {
	notes := note.NewList()

	noteName, _, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}

	n := notes.Note(noteName)

	return n.Dump(os.Stdout)
}
