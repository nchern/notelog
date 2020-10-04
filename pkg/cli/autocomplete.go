package cli

import (
	"fmt"
	"os"
)

func autoCompleteScript() string {
	name := os.Args[0]
	return fmt.Sprintf("# Bash autocompletion for %s. Completes notes\ncomplete -W \"`%s -c=list`\" %s",
		name, name, name)
}
