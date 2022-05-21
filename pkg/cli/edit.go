package cli

import (
	"strings"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
)

const (
	// should be configurable
	defaultSkipLines uint = 2

	defaultFormat = string(note.Org)
)

var (
	readOnly bool

	noteFormat string

	editCmd = &coral.Command{
		Use:   "edit",
		Short: "opens a given note in editor",

		Args: coral.MinimumNArgs(1),

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *coral.Command, args []string) error {
			return edit(args, readOnly)
		},
	}
)

func init() {
	editCmd.Flags().BoolVarP(&readOnly,
		"read-only", "r", false, "opens note in read-only mode")
	addFormatFlag(editCmd)

	doCmd.AddCommand(editCmd)
}

func addFormatFlag(cmd *coral.Command) {
	cmd.Flags().StringVarP(&noteFormat,
		"format", "t", defaultFormat, "note format; currently org or md are supported")
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
	t, err := note.ParseFormat(noteFormat)
	if err != nil {
		return err
	}
	notes := note.NewList()

	var lnum editor.LineNumber
	noteName := noteNameFromArgs(args)
	noteName, lnum = parseNoteNameAndLineNumber(noteName)

	noteName, err = parseNoteName(noteName)
	if err != nil {
		return err
	}
	nt, err := notes.GetOrCreate(noteName, t)
	if err != nil {
		return err
	}

	instantRecord := ""
	if len(args) > 1 {
		instantRecord = strings.TrimSpace(strings.Join(args[1:], " "))
	}
	if instantRecord != "" {
		return nt.WriteInstantRecord(instantRecord, defaultSkipLines)
	}

	return editor.Edit(nt, readOnly, lnum)
}
