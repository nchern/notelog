package cli

import (
	"flag"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
)

func edit() error {
	notes := note.NewList()

	noteName, instantRecord, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}
	if instantRecord != "" {
		return editor.WriteInstantRecord(notes.Note(noteName), instantRecord)
	}
	return editor.Edit(notes.Note(noteName))
}
