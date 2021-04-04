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

// Get returns an existing node from the current collection with a given name
// If the note with a given name does not exit an error is returned
func (l List) Get(name string) (*Note, error) {
	nt := NewNote(name, l.HomeDir())
	found, err := nt.Exists()
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	return nt, nil
}

// MetadataFilename returns full path to the notelog metadata for a given file
func (l List) MetadataFilename(name string) string {
	return filepath.Join(l.HomeDir(), DotNotelogDir, name)
}

// Remove removes a note by name
func (l List) Remove(name string) error {
	nt, err := l.Get(name)
	if err != nil {
		return err
	}

	return os.RemoveAll(nt.dir())
}

// Rename renames a note
func (l List) Rename(oldName string, newName string) error {
	nt, err := l.Get(oldName)
	if err != nil {
		return err
	}

	return os.Rename(nt.dir(), NewNote(newName, l.HomeDir()).dir())
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
		res = append(res, NewNote(dir.Name(), l.HomeDir()))
	}

	return res, nil
}

// Archive puts a given note into archive
func (l List) Archive(name string) error {
	nt, err := l.Get(name)
	if err != nil {
		return err
	}

	archiveDir := filepath.Join(l.HomeDir(), archiveNoteDir)
	if err := os.MkdirAll(archiveDir, defaultDirPerms); err != nil {
		return err
	}
	// os.Rename(path, archiveDir) does not work:
	// it fails with "rename $path $archiveDir file exists"
	return exec.Command("mv", nt.dir(), archiveDir).Run()
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
