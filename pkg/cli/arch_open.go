package cli

import (
	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
)

const archOpenCmdName = "arch-open"

var (
	archOpenCmd = &coral.Command{
		Use:   archOpenCmdName,
		Short: "opens a note from archive",
		Args:  coral.ExactArgs(1),

		SilenceErrors: true,
		SilenceUsage:  true,

		ValidArgsFunction: completeNoteNames,

		RunE: func(cmd *coral.Command, args []string) error {
			return archOpen(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(archOpenCmd)
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
