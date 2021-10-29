package cli

import (
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/search"
	"github.com/spf13/cobra"
)

var grepCmd = &cobra.Command{
	Use:   "grep",
	Short: "runs grep over notes collection",

	Args: cobra.MinimumNArgs(0),

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		notes := note.NewList()

		s := search.NewGrepEngine(notes)
		s.OnlyNames = titlesOnly
		s.CaseSensitive = caseSensitive

		return doSearch(args, s)
	},
}

func init() {
	bindSearchFlags(grepCmd)

	doCmd.AddCommand(grepCmd)
}
