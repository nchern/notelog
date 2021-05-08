package cli

import (
	"errors"
	"io"
	"os"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/searcher"
	"github.com/spf13/cobra"
)

var (
	interactive bool

	titlesOnly bool

	caseSensitive bool
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "runs search over notes collection",

	Args: cobra.MinimumNArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return search(args)
	},
}

func init() {
	searchCmd.Flags().BoolVarP(&interactive,
		"interactive",
		"i",
		false,
		"if set, search results are saved to a file under NOTELOG_HOME dir. Search results in output get numbered.")

	searchCmd.Flags().BoolVarP(&titlesOnly,
		"titles-only",
		"t",
		false,
		"if set, outputs note titles of search results only")

	searchCmd.Flags().BoolVarP(&caseSensitive,
		"case-sensitive",
		"c",
		false,
		"if set, runs case sensitive search")

	doCmd.AddCommand(searchCmd)
}

func search(args []string) error {
	if len(args) < 1 {
		return errors.New("Not enough args. Specify a search term")
	}

	notes := note.NewList()

	var out io.Writer = os.Stdout
	if interactive {
		out = &nlWriter{inner: out}
	}
	s := searcher.NewSearcher(notes, out)

	s.OnlyNames = titlesOnly
	s.SaveResults = interactive
	s.CaseSensitive = caseSensitive
	err := s.Search(args...)
	if err == searcher.ErrNoResults {
		os.Exit(1)
	}

	return err
}
