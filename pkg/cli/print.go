package cli

import (
	"os"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

var prnNoteCmd = &coral.Command{
	Use:     "print",
	Short:   "prints a given note",
	Aliases: []string{"cat"},

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		return printNotes(args)
	},
}

func init() {
	rootCmd.AddCommand(prnNoteCmd)
}

func printNote(notes note.List, name string) error {
	nt, err := notes.Get(name)
	if err != nil {
		return err
	}

	return nt.Dump(os.Stdout)
}

func printNotes(args []string) error {
	notes := note.NewList()

	if len(args) == 0 {
		// handle scratchpad
		noteName := noteNameFromArgs(args)
		return printNote(notes, noteName)
	}

	for _, arg := range args {
		noteName, err := parseNoteName(arg)
		if err != nil {
			return err
		}
		if err := printNote(notes, noteName); err != nil {
			return err
		}
	}
	return nil
}
