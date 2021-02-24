package cli

import (
	"flag"

	"github.com/nchern/notelog/pkg/note"
)

func touch(notes note.List) error {
	name, _, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}
	return notes.Note(name).Touch()
}
