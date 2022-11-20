package cli

import (
	"fmt"
	"io"
)

func listCommands(w io.Writer) error {
	for _, c := range rootCmd.Commands() {
		fmt.Fprintln(w, c.Use)
	}

	return nil
}
