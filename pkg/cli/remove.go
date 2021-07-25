package cli

import (
	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "removes a given note",

	Args: cobra.ExactArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var name string

		if filename {
			name = note.NameFromFilename(args[0])
			err = validateNoteName(name)
		} else {
			name, err = parseNoteName(args[0])
		}
		if err != nil {
			return err
		}
		return notes.Remove(name)
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
