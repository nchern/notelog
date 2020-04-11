package env

import (
	"fmt"
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

	notesRootPath = Get("NOTELOG_HOME", filepath.Join(os.Getenv("HOME"), defaultNotesDir))
)

// NotesRootPath returns notes home dir
func NotesRootPath() string {
	return notesRootPath
}

// NotesFilePath returns full path to the notes file
func NotesFilePath(name string) string {
	return filepath.Join(NotesRootPath(), name, defaultFilename)
}

// NotesMetadataPath returns full path to the notelog metadata for a given file
func NotesMetadataPath(name string) string {
	return filepath.Join(NotesRootPath(), ".notelog", name)
}

// Get returns value of env var with given name. If it's empty, returns defaultVal
func Get(name string, defaultVal string) string {
	res := defaultVal

	defer func() { settings[name] = res }()

	val := os.Getenv(name)
	if val == "" {
		return res
	}
	res = val
	return res
}

// Vars returns string dump of all env variables along with their values
func Vars() string {
	return settings.String()
}

type settingsBag map[string]string

func (b settingsBag) String() string {
	res := []string{}
	for k, v := range b {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(res, "\n")
}
