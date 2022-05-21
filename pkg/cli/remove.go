package cli

import (
	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

var removeCmd = &coral.Command{
	Use:   "rm",
	Short: "removes a given note",

	Args: coral.MinimumNArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		var err error
		var name string

		for _, arg := range args {
			if filename {
				name = note.NameFromFilename(arg)
				err = validateNoteName(name)
			} else {
				name, err = parseNoteName(arg)
			}
			if err != nil {
				return err
			}
			err = notes.Remove(name)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&filename,
		"filename",
		"f",
		false,
		"addresses a note by filename")

	doCmd.AddCommand(removeCmd)
}
