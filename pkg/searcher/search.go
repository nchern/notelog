package searcher

import (
	"os"
	"os/exec"

	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultGrep     = "grep"
	defaultGrepArgs = "-rni"
)

var grepCmd = env.Get("NOTELOG_GREP", defaultGrepArgs)

// Search runs the search over all notes in notes home and prints results to stdout
func Search(terms string) error {

	cmd := exec.Command(grepCmd, defaultGrepArgs, terms, env.NotesRootPath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	switch e := err.(type) {
	case *exec.ExitError:
		if e.ExitCode() == 1 {
			os.Exit(1)
		}
	}
	return err
}
