package cli

import (
	"flag"
	"io"
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

	f, err := os.Open(n.FullPath())
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, f)
	return err
}
