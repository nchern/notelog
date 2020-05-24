package cli

import (
	"flag"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/searcher"
)

func search() error {
	notes := note.NewList()

	terms, err := parseSearchArgs(flag.Args())
	if err != nil {
		return err
	}
	return searcher.Search(notes, terms)
}
