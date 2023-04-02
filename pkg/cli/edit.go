package cli

import (
	"io/ioutil"
	"os"
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
	readOnly  bool
	fromStdin bool

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
	editCmd.Flags().BoolVarP(&fromStdin,
		"stdin", "", false,
		"reads data from stdin and writes this data into the note")
	editCmd.Flags().BoolVarP(&readOnly,
		"read-only", "r", false,
		"opens note in read-only mode")
	addFormatFlag(editCmd)

	rootCmd.AddCommand(editCmd)
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

func writeRecord(nt *note.Note, record string) error {
	rx, err := conf.skipLinesRegex()
	if err != nil {
		return err
	}
	return nt.WriteInstantRecord(record, conf.SkipLines, rx)
}

func writeFromStdin(nt *note.Note) error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	return writeRecord(nt, string(b))
}

func edit(args []string, readOnly bool) error {
	t, err := note.ParseFormat(conf.NoteFormat)
	if err != nil {
		return err
	}

	var lnum editor.LineNumber
	noteName := noteNameFromArgs(args)
	noteName, lnum, _ = parseNoteNameAndLineNumber(noteName)

	noteName, err = parseNoteName(noteName)
	if err != nil {
		return err
	}

	notes := note.NewList()
	nt, err := notes.GetOrCreate(noteName, t)
	if err != nil {
		return err
	}

	if fromStdin {
		return writeFromStdin(nt)
	}

	if len(args) > 1 {
		instantRecord := strings.TrimSpace(strings.Join(args[1:], " "))
		if instantRecord != "" {
			return writeRecord(nt, instantRecord)
		}
	}

	return editor.Edit(nt, readOnly, lnum)
}
