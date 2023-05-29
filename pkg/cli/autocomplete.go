package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

var autocompleteCmd = &coral.Command{
	Use:   "autocomplete",
	Short: "uses by bash to return autocompletions",

	Args: coral.ArbitraryArgs,

	SilenceErrors:      true,
	SilenceUsage:       true,
	DisableFlagParsing: true,

	RunE: func(cmd *coral.Command, args []string) error {
		pos, err := strconv.Atoi(os.Getenv("COMP_POINT"))
		if err != nil {
			return err
		}
		pos-- // bash sets position as 1- array based
		return autoComplete(note.NewList(), os.Getenv("COMP_LINE"), pos, cmd)
	},
}

func init() {
	rootCmd.AddCommand(autocompleteCmd)
}

func autoComplete(list note.List, line string, i int, cmd *coral.Command) error {
	w := cmd.OutOrStdout()
	beforeCursor := line[0 : i+1]
	curTok := getCurrentCompletingToken(beforeCursor)
	prevToks := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(beforeCursor), curTok))

	if strings.HasSuffix(prevToks, rootCmd.Use) {
		return printCommands(w, func(s string) bool { return strings.HasPrefix(s, curTok) })
	}

	// Autocomplete for subcommand
	return autoCompleteForSubcommand(cmd.Parent(), prevToks, curTok, w)
}

func autoCompleteForSubcommand(root *coral.Command, prevToks string, curTok string, w io.Writer) error {
	toks := strings.Split(prevToks, " ")
	if len(toks) < 2 {
		return nil
	}
	complCmdName := toks[1] // second command, the 1st is "notelog"
	complCmd, _, err := root.Find([]string{complCmdName})
	if err != nil {
		return err
	}
	if complCmd == nil || complCmd.ValidArgsFunction == nil {
		return nil
	}
	names, _ := complCmd.ValidArgsFunction(complCmd, []string{complCmdName}, curTok)
	for _, name := range names {
		if _, err := fmt.Fprintln(w, name); err != nil {
			return err
		}
	}
	return nil
}

func completeNoteNames(cmd *coral.Command, args []string, toComplete string) ([]string, coral.ShellCompDirective) {
	prevTok := ""
	list := note.NewList()
	if len(args) > 0 {
		prevTok = args[0]
	}
	if strings.HasSuffix(prevTok, archOpenCmdName) {
		list = list.GetArchive()
	}
	notes, err := list.All()
	if err != nil {
		// HACK: better to return the error, but API does not support it
		log.Fatalf("fatal: %T %s", err, err)
		return nil, coral.ShellCompDirectiveError
	}
	names := make([]string, 0, len(notes))
	for _, note := range notes {
		if !strings.HasPrefix(note.Name(), toComplete) {
			continue
		}
		names = append(names, note.Name())
	}
	return names, coral.ShellCompDirectiveDefault
}

func getCurrentCompletingToken(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ' ' {
			return s[i+1:]
		}
	}
	return s
}
