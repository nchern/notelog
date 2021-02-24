package editor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultEditor = "vim"
)

var (
	editorCmd = env.Get("EDITOR", defaultEditor)

	editorFlags = env.Get("EDITOR_FLAGS", "")
)

// Note abstracts the note to edit
type Note interface {
	Init() error
	FullPath() string
	RemoveIfEmpty() error
}

// Edit calls an editor to interactively edit given note
func Edit(note Note, readOnly bool) error {
	defer note.RemoveIfEmpty()

	if err := note.Init(); err != nil {
		return err
	}

	args := []string{}
	if readOnly {
		// TODO: customize, this works in vim only
		args = append(args, "-R")
	}
	args = append(args, note.FullPath())
	ed := Shellout(args...)
	return ed.Run()
}

// Shellout creates a ready to shellout exec.Command editor to edit given filename.
// It inherits all std* streams from the current process
func Shellout(flags ...string) *exec.Cmd {
	// HACK: this will not work properly if flags contain values with spaces
	args := strings.Fields(strings.TrimSpace(editorFlags))
	args = append(args, flags...)

	cmd := exec.Command(editorCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
