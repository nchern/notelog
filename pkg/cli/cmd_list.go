package cli

import (
	"fmt"
	"io"
	"sort"
)

func printCommands(w io.Writer, filter func(string) bool) error {
	cmds := make([]string, 0, 2*len(rootCmd.Commands()))
	for _, c := range rootCmd.Commands() {
		// handle the case when Use has a format of "help [command]"
		cmd, _ := cutString(c.Use, " ")
		cmds = append(cmds, cmd)
		cmds = append(cmds, c.Aliases...)
	}
	sort.Strings(cmds)
	for _, cmd := range cmds {
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
