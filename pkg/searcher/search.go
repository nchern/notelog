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

var grepCmd = env.Get("NOTELOG_GREP", defaultGrep)

// Notes abstracts note collection to search in
type Notes interface {
	HomeDir() string
}

// Search runs the search over all notes in notes home and prints results to stdout
func Search(notes Notes, terms string) error {
	cmd := exec.Command(grepCmd, defaultGrepArgs, terms, notes.HomeDir())
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
