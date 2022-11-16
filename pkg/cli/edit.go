package cli

import (
	"regexp"
	"strings"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
)

const (
	defaultSkipLines uint = 2

	defaultFormat = string(note.Org)
)

var (
	readOnly bool

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
		"read-only", "r", false,
		"opens note in read-only mode")
	addFormatFlag(editCmd)

	doCmd.AddCommand(editCmd)
}

func addFormatFlag(cmd *coral.Command) {
	cmd.Flags().StringVarP(&conf.NoteFormat,
		"format", "t", defaultFormat,
		"note format; currently org or md are supported")
	cmd.Flags().StringVarP(&conf.SkipLinesAfterMatch,
		"skip-lines-after", "s", "",
		"sets regexp to define where to write instant records")
}

func parseNoteNameAndLineNumber(rawName string) (name string, lnum editor.LineNumber, isArchive bool) {
	nameAndLine := strings.Split(rawName, ":")
	name = nameAndLine[0]
	isArchive = false
	if len(nameAndLine) > 1 {
		lnum = editor.LineNumber(nameAndLine[1])
	}
	if len(nameAndLine) > 2 && nameAndLine[2] == "a" {
		isArchive = true
	}
	return
}

func edit(args []string, readOnly bool) error {
	t, err := note.ParseFormat(conf.NoteFormat)
	if err != nil {
		return err
	}
	notes := note.NewList()

	var lnum editor.LineNumber
	noteName := noteNameFromArgs(args)
	noteName, lnum, _ = parseNoteNameAndLineNumber(noteName)

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
		var rx *regexp.Regexp
		if conf.SkipLinesAfterMatch != "" {
			rx, err = regexp.Compile(conf.SkipLinesAfterMatch)
			if err != nil {
				return err
			}
		}
		return nt.WriteInstantRecord(instantRecord, conf.SkipLines, rx)
	}

	return editor.Edit(nt, readOnly, lnum)
}
