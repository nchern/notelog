package cli

import (
	"flag"
	"fmt"

	"github.com/nchern/notelog/pkg/note"
)

func printFullPath() error {
	notes := note.NewList()

	noteName, _, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}
	n := notes.Note(noteName)

	if ok, err := n.Exists(); !ok {
		if err != nil {
			return err
		}

		return fmt.Errorf("'%s' does not exist", noteName)
	}

	_, err = fmt.Print(n.FullPath())
	return err
}
