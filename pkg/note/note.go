package note

import (
	"io"
	"os"
	"path/filepath"

	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultNotesDir = "notes"
	defaultFilename = "main.org"

	// DotNotelogDir is a dir where notelog store its files
	DotNotelogDir = ".notelog"

	archiveNoteDir = ".archive"

	defaultDirPerms = 0700

	defaultFilePerms = 0644
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

// MetadataFilename returns full path to the metadata file in this note namespace
func (n *Note) MetadataFilename(name string) string {
	return filepath.Join(n.dir(), name)
}

// Name returns a note's name
func (n *Note) Name() string {
	return n.name
}

// Exists tells if this note exists
func (n *Note) Exists() (bool, error) {
	_, err := os.Stat(n.FullPath())
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// RemoveIfEmpty cleans up note resources if the note is empty
func (n *Note) RemoveIfEmpty() error {
	if ok, _ := n.Exists(); ok {
		return nil
	}

	return os.Remove(n.dir())
}

// Init initializes this note
func (n *Note) Init() error {
	return os.MkdirAll(n.dir(), defaultDirPerms)
}

// Touch creates given note if it does not exist otherwise does nothing
func (n *Note) Touch() error {
	if err := n.Init(); err != nil {
		return err
	}
	if ok, _ := n.Exists(); ok {
		return nil
	}

	f, err := os.OpenFile(n.FullPath(), os.O_RDWR|os.O_CREATE, defaultFilePerms)
	if err != nil {
		return err
	}
	return f.Close()
}

// Dump writes this note contents to a given writer
func (n *Note) Dump(w io.Writer) error {
	f, err := os.Open(n.FullPath())
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}

func (n *Note) dir() string {
	return filepath.Join(n.homeDir, n.name)
}
