package cli

import (
	"os"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/repo"
)

var syncCmd = &coral.Command{
	Use:   "sync",
	Short: "syncs repo with the remote",

	Args: coral.MinimumNArgs(0),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		comment := ""
		if len(args) > 0 {
			comment = args[0]
		}
		return repo.Sync(notes, comment, os.Stderr)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
