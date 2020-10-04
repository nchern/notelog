package cli

import (
	"fmt"
	"io"
)

func listCommands(w io.Writer) error {
	for _, c := range commands {
		fmt.Fprintln(w, c)
	}

	return nil
}
