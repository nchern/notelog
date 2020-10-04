package cli

import (
	"fmt"
	"io"

	"github.com/nchern/notelog/pkg/note"
)

func listNotes(list note.List, w io.Writer) error {
	notes, err := list.All()
	if err != nil {
		return err
	}
	for _, note := range notes {
		if _, err := fmt.Fprintln(w, note.Name()); err != nil {
			return err
		}
	}
	return nil
}
