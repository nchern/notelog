package cli

import (
	"fmt"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/env"
)

var envCmd = &coral.Command{
	Use:   "env",
	Short: "prints env vars",

	Args: coral.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		_, err := fmt.Println(env.Vars())
		return err
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
