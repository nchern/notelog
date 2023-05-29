package cli

import (
	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

var (
	archiveCmd = &coral.Command{
		Use:   "archive",
		Short: "moves a note to archive",
		Args:  coral.ExactArgs(1),

		SilenceErrors: true,
		SilenceUsage:  true,

		ValidArgsFunction: completeNoteNames,

		RunE: func(cmd *coral.Command, args []string) error {
			return archive(args)
		},
	}

	filename bool
)

func init() {
	archiveCmd.Flags().BoolVarP(&filename,
		"filename",
		"f",
		false,
		"addresses a note by filename")

	rootCmd.AddCommand(archiveCmd)
}

func archive(args []string) error {
	notes := note.NewList()

	var err error
	var noteName string

	if filename {
		noteName = note.NameFromFilename(args[0])
		err = validateNoteName(noteName)
	} else {
		noteName, err = parseNoteName(args[0])
	}
	if err != nil {
		return err
	}

	return notes.Archive(noteName)
}
