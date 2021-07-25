package cli

import (
	"strings"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

const (
	// should be configurable
	skipLines uint = 2
)

var (
	readOnly bool

	editCmd = &cobra.Command{
		Use:   "edit",
		Short: "opens a given note in editor",

		Args: cobra.MinimumNArgs(1),

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return edit(args, readOnly)
		},
	}
)

func init() {
	editCmd.Flags().BoolVarP(&readOnly, "read-only", "r", false, "opens note in read-only mode")

	doCmd.AddCommand(editCmd)
}

func edit(args []string, readOnly bool) error {
	notes := note.NewList()

	var lnum editor.LineNumber
	if len(args) > 0 {
		nameAndLine := strings.SplitN(args[0], ":", 2)
		args[0] = nameAndLine[0]
		if len(nameAndLine) > 1 {
			lnum = editor.LineNumber(nameAndLine[1])
		}
	}

	noteName, instantRecord, err := parseArgs(args)
	if err != nil {
		return err
	}

	nt, err := notes.GetOrCreate(noteName)
	if err != nil {
		return err
	}

	if instantRecord != "" {
		return nt.WriteInstantRecord(instantRecord, skipLines)
	}

	return editor.Edit(nt, readOnly, lnum)
}
