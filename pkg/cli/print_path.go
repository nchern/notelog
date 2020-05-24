package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/nchern/notelog/pkg/note"
)

func printFullPath() error {
	notes := note.NewList()

	noteName, _, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}
	n := notes.Note(noteName)
	path := n.FullPath()
	if _, err := os.Stat(path); err != nil {
		return err
	}
	_, err = fmt.Print(path)
	return err
}
