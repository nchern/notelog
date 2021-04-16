package cli

import (
	"github.com/nchern/notelog/pkg/repo"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "syncs repo with the remote",

	Args: cobra.MinimumNArgs(0),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		comment := ""
		if len(args) > 0 {
			comment = args[0]
		}
		return repo.Sync(notes, comment)
	},
}

func init() {
	doCmd.AddCommand(syncCmd)
}
