package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/env"
	"github.com/nchern/notelog/pkg/searcher"
	"github.com/nchern/notelog/pkg/todos"
)

func fatal(s string) { log.Fatalf("FATAL: %s\n", s) }

func must(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func dieOnError(err error) {
	must(err)
}

func init() {
	log.SetFlags(0)
}

func parseArgs(args []string) (filename string, instantRecord string, err error) {
	if len(args) < 1 {
		err = errors.New("Not enough args. Specify at least notes file name")
		return
	}

	filename = args[0]
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
	cmdPrintHome    = c("print-home")
	cmdGetFullPath  = c("get-full-path")
	cmdBashComplete = c("bash-complete")
	cmdSortTodoList = c("sort-todos")
	cmdListVars     = c("list-vars")

	cmd = flag.String("cmd", cmdEdit, fmt.Sprintf("One of: %s", commands))
)

func main() {
	flag.Parse()

	switch *cmd {
	case cmdEdit:
		noteName, instantRecord, err := parseArgs(flag.Args())
		dieOnError(err)
		must(editor.Edit(noteName, instantRecord))
	case cmdLs:
		dirs, err := ioutil.ReadDir(env.NotesRootPath())
		dieOnError(err)
		for _, dir := range dirs {
			fmt.Println(dir.Name())
		}
	case cmdBashComplete:
		fmt.Println(autoCompleteScript())
	case cmdPrintHome:
		fmt.Print(env.NotesRootPath())
	case cmdGetFullPath:
		noteName, _, err := parseArgs(flag.Args())
		dieOnError(err)
		fmt.Print(env.NotesFilePath(noteName))
	case cmdSortTodoList:
		dieOnError(todos.Sort(os.Stdin, os.Stdout))
	case cmdSearch:
		terms, err := parseSearchArgs(flag.Args())
		dieOnError(err)
		must(searcher.Search(terms))
	case cmdListVars:
		fmt.Println(env.VarNames())
	default:
		fatal(fmt.Sprintf("Bad cmd: '%s'", *cmd))
	}
}
