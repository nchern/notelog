package cli

import (
	"github.com/muesli/coral"
)

var renameCmd = &coral.Command{
	Use:     "rename",
	Aliases: []string{"mv"},
	Short:   "renames a given note",

	Args: coral.ExactArgs(2),

	SilenceErrors: true,
	SilenceUsage:  true,

	ValidArgsFunction: completeNoteNames,

	RunE: func(cmd *coral.Command, args []string) error {
		name, err := parseNoteName(args[0])
		if err != nil {
			return err
		}
		newName, err := parseNoteName(args[1])
		if err != nil {
			return err
		}
		return notes.Rename(name, newName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
