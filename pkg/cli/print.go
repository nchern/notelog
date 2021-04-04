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

	nt, err := notes.Get(noteName)
	if err != nil {
		return err
	}

	return nt.Dump(os.Stdout)
}
