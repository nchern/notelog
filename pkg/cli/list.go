package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "lists all notes",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listNotes(note.NewList(), os.Stdout)
	},
}

func init() {
	doCmd.AddCommand(listCmd)
}

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
