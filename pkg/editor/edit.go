package editor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/nchern/notelog/pkg/env"
)

const (
	// DefaultDirPerms
	DefaultDirPerms = 0700

	// DefaultFilePerms
	DefaultFilePerms = 0644

	defaultEditor = "vim"
)

var (
	editorCmd = env.Get("EDITOR", defaultEditor)

	editorFlags = env.Get("EDITOR_FLAGS", "")
)

// Note abstracts the note to edit
type Note interface {
	Dir() string
	FullPath() string
}

// EditNote calls an editor to interactively edit `noteName` or directly writes an `instant` string to it
func EditNote(note Note, instantRecord string) error {
	defer removeDirIfNotesFileNotExists(note.Dir(), note.FullPath())

	if err := os.MkdirAll(note.Dir(), DefaultDirPerms); err != nil {
		return err
	}

	if instantRecord != "" {
		return writeInstantRecord(note.FullPath(), instantRecord)
	}

	ed := Shellout(note.FullPath())
	return ed.Run()
}

// Shellout creates a ready to shellout exec.Command editor to edit given filename.
// It inherits all std* streams from the current process
func Shellout(fileName string) *exec.Cmd {
	// HACK: this will not work properly if flags contain values with spaces
	args := strings.Fields(strings.TrimSpace(editorFlags))
	args = append(args, fileName)

	cmd := exec.Command(editorCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func removeDirIfNotesFileNotExists(dirName string, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		os.Remove(dirName)
	}
}
