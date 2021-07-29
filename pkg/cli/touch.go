package cli

import (
	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var touchCmd = &cobra.Command{
	Use:   "touch",
	Short: "runs a given command to manipulate notes",

	Args: cobra.ExactArgs(1),

	SilenceUsage:  true,
	SilenceErrors: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return touch(notes, args)
	},
}

func init() {
	addFormatFlag(touchCmd)
	doCmd.AddCommand(touchCmd)
}

func touch(notes note.List, args []string) error {
	t, err := note.ParseFormat(noteFormat)
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
