package cli

import (
	"fmt"

	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var printPathCmd = &cobra.Command{
	Use:   "path",
	Short: "shows full path to a given note",

	Args: cobra.MaximumNArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return printFullPath(args)
	},
}

func init() {
	doCmd.AddCommand(printPathCmd)
}

func printFullPath(args []string) error {
	notes := note.NewList()

	noteName, err := parseNoteName(noteNameFromArgs(args))
	if err != nil {
		return err
	}
	nt, err := notes.Get(noteName)
	if err != nil {
		return err
	}

	_, err = fmt.Print(nt.FullPath())
	return err
}
