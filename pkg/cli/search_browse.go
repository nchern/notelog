package cli

import (
	"errors"
	"strconv"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/search"
	"github.com/spf13/cobra"
)

var browseSearchCmd = &cobra.Command{
	Use:   "search-browse",
	Short: "runs search over notes collection",

	Args: cobra.ExactArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return browseSearch(args)
	},
}

func init() {
	doCmd.AddCommand(browseSearchCmd)
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

	noteName, lnum := parseNoteNameAndLineNumber(r)
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
