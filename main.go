package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func fatal(s string) { log.Fatalf("FATAL: %s\n", s) }

func must(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func dieIf(err error) {
	must(err)
}

func init() {
	log.SetFlags(0)
}

type commandList []string

func (l commandList) String() string {
	return strings.Join(l, ", ")
}

func c(s string) string {
	commands = append(commands, s)
	return s
}

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

	cmd = flag.String("c", cmdEdit, fmt.Sprintf("One of: %s", commands))
)

func main() {
	flag.Parse()

	notes := note.NewList()

	switch *cmd {
	case cmdEdit:
		noteName, instantRecord, err := parseArgs(flag.Args())
		dieIf(err)
		must(editor.EditNote(notes.Note(noteName), instantRecord))
	case cmdLs:
		must(listNotes())
	case cmdBashComplete:
		fmt.Println(autoCompleteScript())
	case cmdPrint:
		noteName, _, err := parseArgs(flag.Args())
		dieIf(err)
		must(printNote(notes.Note(noteName)))
	case cmdPrintHome:
		fmt.Print(notes.HomeDir())
	case cmdGetFullPath:
		noteName, _, err := parseArgs(flag.Args())
		dieIf(err)
		must(printFullPath(notes.Note(noteName)))
	case cmdSortTodoList:
		must(todos.Sort(os.Stdin, os.Stdout))
	case cmdSearch:
		terms, err := parseSearchArgs(flag.Args())
		dieIf(err)
		must(searcher.Search(notes, terms))
	case cmdEnv:
		fmt.Println(env.Vars())
	case cmdRemotePush:
		must(handleNoRemoteConfig(remote.Push(notes)))
	case cmdRemotePull:
		must(handleNoRemoteConfig(remote.Pull(notes)))
	default:
		fatal(fmt.Sprintf("Bad cmd: '%s'", *cmd))
	}
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
