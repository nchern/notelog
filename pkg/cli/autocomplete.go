package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

const cmdDo = "do"

var autocompleteCmd = &coral.Command{
	Use:   "autocomplete",
	Short: "uses by bash to return autocompletions",

	Args: coral.ArbitraryArgs,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *coral.Command, args []string) error {
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
	beforeCursor := line[0 : i+1]
	curTok := getCurrentCompletingToken(beforeCursor)
	prevToks := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(beforeCursor), curTok))

	if strings.HasSuffix(prevToks, cmdDo) {
		return printCommandsWithPrefix(curTok, w)
	}

	if strings.HasSuffix(prevToks, archOpenCmd.Use) {
		list = list.GetArchive()
	}
	notes, err := list.All()
	if err != nil {
		return err
	}
	return printNotesWithPrefix(notes, curTok, prevToks, w)
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

func autoCompleteDoCommand(curTok string, prevToks string, w io.Writer) error {
	// Hack: this hacky function attempts to autocomplete do command
	// if this is required
	for _, cmd := range doCmd.Commands() {
		// no need to autocomplete "do" if subcommands are already entered
		if strings.HasSuffix(prevToks, cmd.Use) {
			return nil
		}
	}
	if strings.HasPrefix(cmdDo, curTok) {
		_, err := fmt.Fprintln(w, cmdDo)
		if err != nil {
			return err
		}
	}
	return nil
}

func printNotesWithPrefix(notes []*note.Note, curTok string, prevToks string, w io.Writer) error {
	if err := autoCompleteDoCommand(curTok, prevToks, w); err != nil {
		return err
	}

	for _, note := range notes {
		if !strings.HasPrefix(note.Name(), curTok) {
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
