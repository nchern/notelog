package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "prints current version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printVersion()
		},
	}
)

func init() {
	doCmd.AddCommand(versionCmd)
}

func printVersion() error {
	_, err := fmt.Printf("notelog version %s\n", version)
	return err
}
