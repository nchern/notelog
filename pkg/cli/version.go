package cli

import (
	"fmt"

	"github.com/muesli/coral"
)

var (
	version string

	versionCmd = &coral.Command{
		Use:   "version",
		Short: "prints current version",
		Args:  coral.NoArgs,
		RunE: func(cmd *coral.Command, args []string) error {
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
