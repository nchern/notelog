package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var printHomeCmd = &cobra.Command{
	Use:   "print-home",
	Short: "prints current NOTELOG_HOME value",

	Args: cobra.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Print(notes.HomeDir())
		return err
	},
}

func init() {
	doCmd.AddCommand(printHomeCmd)
}
