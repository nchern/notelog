package cli

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "removes a given note",

	Args: cobra.ExactArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		name, _, err := parseArgs(args)
		if err != nil {
			return err
		}
		return notes.Remove(name)
	},
}

func init() {
	doCmd.AddCommand(removeCmd)
}
