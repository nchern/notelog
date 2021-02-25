package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nchern/notelog/pkg/checklist"
	"github.com/nchern/notelog/pkg/env"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/repo"
)

const (
	subCommand     = "c"
	scratchpadName = ".scratchpad"
)

var (
	cmdArchive       = c("archive")
	cmdAutoComplete  = c("autocomplete")
	cmdBashComplete  = c("bash-complete")
	cmdEnv           = c("env")
	cmdEdit          = c("edit")
	cmdLs            = c("list")
	cmdLsCmds        = c("list-cmds")
	cmdGetFullPath   = c("path")
	cmdPrint         = c("print")
	cmdPrintHome     = c("print-home")
	cmdSync          = c("sync")
	cmdInitRepo      = c("init-repo")
	cmdRemove        = c("rm")
	cmdRename        = c("rename")
	cmdSearch        = c("search")
	cmdSearchBrowse  = c("search-browse")
	cmdSortCheckList = c("sort-checklist")
	cmdTouch         = c("touch")
	cmdVersion       = c("version")

	// Command is a user subcommand
	Command = flag.String(subCommand, cmdEdit, fmt.Sprintf("One of: %s", commands))

	readOnly = flag.Bool("r", false, "opens note in read-only mode")
)

// Execute runs specified command
func Execute(cmd string) error {
	var err error
	notes := note.NewList()

	switch cmd {
	case cmdAutoComplete:
		pos, err := strconv.Atoi(os.Getenv("COMP_POINT"))
		if err != nil {
			return err
		}
		pos-- // bash sets position as 1- array based
		return autoComplete(note.NewList(), os.Getenv("COMP_LINE"), pos, os.Stdout)
	case cmdEdit:
		return edit(*readOnly)
	case cmdLs:
		return listNotes(note.NewList(), os.Stdout)
	case cmdLsCmds:
		return listCommands(os.Stdout)
	case cmdBashComplete:
		_, err = fmt.Println(autoCompleteScript())
		return err
	case cmdPrint:
		return printNote()
	case cmdPrintHome:
		_, err = fmt.Print(notes.HomeDir())
		return err
	case cmdGetFullPath:
		return printFullPath()
	case cmdRemove:
		name, _, err := parseArgs(flag.Args())
		if err != nil {
			return err
		}
		return notes.Remove(name)
	case cmdRename:
		name, newName, err := parseArgs(flag.Args())
		if err != nil {
			return err
		}
		return notes.Rename(name, newName)
	case cmdSortCheckList:
		return checklist.Sort(os.Stdin, os.Stdout)
	case cmdSearch:
		return search()
	case cmdSearchBrowse:
		return browseSearch()
	case cmdEnv:
		_, err = fmt.Println(env.Vars())
		return err
	case cmdInitRepo:
		return repo.Init(notes)
	case cmdSync:
		return repo.Sync(notes, flag.Arg(0))
	case cmdTouch:
		return touch(notes)
	case cmdVersion:
		return printVersion()
	case cmdArchive:
		return archive()
	default:
		return fmt.Errorf("Bad cmd: '%s'", cmd)
	}
}

func parseArgs(args []string) (noteName string, instantRecord string, err error) {
	if len(args) < 1 {
		noteName = scratchpadName
		return
	}

	noteName = strings.TrimSpace(args[0])
	if err = validateNoteName(noteName); err != nil {
		return
	}
	instantRecord = strings.TrimSpace(strings.Join(args[1:], " "))
	return
}

func validateNoteName(name string) error {
	if name == "" {
		return errors.New("Empty note name. Specify the real name")
	}
	if strings.HasPrefix(name, ".") {
		return errors.New("Note name can not start with '.'")
	}
	return nil
}
