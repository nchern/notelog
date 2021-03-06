package cli

import (
	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var (
	archiveCmd = &cobra.Command{
		Use:   "archive",
		Short: "moves a note to archive",
		Args:  cobra.ExactArgs(1),

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
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

	doCmd.AddCommand(archiveCmd)
}

func archive(args []string) error {
	notes := note.NewList()

	var err error
	var noteName string

	if filename {
		noteName = note.NameFromFilename(args[0])
		err = validateNoteName(noteName)
	} else {
		noteName, _, err = parseArgs(args)
	}
	if err != nil {
		return err
	}

	return notes.Archive(noteName)
}
