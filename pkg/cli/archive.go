package cli

import (
	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "moves a note to archive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return archive(args)
	},
}

func init() {
	doCmd.AddCommand(archiveCmd)
}

func archive(args []string) error {
	notes := note.NewList()

	noteName, _, err := parseArgs(args)
	if err != nil {
		return err
	}

	return notes.Archive(noteName)
}
