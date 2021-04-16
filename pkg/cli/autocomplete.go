package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/nchern/notelog/pkg/note"
	"github.com/spf13/cobra"
)

var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "uses by bash to return autocompletions",

	Args: cobra.ArbitraryArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		pos, err := strconv.Atoi(os.Getenv("COMP_POINT"))
		if err != nil {
			return err
		}
		pos-- // bash sets position as 1- array based
		return autoComplete(note.NewList(), os.Getenv("COMP_LINE"), pos, os.Stdout)
	},
}

func init() {
	doCmd.AddCommand(autocompleteCmd)
}

func autoComplete(list note.List, line string, i int, w io.Writer) error {
	const cmdDo = "do"

	beforeCursor := line[0 : i+1]
	curTok := getCurrentCompletingToken(beforeCursor)
	prevToks := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(beforeCursor), curTok))
	if strings.HasPrefix(curTok, "d") {
		_, err := fmt.Fprintln(w, cmdDo)
		return err
	}

	if strings.HasSuffix(prevToks, cmdDo) {
		return printCommandsWithPrefix(curTok, w)
	}

	notes, err := list.All()
	if err != nil {
		return err
	}
	return printNotesWithPrefix(notes, curTok, w)
}

func printCommandsWithPrefix(prefix string, w io.Writer) error {
	for _, c := range doCmd.Commands() {
		if !strings.HasPrefix(c.Use, prefix) {
			continue
		}
		if _, err := fmt.Fprintln(w, c.Use); err != nil {
			return err
		}
	}
	return nil
}

func printNotesWithPrefix(notes []*note.Note, prefix string, w io.Writer) error {
	for _, note := range notes {
		if !strings.HasPrefix(note.Name(), prefix) {
			continue
		}
		if _, err := fmt.Fprintln(w, note.Name()); err != nil {
			return err
		}
	}
	return nil
}

func getCurrentCompletingToken(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ' ' {
			return s[i+1:]
		}
	}
	return s
}
