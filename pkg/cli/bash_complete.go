package cli

import (
	"fmt"
	"os"

	"github.com/muesli/coral"
)

// TODO: This command is OBSOLETE and not used. TO BE REMOVED.

var bashCompleteCmd = &coral.Command{
	Use:   "bash-complete",
	Short: "returns autocompete initialization script for bashrc",

	Args: coral.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		_, err := fmt.Println(autoCompleteScript())
		return err
	},
}

func init() {
	rootCmd.AddCommand(bashCompleteCmd)
}

func autoCompleteScript() string {
	name := os.Args[0]
	return fmt.Sprintf("# Bash autocompletion for %s. Completes notes\ncomplete -C \"`%s do autocomplete`\" %s",
		name, name, name)
}
