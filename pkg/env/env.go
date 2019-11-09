package env

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultNotesDir = "notes"
	defaultFilename = "main.org"
)

var (
	settings = settingsBag{}
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
	settings[name] = defaultVal
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}

func VarNames() string {
	return settings.String()
}

type settingsBag map[string]string

func (b settingsBag) String() string {
	res := []string{}
	for k, _ := range b {
		res = append(res, k)
	}
	return strings.Join(res, "\n")
}
