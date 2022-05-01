package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	sortByDate bool

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "lists all notes. Output is sorted by note's name alphabetically by default",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listNotes(note.NewList(), os.Stdout)
		},
	}
)

type byDate []*note.Note

func (b byDate) Len() int           { return len(b) }
func (b byDate) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byDate) Less(i, j int) bool { return b[i].ModifiedAt().Before(b[j].ModifiedAt()) }

func init() {
	listCmd.Flags().BoolVarP(&sortByDate, "by-date", "d", false, "sorts notes by last modified date in asc order")
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output: includes note modified dates")

	doCmd.AddCommand(listCmd)
}

func formatNote(nt *note.Note) string {
	if verbose {
		return fmt.Sprintf("%s\t%s",
			nt.ModifiedAt().Format("2006-01-02T15:04"),
			nt.Name())
	}
	return nt.Name()
}

func listNotes(list note.List, w io.Writer) error {
	notes, err := list.All()
	if err != nil {
		return err
	}
	if sortByDate {
		sort.Sort(byDate(notes))
	}
	for _, note := range notes {
		if _, err := fmt.Fprintln(w, formatNote(note)); err != nil {
			return err
		}
	}
	return nil
}
