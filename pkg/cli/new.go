package cli

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

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
		if err := createFromTemplate(nt, fromName); err != nil {
			return err
		}
	}
	return editor.Edit(nt, false, lnum)
}

func createFromTemplate(nt *note.Note, fromName string) error {
	tplSrc, err := notes.Get(fromName)
	if err != nil {
		return err
	}
	src, err := tplSrc.ReadAll()
	if err != nil {
		return err
	}
	tmpl, err := template.New("t").Funcs(
		template.FuncMap{
			"title": strings.Title,
			"upper": strings.ToUpper,
		}).Parse(src)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	vars := struct {
		Title string
		Date  string
	}{nt.Name(), time.Now().Format("2006-01-02")}
	if err := tmpl.Execute(&buf, &vars); err != nil {
		return err
	}
	return nt.Overwrite(buf.String())
}
