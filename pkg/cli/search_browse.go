package cli

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/searcher"
)

func browseSearch() error {
	notes := note.NewList()
	n, err := parseNumber(flag.Args())
	if err != nil {
		return err
	}
	r, err := searcher.GetLastNthResult(notes, int(n))
	if err != nil {
		return err
	}
	if r == "" {
		return nil
	}
	toks := strings.Split(r, ":")

	return editor.Shellout(toks[0], "+"+toks[1]).Run()
}

func parseNumber(args []string) (int64, error) {
	if len(args) < 1 {
		return -1, errors.New("Not enough args. Specify a search result ordinal number")
	}
	return strconv.ParseInt(args[0], 10, 64)
}
