package cli

import (
	"fmt"
	"io"
)

func printCommands(w io.Writer, filter func(string) bool) error {
	for _, c := range rootCmd.Commands() {
		// handle the case when Use has a format of "help [command]"
		cmd, _ := cutString(c.Use, " ")
		if !filter(cmd) {
			continue
		}
		if _, err := fmt.Fprintln(w, cmd); err != nil {
			return err
		}
	}
	return nil
}

func listCommands(w io.Writer) error {
	return printCommands(w, func(s string) bool { return true })
}
