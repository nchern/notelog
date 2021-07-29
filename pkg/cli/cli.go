package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

const (
	scratchpadName = ".scratchpad"
)

var (
	notes = note.NewList()

	doCmd = &cobra.Command{
		Use:   "do",
		Short: "runs a given command to manipulate notes",
		Args:  cobra.ExactArgs(1),

		SilenceErrors: true,
		SilenceUsage:  false,

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd = &cobra.Command{
		Use:   "notelog",
		Short: "Efficient CLI personal note manager",
		Args:  cobra.MinimumNArgs(0),

		SilenceUsage:  true,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return edit(args, false)
		},
	}

	// HACK
	lsCmdsCmd = &cobra.Command{
		Use:   "list-cmds",
		Short: "lists all subcommands",
		Args:  cobra.NoArgs,

		SilenceErrors: true,
		SilenceUsage:  false,

		RunE: func(cmd *cobra.Command, args []string) error {
			return listCommands(os.Stdout)
		},
	}

	defaultHelp = rootCmd.HelpFunc()
)

func init() {
	doCmd.AddCommand(lsCmdsCmd)

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, s []string) {
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
