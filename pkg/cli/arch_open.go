package cli

import (
	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
)

var (
	archOpenCmd = &coral.Command{
		Use:   "arch-open",
		Short: "opens a note from archive",
		Args:  coral.ExactArgs(1),

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *coral.Command, args []string) error {
			return archOpen(args)
		},
	}
)

func init() {
	doCmd.AddCommand(archOpenCmd)
}

func archOpen(args []string) error {
	var err error
	var noteName string

	notes := note.NewList().GetArchive()

	noteName, err = parseNoteName(args[0])
	if err != nil {
		return err
	}

	nt, err := notes.Get(noteName)
	if err != nil {
		return err
	}

	var lnum editor.LineNumber
	return editor.Edit(nt, readOnly, lnum)
}
