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
	cmdGetFullPath  = c("get-full-path")
	cmdBashComplete = c("bash-complete")
	cmdSortTodoList = c("sort-todos")
	cmdListVars     = c("list-vars")
	cmdRemotePush   = c("push")
	cmdRemotePull   = c("pull")

	cmd = flag.String("cmd", cmdEdit, fmt.Sprintf("One of: %s", commands))
)

func main() {
	flag.Parse()

	switch *cmd {
	case cmdEdit:
		noteName, instantRecord, err := parseArgs(flag.Args())
		dieIf(err)
		must(editor.EditNote(noteName, instantRecord))
	case cmdLs:
		must(listNotes())
	case cmdBashComplete:
		fmt.Println(autoCompleteScript())
	case cmdPrint:
		noteName, _, err := parseArgs(flag.Args())
		dieIf(err)
		must(printNote(noteName))
	case cmdPrintHome:
		fmt.Print(env.NotesRootPath())
	case cmdGetFullPath:
		noteName, _, err := parseArgs(flag.Args())
		dieIf(err)
		fmt.Print(env.NotesFilePath(noteName))
	case cmdSortTodoList:
		must(todos.Sort(os.Stdin, os.Stdout))
	case cmdSearch:
		terms, err := parseSearchArgs(flag.Args())
		dieIf(err)
		must(searcher.Search(terms))
	case cmdListVars:
		fmt.Println(env.VarNames())
	case cmdRemotePush:
		must(handleNoRemoteConfig(remote.Push()))
	case cmdRemotePull:
		must(handleNoRemoteConfig(remote.Pull()))
	default:
		fatal(fmt.Sprintf("Bad cmd: '%s'", *cmd))
	}
}

func handleNoRemoteConfig(err error) error {
	if os.IsNotExist(err) {
		configPath := env.NotesMetadataPath(remote.ConfigName)
		if err := os.MkdirAll(path.Dir(configPath), editor.DefaultDirPerms); err != nil {
			return err
		}
		if err := ioutil.WriteFile(configPath, []byte(remote.DefaultConfig), editor.DefaultFilePerms); err != nil {
			return err
		}

		return editor.Command(configPath).Run()
	}
	return err
}

func listNotes() error {
	dirs, err := ioutil.ReadDir(env.NotesRootPath())
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir.Name() == scratchpadName {
			continue
		}
		fmt.Println(dir.Name())
	}
}

func parseArgs(args []string) (filename string, instantRecord string, err error) {
	if len(args) < 1 {
		filename = scratchpadName
		return
	}

	filename = args[0]
	if strings.HasPrefix(filename, ".") {
		err = errors.New("Note name can not start with '.'")
		return
	}
	instantRecord = strings.TrimSpace(strings.Join(args[1:], " "))
	return
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

func printNote(noteName string) error {
	filename := env.NotesFilePath(noteName)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, f)
	return err
}
