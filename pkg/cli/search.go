package cli

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/searcher"
)

var (
	interactive = flag.Bool(
		"interactive",
		false,
		"(search only) if set, search results are saved to a file under NOTELOG_HOME dir. Search results in output get numbered.")

	titlesOnly = flag.Bool(
		"titles-only",
		false,
		"(search only) if set, outputs note titles of search results only",
	)
)

func search() error {
	if len(flag.Args()) < 1 {
		return errors.New("Not enough args. Specify a search term")
	}

	notes := note.NewList()

	var out io.Writer = os.Stdout
	if *interactive {
		out = &nlWriter{inner: out}
	}
	s := searcher.NewSearcher(notes, out)

	s.OnlyNames = *titlesOnly
	s.SaveResults = *interactive

	err := s.Search(flag.Args()...)
	if err == searcher.ErrNoResults {
		os.Exit(1)
	}

	return err
}
