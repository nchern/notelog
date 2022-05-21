package cli

import (
	"fmt"

	"github.com/muesli/coral"
)

var printHomeCmd = &coral.Command{
	Use:   "print-home",
	Short: "prints current NOTELOG_HOME value",

	Args: coral.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		_, err := fmt.Print(notes.HomeDir())
		return err
	},
}

func init() {
	doCmd.AddCommand(printHomeCmd)
}
