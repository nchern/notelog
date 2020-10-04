package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nchern/notelog/pkg/env"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/remote"
	"github.com/nchern/notelog/pkg/todos"
)

const scratchpadName = ".scratchpad"

var (
	cmdBashComplete = c("bash-complete")
	cmdEnv          = c("env")
	cmdEdit         = c("edit")
	cmdLs           = c("list")
	cmdGetFullPath  = c("path")
	cmdPrint        = c("print")
	cmdPrintHome    = c("print-home")
	cmdRemotePull   = c("pull")
	cmdRemotePush   = c("push")
	cmdRemove       = c("rm")
	cmdRename       = c("rename")
	cmdSearch       = c("search")
	cmdSearchBrowse = c("search-browse")
	cmdSortTodoList = c("sort-todos")

	// Command is a user subcommand
	Command = flag.String("c", cmdEdit, fmt.Sprintf("One of: %s", commands))
)

// Execute runs specified command
func Execute(cmd string) error {
	var err error
	notes := note.NewList()

	switch cmd {
	case cmdEdit:
		return edit()
	case cmdLs:
		return listNotes()
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
	case cmdSortTodoList:
		return todos.Sort(os.Stdin, os.Stdout)
	case cmdSearch:
		return search()
	case cmdSearchBrowse:
		return browseSearch()
	case cmdEnv:
		_, err = fmt.Println(env.Vars())
		return err
	case cmdRemotePush:
		return handleNoRemoteConfig(remote.Push(notes))
	case cmdRemotePull:
		return handleNoRemoteConfig(remote.Pull(notes))
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
