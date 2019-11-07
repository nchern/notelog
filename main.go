package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultDirPerms  = 0700
	defaultFilePerms = 0644
	defaultNotesDir  = "notes"
	defaultFilename  = "main.org"
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

func notesRootPath() string {
	return filepath.Join(os.Getenv("HOME"), defaultNotesDir)
}

func currentNotesFilePath(name string) string {
	return filepath.Join(notesRootPath(), name, defaultFilename)
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

	cmd = flag.String("cmd", cmdEdit, fmt.Sprintf("One of: %s", commands))
)

func removeDirIfNotesFileNotExists(dirName string, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		os.Remove(dirName)
	}
}

func edit() error {
	noteName, instantRecord, err := parseArgs(flag.Args())
	if err != nil {
		return err
	}

	filename := currentNotesFilePath(noteName)
	dirName := filepath.Dir(filename)

	defer removeDirIfNotesFileNotExists(dirName, filename)

	if err := os.MkdirAll(dirName, defaultDirPerms); err != nil {
		return err
	}

	if instantRecord != "" {
		return writeInstantRecord(filename, instantRecord)
	}

	ed := editor(filename)
	return ed.Run()
}

func main() {
	flag.Parse()

	if *cmd == cmdEdit {
		must(edit())
	} else if *cmd == cmdLs {
		dirs, err := ioutil.ReadDir(notesRootPath())
		dieOnError(err)
		for _, dir := range dirs {
			fmt.Println(dir.Name())
		}
	} else if *cmd == cmdBashComplete {
		fmt.Println(autoCompleteScript())
	} else if *cmd == cmdPrintHome {
		fmt.Print(notesRootPath())
	} else if *cmd == cmdGetFullPath {
		noteName, _, err := parseArgs(flag.Args())
		dieOnError(err)
		fmt.Print(currentNotesFilePath(noteName))
	} else if *cmd == cmdSortTodoList {
		dieOnError(sortTODOList(os.Stdin, os.Stdout))
	} else if *cmd == cmdSearch {
		terms, err := parseSearchArgs(flag.Args())
		dieOnError(err)
		must(search(terms))
	} else {
		fatal(fmt.Sprintf("Bad cmd: '%s'", *cmd))
	}
}
