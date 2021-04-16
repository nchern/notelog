package cli

import (
	"errors"
	"strconv"
	"strings"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/searcher"
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
