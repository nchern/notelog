package cli

import (
	"flag"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
)

const (
	// should be configurable
	skipLines uint = 2
)

func edit(readOnly bool) error {
	notes := note.NewList()

	noteName, instantRecord, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}
	nt := notes.Note(noteName)
	if instantRecord != "" {
		return nt.WriteInstantRecord(instantRecord, skipLines)
	}
	return editor.Edit(nt, readOnly)
}
