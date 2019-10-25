package main

import (
	"os"
	"os/exec"
)

const defaultEditor = "vim"

func getEditorName() string {
	name := os.Getenv("EDITOR")
	if name == "" {
		return defaultEditor
	}
	return name
}

func editor(fileName string) *exec.Cmd {
	ed := getEditorName()

	cmd := exec.Command(ed, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
