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

// Edit calls an editor to interactively edit given note or directly writes an `instant` string to it
func Edit(note Note, instantRecord string) error {
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

// Touch creates given note if it does not exist otherwise does nothing
func Touch(note Note) error {
	if err := os.MkdirAll(note.Dir(), DefaultDirPerms); err != nil {
		return err
	}
	_, err := os.Stat(note.FullPath())
	if os.IsNotExist(err) {
		f, err := os.OpenFile(note.FullPath(), os.O_RDWR|os.O_CREATE, DefaultFilePerms)
		if err != nil {
			return err
		}
		return f.Close()
	}
	return nil
}

func removeDirIfNotesFileNotExists(dirName string, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		os.Remove(dirName)
	}
}
