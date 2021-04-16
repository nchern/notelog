package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var bashCompleteCmd = &cobra.Command{
	Use:   "bash-complete",
	Short: "returns autocompete initialization script for bashrc",

	Args: cobra.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Println(autoCompleteScript())
		return err
	},
}

func init() {
	doCmd.AddCommand(bashCompleteCmd)
}

func autoCompleteScript() string {
	name := os.Args[0]
	return fmt.Sprintf("# Bash autocompletion for %s. Completes notes\ncomplete -C \"`%s do autocomplete`\" %s",
		name, name, name)
}
