package cli

import (
	"errors"
	"fmt"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
)

var (
	newCmd = &coral.Command{
		Use:   "new",
		Short: "creates a new note. Fails if the note already exists",

		Args: coral.MinimumNArgs(1),

		SilenceErrors: true,
		SilenceUsage:  true,

		ValidArgsFunction: completeNoteNames,

		RunE: func(cmd *coral.Command, args []string) error {
			return newNote(args)
		},
	}

	fromName string
)

func init() {
	newCmd.Flags().StringVarP(&fromName,
		"from", "f", "",
		"create a new note from a given one as a template")
	addFormatFlag(newCmd)
	rootCmd.AddCommand(newCmd)
}

func newNote(args []string) error {
	t, err := note.ParseFormat(conf.NoteFormat)
	if err != nil {
		return err
	}
	rawName, err := parseNoteName(noteNameFromArgs(args))
	if err != nil {
		return err
	}
	noteName, lnum, _ := parseNoteNameAndLineNumber(rawName)
	notes := note.NewList()
	nt, err := notes.Get(noteName)
	if nt != nil {
		return fmt.Errorf("%s: already exists", noteName)
	}
	if !errors.Is(err, note.ErrNotExist) {
		return err
	}
	if nt, err = notes.GetOrCreate(noteName, t); err != nil {
		return err
	}
	if fromName != "" {
		if err := notes.Copy(fromName, noteName); err != nil {
			return err
		}
	}
	return editor.Edit(nt, false, lnum)
}
