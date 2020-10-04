package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nchern/notelog/pkg/note"
)

func autoCompleteScript() string {
	name := os.Args[0]
	return fmt.Sprintf("# Bash autocompletion for %s. Completes notes\ncomplete -W \"`%s -c=list`\" %s",
		name, name, name)
}

func autoComplete(list note.List, line string, i int, w io.Writer) error {
	line = strings.TrimSpace(line)
	if strings.HasSuffix(line, "-"+subCommand) {
		return listCommands(w)
	}

	return listNotes(list, w)
}
