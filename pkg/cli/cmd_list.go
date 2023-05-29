package cli

import (
	"fmt"
	"io"
	"strings"
)

func cutString(s string, sep string) (a string, b string) {
	toks := strings.SplitN(s, " ", 2)
	if len(toks) > 0 {
		a = toks[0]
	}
	if len(toks) > 1 {
		b = toks[1]
	}
	return
}

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
