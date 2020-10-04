package note

import (
	"os"
	"path/filepath"

	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultNotesDir = "notes"
	defaultFilename = "main.org"
)

var notesRootPath = env.Get("NOTELOG_HOME", filepath.Join(os.Getenv("HOME"), defaultNotesDir))

// Note represents a note in the system. A directory with the main.org file as note file as of now.
type Note struct {
	name    string
	homeDir string
}

// FullPath returns full path to the notes file
func (n *Note) FullPath() string {
	return filepath.Join(n.homeDir, n.name, defaultFilename)
}

// Dir returns directory where note is stored
func (n *Note) Dir() string {
	return filepath.Join(n.homeDir, n.name)
}

// MetadataFilename returns full path to the metadata file in this note namespace
func (n *Note) MetadataFilename(name string) string {
	return filepath.Join(n.Dir(), name)
}

// Name returns a note's name
func (n *Note) Name() string {
	return n.name
}
