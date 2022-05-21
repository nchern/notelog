package cli

import (
	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

var touchCmd = &coral.Command{
	Use:   "touch",
	Short: "runs a given command to manipulate notes",

	Args: coral.ExactArgs(1),

	SilenceUsage:  true,
	SilenceErrors: true,

	RunE: func(cmd *coral.Command, args []string) error {
		return touch(notes, args)
	},
}

func init() {
	addFormatFlag(touchCmd)
	doCmd.AddCommand(touchCmd)
}

func touch(notes note.List, args []string) error {
	t, err := note.ParseFormat(conf.NoteFormat)
	if err != nil {
		return err
	}

	name, err := parseNoteName(args[0])
	if err != nil {
		return err
	}
	_, err = notes.GetOrCreate(name, t)
	return err
}
