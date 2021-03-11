package note

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	return NewNote(name, l.HomeDir())
}

// MetadataFilename returns full path to the notelog metadata for a given file
func (l List) MetadataFilename(name string) string {
	return filepath.Join(l.HomeDir(), DotNotelogDir, name)
}

// Remove removes a note by name
func (l List) Remove(name string) error {
	path, err := getExistingNotePath(l.Note(name))
	if err != nil {
		return err
	}

	return os.RemoveAll(path)
}

// Rename renames a note
func (l List) Rename(oldName string, newName string) error {
	path, err := getExistingNotePath(l.Note(oldName))
	if err != nil {
		return err
	}

	return os.Rename(path, l.Note(newName).dir())
}

// All returns all notes from this list
func (l List) All() ([]*Note, error) {
	res := []*Note{}
	dirs, err := ioutil.ReadDir(l.HomeDir())
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}
		res = append(res, l.Note(dir.Name()))
	}

	return res, nil
}

// Archive puts a given note into archive
func (l List) Archive(name string) error {
	path, err := getExistingNotePath(l.Note(name))
	if err != nil {
		return err
	}

	archiveDir := filepath.Join(l.HomeDir(), archiveNoteDir)
	if err := os.MkdirAll(archiveDir, defaultDirPerms); err != nil {
		return err
	}
	// This does not work: os.Rename(path, archiveDir)
	// fails with "rename $path $archiveDir file exists"
	return exec.Command("mv", path, archiveDir).Run()
}

// Add creates a note in this list if it does not exist otherwise does nothing
func (l List) Add(name string) (*Note, error) {
	nt := NewNote(name, l.HomeDir())
	if err := nt.Init(); err != nil {
		return nil, err
	}
	if ok, _ := nt.Exists(); ok {
		return nt, nil
	}

	f, err := os.OpenFile(nt.FullPath(), os.O_RDWR|os.O_CREATE, defaultFilePerms)
	if err != nil {
		return nil, err
	}
	return nt, f.Close()
}

// NewList returns a list of notes
func NewList() List {
	return List(notesRootPath)
}

func getExistingNotePath(nt *Note) (string, error) {
	found, err := nt.Exists()
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("%s does not exist", nt.name)
	}
	return nt.dir(), nil
}
