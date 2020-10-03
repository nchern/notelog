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
	commands = commandList{}

	cmdLs           = c("list")
	cmdEdit         = c("edit")
	cmdSearch       = c("search")
	cmdPrint        = c("print")
	cmdPrintHome    = c("print-home")
	cmdGetFullPath  = c("path")
	cmdBashComplete = c("bash-complete")
	cmdSortTodoList = c("sort-todos")
	cmdEnv          = c("env")
	cmdRemotePush   = c("push")
	cmdRemotePull   = c("pull")

	// Command is a user subcommand
	Command = flag.String("c", cmdEdit, fmt.Sprintf("One of: %s", commands))
)

// Execute runs specified command
func Execute(cmd string) error {
	notes := note.NewList()

	switch cmd {
	case cmdEdit:
		return edit()
	case cmdLs:
		return listNotes()
	case cmdBashComplete:
		fmt.Println(autoCompleteScript())
	case cmdPrint:
		return printNote()
	case cmdPrintHome:
		fmt.Print(notes.HomeDir())
	case cmdGetFullPath:
		return printFullPath()
	case cmdSortTodoList:
		return todos.Sort(os.Stdin, os.Stdout)
	case cmdSearch:
		return search()
	case cmdEnv:
		fmt.Println(env.Vars())
	case cmdRemotePush:
		return handleNoRemoteConfig(remote.Push(notes))
	case cmdRemotePull:
		return handleNoRemoteConfig(remote.Pull(notes))
	default:
		return fmt.Errorf("Bad cmd: '%s'", cmd)
	}
	return nil
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
