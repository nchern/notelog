package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/env"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/remote"
	"github.com/nchern/notelog/pkg/searcher"
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
		noteName, instantRecord, err := parseArgs(flag.Args())
		if err != nil {
			return err
		}
		if err := editor.Edit(notes.Note(noteName), instantRecord); err != nil {
			return err
		}
	case cmdLs:
		if err := listNotes(); err != nil {
			return err
		}
	case cmdBashComplete:
		fmt.Println(autoCompleteScript())
	case cmdPrint:
		noteName, _, err := parseArgs(flag.Args())
		if err != nil {
			return err
		}
		if err := printNote(notes.Note(noteName)); err != nil {
			return err
		}
	case cmdPrintHome:
		fmt.Print(notes.HomeDir())
	case cmdGetFullPath:
		noteName, _, err := parseArgs(flag.Args())
		if err != nil {
			return err
		}
		if err := printFullPath(notes.Note(noteName)); err != nil {
			return err
		}
	case cmdSortTodoList:
		if err := todos.Sort(os.Stdin, os.Stdout); err != nil {
			return err
		}
	case cmdSearch:
		terms, err := parseSearchArgs(flag.Args())
		if err != nil {
			return err
		}
		if err := searcher.Search(notes, terms); err != nil {
			return err
		}
	case cmdEnv:
		fmt.Println(env.Vars())
	case cmdRemotePush:
		if err := handleNoRemoteConfig(remote.Push(notes)); err != nil {
			return err
		}
	case cmdRemotePull:
		if err := handleNoRemoteConfig(remote.Pull(notes)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Bad cmd: '%s'", cmd)
	}
	return nil
}

func handleNoRemoteConfig(err error) error {
	if os.IsNotExist(err) {
		configPath := note.NewList().MetadataFilename(remote.ConfigName)
		if err := os.MkdirAll(path.Dir(configPath), editor.DefaultDirPerms); err != nil {
			return err
		}
		if err := ioutil.WriteFile(configPath, []byte(remote.DefaultConfig), editor.DefaultFilePerms); err != nil {
			return err
		}

		return editor.Shellout(configPath).Run()
	}
	return err
}

func listNotes() error {
	dirs, err := ioutil.ReadDir(note.NewList().HomeDir())
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if validateNoteName(dir.Name()) != nil {
			continue
		}
		fmt.Println(dir.Name())
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

func parseSearchArgs(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Not enough args. Specify a search term")
	}
	return args[0], nil
}

func autoCompleteScript() string {
	name := os.Args[0]
	return fmt.Sprintf("# Bash autocompletion for %s. Completes notes\ncomplete -W \"`%s -cmd=list`\" %s",
		name, name, name)
}

func printNote(n *note.Note) error {
	f, err := os.Open(n.FullPath())
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, f)
	return err
}

func printFullPath(n *note.Note) error {
	path := n.FullPath()
	if _, err := os.Stat(path); err != nil {
		return err
	}
	_, err := fmt.Print(path)
	return err
}
