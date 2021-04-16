package cli

import (
	"os"

	"github.com/nchern/notelog/pkg/checklist"
	"github.com/spf13/cobra"
)

var sortChecklistCmd = &cobra.Command{
	Use:   "sort-checklist",
	Short: "Sorts checklist given on stdin",

	Args: cobra.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return checklist.Sort(os.Stdin, os.Stdout)
	},
}

func init() {
	doCmd.AddCommand(sortChecklistCmd)
}
