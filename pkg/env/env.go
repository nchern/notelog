package env

import (
	"os"
	"path/filepath"
)

const (
	defaultNotesDir = "notes"
	defaultFilename = "main.org"
)

// NotesRootPath returns notes home dir
func NotesRootPath() string {
	return filepath.Join(os.Getenv("HOME"), defaultNotesDir)
}

// NotesFilePath returns full path to the notes file
func NotesFilePath(name string) string {
	return filepath.Join(NotesRootPath(), name, defaultFilename)
}

// Get returns value of env var with given name. If it's empty, returns defaultVal
func Get(name string, defaultVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}
