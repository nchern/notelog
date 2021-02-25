package cli

import (
	"flag"

	"github.com/nchern/notelog/pkg/note"
)

func archive() error {
	notes := note.NewList()

	noteName, _, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}

	return notes.Archive(noteName)
}
