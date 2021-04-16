package cli

import (
	"fmt"

	"github.com/nchern/notelog/pkg/env"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "prints env vars",

	Args: cobra.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Println(env.Vars())
		return err
	},
}

func init() {
	doCmd.AddCommand(envCmd)
}
