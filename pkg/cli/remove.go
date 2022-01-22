package cli

import (
	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "removes a given note",

	Args: cobra.MinimumNArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
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
