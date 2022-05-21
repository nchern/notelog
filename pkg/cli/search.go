package cli

import (
	"errors"
	"os"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/search"
)

var (
	interactive bool

	titlesOnly bool

	caseSensitive bool

	color bool
)

var searchCmd = &coral.Command{
	Use:   "search",
	Short: "runs search over notes collection",

	Args: coral.MinimumNArgs(1),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
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

func bindSearchFlags(cmd *coral.Command) {
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
	cmd.Flags().BoolVarP(&color,
		"color",
		"l",
		false,
		"if set, enables colored output")
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

	simpleRenderer := &search.StreamRenderer{
		W:         os.Stdout,
		OnlyNames: titlesOnly,
		Colorize:  color,
	}
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
