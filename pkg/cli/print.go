package cli

import (
	"os"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

var prnNoteCmd = &coral.Command{
	Use:     "print",
	Short:   "prints a given note",
	Args:    coral.MinimumNArgs(1),
	Aliases: []string{"cat"},

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		return printNote(args)
	},
}

func init() {
	doCmd.AddCommand(prnNoteCmd)
}

func printNote(args []string) error {
	notes := note.NewList()

	for _, arg := range args {
		noteName, err := parseNoteName(arg)
		if err != nil {
			return err
		}

		nt, err := notes.Get(noteName)
		if err != nil {
			return err
		}

		if err = nt.Dump(os.Stdout); err != nil {
			return err
		}
	}
	return nil
}
