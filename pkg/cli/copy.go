package cli

import (
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "cp",
	Short: "copies a given note",

	Args: cobra.ExactArgs(2),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := parseNoteName(args[0])
		if err != nil {
			return err
		}
		newName, err := parseNoteName(args[1])
		if err != nil {
			return err
		}

		return notes.Copy(name, newName)
	},
}

func init() {
	doCmd.AddCommand(copyCmd)
}
