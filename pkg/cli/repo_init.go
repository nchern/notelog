package cli

import (
	"os"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/repo"
)

var repoInitCmd = &coral.Command{
	Use:   "init-repo",
	Short: "Initialiazed a new repo for this note collection",

	Args: coral.MinimumNArgs(0),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		return repo.Init(notes, os.Stderr)
	},
}

func init() {
	doCmd.AddCommand(repoInitCmd)
}
