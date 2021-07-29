package note

import (
	"errors"
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

// Init creates initial resources required for this node collection
func (l List) Init() error {
	dirs := []string{
		l.HomeDir(),
		filepath.Join(l.HomeDir(), DotNotelogDir),
		filepath.Join(l.HomeDir(), archiveNoteDir),
	}
	for _, d := range dirs {
		err := os.Mkdir(d, defaultDirPerms)
		if errors.Is(err, os.ErrExist) {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Get returns an existing node from the current collection with a given name
// If the note with a given name does not exit an error is returned
func (l List) Get(name string) (*Note, error) {
	nt := newNote(name, l.HomeDir())
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

func (l List) metadataRoot() string {
	return filepath.Join(l.HomeDir(), DotNotelogDir)
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

	return os.Rename(nt.dir(), newNote(newName, l.HomeDir()).dir())
}

// Copy copies a note
func (l List) Copy(srcName string, dstName string) error {
	src, err := l.Get(srcName)
	if err != nil {
		return err
	}

	dst, err := l.GetOrCreate(dstName)
	if err != nil {
		return err
	}

	w, err := dst.writer()
	if err != nil {
		return err
	}
	defer w.Close()

	return src.Dump(w)
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
		nt := newNote(dir.Name(), l.HomeDir())
		// HACK: this works only as a whole note file gets re-created.
		// Vim does it when writes the file
		nt.modifiedAt = dir.ModTime()
		res = append(res, nt)
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

// GetOrCreate returns a note with a given name. If the note does not exist it creates it.
func (l List) GetOrCreate(name string) (*Note, error) {
	nt := newNote(name, l.HomeDir())

	// Init call ensures atomic note dir creation
	err := nt.Init()
	if errors.Is(err, os.ErrExist) {
		return nt, nil
	}
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(nt.FullPath(), os.O_RDWR|os.O_CREATE, defaultFilePerms)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return nt, nil
}

// NewList returns a list of notes
func NewList() List {
	return List(notesRootPath)
}
