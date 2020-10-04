package note

import (
	"fmt"
	"os"
	"path/filepath"
)

// List represents a collection of notes.
// As of now this is a collection of folders in NOTELOG_HOME
// Each folder represents a note: main.org is a main note file
type List string

// HomeDir returns notes home dir
func (l List) HomeDir() string {
	return string(l)
}

// Note returns a node from the current collection with a given name
func (l List) Note(name string) *Note {
	return &Note{name: name, homeDir: l.HomeDir()}
}

// MetadataFilename returns full path to the notelog metadata for a given file
func (l List) MetadataFilename(name string) string {
	return filepath.Join(l.HomeDir(), ".notelog", name)
}

// Remove removes a note by name
func (l List) Remove(name string) error {
	filename := l.Note(name).Dir()
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", name)
		}
		return err
	}
	return os.RemoveAll(filename)
}

// NewList returns a list of notes
func NewList() List {
	return List(notesRootPath)
}
