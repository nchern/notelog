package cli

import (
	"github.com/muesli/coral"
)

var (
	// HACK
	lsCmdsCmd = &coral.Command{
		Use:     "list-cmds",
		Short:   "lists all subcommands",
		Aliases: []string{"ls-cmds"},
		Args:    coral.NoArgs,

		SilenceErrors: true,
		SilenceUsage:  false,

		RunE: func(cmd *coral.Command, args []string) error {
			return listCommands(cmd.OutOrStdout())
		},
	}
)

func init() {
	rootCmd.AddCommand(lsCmdsCmd)
}
