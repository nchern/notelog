package editor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nchern/notelog/pkg/env"
	"github.com/pkg/errors"
)

const (
	defaultDirPerms  = 0700
	defaultFilePerms = 0644

	defaultEditor = "vim"
)

var (
	editorCmd = env.Get("EDITOR", defaultEditor)
)

// Edit note: calls editor or writes instant
func Edit(noteName string, instantRecord string) error {
	noteName = strings.TrimSpace(noteName)
	if noteName == "" {
		return errors.New("Empty note name. Specify the real name")
	}

	filename := env.NotesFilePath(noteName)
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

func editor(fileName string) *exec.Cmd {
	cmd := exec.Command(editorCmd, fileName)
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
