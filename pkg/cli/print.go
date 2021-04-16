package cli

import (
	"os"

	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var prnNoteCmd = &cobra.Command{
	Use:   "print",
	Short: "prints a given note",
	Args:  cobra.ExactArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return printNote(args)
	},
}

func init() {
	doCmd.AddCommand(prnNoteCmd)
}

func printNote(args []string) error {
	notes := note.NewList()

	noteName, _, err := parseArgs(args)
	if err != nil {
		return err
	}

	nt, err := notes.Get(noteName)
	if err != nil {
		return err
	}

	return nt.Dump(os.Stdout)
}
