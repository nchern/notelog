package cli

import (
	"errors"
	"flag"
	"os"
	"os/exec"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/searcher"
)

func search() error {
	if len(flag.Args()) < 1 {
		return errors.New("Not enough args. Specify a search term")
	}

	notes := note.NewList()
	err := searcher.Search(notes, os.Stdout, flag.Args()...)
	switch e := err.(type) {
	case *exec.ExitError:
		if e.ExitCode() == 1 {
			os.Exit(1)
		}
	}
	return err
}
