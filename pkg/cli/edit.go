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

func parseNoteNameAndLineNumber(rawName string) (name string, lnum editor.LineNumber) {
	nameAndLine := strings.SplitN(rawName, ":", 2)
	name = nameAndLine[0]
	if len(nameAndLine) > 1 {
		lnum = editor.LineNumber(nameAndLine[1])
	}
	return
}

func edit(args []string, readOnly bool) error {
	notes := note.NewList()

	var lnum editor.LineNumber
	noteName := noteNameFromArgs(args)
	noteName, lnum = parseNoteNameAndLineNumber(noteName)

	noteName, err := parseNoteName(noteName)
	if err != nil {
		return err
	}

	nt, err := notes.GetOrCreate(noteName)
	if err != nil {
		return err
	}

	instantRecord := ""
	if len(args) > 1 {
		instantRecord = strings.TrimSpace(strings.Join(args[1:], " "))
	}
	if instantRecord != "" {
		return nt.WriteInstantRecord(instantRecord, skipLines)
	}

	return editor.Edit(nt, readOnly, lnum)
}
