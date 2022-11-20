package cli

import (
	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/search"
)

var grepCmd = &coral.Command{
	Use:   "grep",
	Short: "runs grep over notes collection",

	Args: coral.MinimumNArgs(0),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
		notes := note.NewList()

		s := search.NewGrepEngine(notes)
		s.OnlyNames = titlesOnly
		s.CaseSensitive = caseSensitive

		return doSearch(args, s)
	},
}

func init() {
	bindSearchFlags(grepCmd)

	rootCmd.AddCommand(grepCmd)
}
