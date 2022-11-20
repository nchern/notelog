package cli

import (
	"errors"
	"strconv"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/search"
)

var browseSearchCmd = &coral.Command{
	Use:   "search-browse",
	Short: "runs search over notes collection",

	Args: coral.ExactArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		return browseSearch(args)
	},
}

func init() {
	rootCmd.AddCommand(browseSearchCmd)
}

func browseSearch(args []string) error {
	notes := note.NewList()
	n, err := parseNumber(args)
	if err != nil {
		return err
	}
	r, err := search.GetLastNthResult(notes, int(n))
	if err != nil {
		return err
	}
	if r == "" {
		return nil
	}

	noteName, lnum, isArchive := parseNoteNameAndLineNumber(r)
	if isArchive {
		notes = notes.GetArchive()
	}
	nt, err := notes.Get(noteName)
	if err != nil {
		return err
	}

	return editor.Edit(nt, false, lnum)
}

func parseNumber(args []string) (int64, error) {
	if len(args) < 1 {
		return -1, errors.New("Not enough args. Specify a search result ordinal number")
	}
	return strconv.ParseInt(args[0], 10, 64)
}
