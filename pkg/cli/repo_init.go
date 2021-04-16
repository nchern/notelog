package cli

import (
	"os"

	"github.com/nchern/notelog/pkg/repo"
	"github.com/spf13/cobra"
)

var repoInitCmd = &cobra.Command{
	Use:   "init-repo",
	Short: "Initialiazed a new repo for this note collection",

	Args: cobra.MinimumNArgs(0),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return repo.Init(notes, os.Stderr)
	},
}

func init() {
	doCmd.AddCommand(repoInitCmd)
}
