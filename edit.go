package main

import (
	"os"
	"os/exec"
)

func editor(fileName string) *exec.Cmd {
	ed := "nvim"

	cmd := exec.Command(ed, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
