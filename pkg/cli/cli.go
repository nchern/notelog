package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

const (
	scratchpadName = ".scratchpad"
)

var (
	notes = note.NewList()

	doCmd = &coral.Command{
		Use:   cmdDo,
		Short: "runs a given command to manipulate notes",
		Args:  coral.ExactArgs(1),

		SilenceErrors: true,
		SilenceUsage:  false,

		Run: func(cmd *coral.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd = &coral.Command{
		Use:   "notelog",
		Short: "Efficient CLI personal note manager",
		Args:  coral.MinimumNArgs(0),

		SilenceUsage:  true,
		SilenceErrors: true,

		RunE: func(cmd *coral.Command, args []string) error {
			return edit(args, false)
		},
	}

	// HACK
	lsCmdsCmd = &coral.Command{
		Use:     "list-cmds",
		Short:   "lists all subcommands",
		Aliases: []string{"ls-cmds"},
		Args:    coral.NoArgs,

		SilenceErrors: true,
		SilenceUsage:  false,

		RunE: func(cmd *coral.Command, args []string) error {
			return listCommands(os.Stdout)
		},
	}

	defaultHelp = rootCmd.HelpFunc()
)

func init() {
	doCmd.AddCommand(lsCmdsCmd)

	rootCmd.SetHelpFunc(func(cmd *coral.Command, s []string) {
		defaultHelp(cmd, s)

		fmt.Println()
		fmt.Println("Use \"notelog <note-name>\" as a shortcut of \"notelog do edit <note-name>\"")
	})

	rootCmd.AddCommand(doCmd)
}

// Execute is an entry point of CLI
func Execute() error {
	return rootCmd.Execute()
}

func parseNoteName(name string) (string, error) {
	// FIXME: this is a hack, need more elegant solution than double if
	if name == scratchpadName {
		return name, nil
	}
	name = strings.TrimSpace(name)
	if err := validateNoteName(name); err != nil {
		return "", err
	}
	return name, nil
}

func noteNameFromArgs(args []string) string {
	if len(args) < 1 {
		return scratchpadName
	}
	return args[0]
}
