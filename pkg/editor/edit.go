package editor

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultDirPerms  = 0700
	defaultFilePerms = 0644

	defaultEditor = "vim"
)

// Edit note: calls editor or writes instant
func Edit(noteName string, instantRecord string) error {

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

func removeDirIfNotesFileNotExists(dirName string, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		os.Remove(dirName)
	}
}
