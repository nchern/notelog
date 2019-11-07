package main

import (
	"errors"
	"os"
	"os/exec"
)

const defaultGrep = "grep"
const defaultGrepArgs = "-rni"

func getGrepCmd() string {
	name := os.Getenv("NOTELOG_GREP")
	if name == "" {
		return defaultGrep
	}
	return name
}

func parseSearchArgs(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Not enough args. Specify a search term")
	}
	return args[0], nil
}

func search(terms string) error {

	cmd := exec.Command(getGrepCmd(), defaultGrepArgs, terms, notesRootPath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	switch err.(type) {
	case *exec.ExitError:
		// TODO: check for error code 1
		os.Exit(1)
	default:
		return err
	}

	return nil
}
