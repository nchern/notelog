package cli

import (
	"errors"
	"os"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/search"
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
		if len(args) < 1 {
			return errors.New("Not enough args. Specify a search term")
		}

		notes := note.NewList()

		s := search.NewEngine(notes)

		s.OnlyNames = titlesOnly
		s.CaseSensitive = caseSensitive

		return doSearch(args, s)
	},
}

func bindSearchFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&interactive,
		"interactive",
		"i",
		false,
		"if set, search results are saved to a file under NOTELOG_HOME dir. Search results in output get numbered.")

	cmd.Flags().BoolVarP(&titlesOnly,
		"titles-only",
		"t",
		false,
		"if set, outputs note titles of search results only")

	cmd.Flags().BoolVarP(&caseSensitive,
		"case-sensitive",
		"c",
		false,
		"if set, runs case sensitive search")
}

func init() {
	bindSearchFlags(searchCmd)

	doCmd.AddCommand(searchCmd)
}

type searcher interface {
	Search(...string) ([]*search.Result, error)
}

func doSearch(args []string, s searcher) error {
	res, err := s.Search(args...)
	if err != nil {
		return err
	}
	if len(res) == 0 {
		os.Exit(1)
	}

	simpleRenderer := &search.StreamRenderer{W: os.Stdout, OnlyNames: titlesOnly}
	var renderer search.Renderer = simpleRenderer
	if interactive {
		simpleRenderer.W = &nlWriter{inner: os.Stdout}
		renderer, err = search.NewPersistentRenderer(notes, simpleRenderer)
		if err != nil {
			return err
		}
	}

	return search.Render(renderer, res, titlesOnly)
}
