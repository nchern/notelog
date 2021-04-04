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
	nt, err := notes.Get(noteName)
	if err != nil {
		return err
	}

	_, err = fmt.Print(nt.FullPath())
	return err
}
