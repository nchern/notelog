package cli

import (
	"errors"
	"flag"
	"os"
	"os/exec"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/searcher"
)

var saveResults = flag.Bool("save-results", false, "(search only) if set, search results are saved to a file under NOTELOG_HOME dir")

func search() error {
	if len(flag.Args()) < 1 {
		return errors.New("Not enough args. Specify a search term")
	}

	notes := note.NewList()

	s := searcher.NewSearcher(notes, os.Stdout)
	s.SaveResults = *saveResults

	err := s.Search(flag.Args()...)

	switch e := err.(type) {
	case *exec.ExitError:
		if e.ExitCode() == 1 {
			os.Exit(1)
		}
	}
	return err
}
