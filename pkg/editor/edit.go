package editor

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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

// EditNote calls an editor to interactively edit `noteName` or directly writes an `instant` string to it
func EditNote(noteName string, instantRecord string) error {
	noteName = strings.TrimSpace(noteName)
	if noteName == "" {
		return errors.New("Empty note name. Specify the real name")
	}

	filename := env.NotesFilePath(noteName)
	dirName := filepath.Dir(filename)

	defer removeDirIfNotesFileNotExists(dirName, filename)

	if err := os.MkdirAll(dirName, DefaultDirPerms); err != nil {
		return err
	}

	if instantRecord != "" {
		return writeInstantRecord(filename, instantRecord)
	}

	ed := Command(filename)
	return ed.Run()
}

// Command creates exec.Command with editor to edit given filename
func Command(fileName string) *exec.Cmd {
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
