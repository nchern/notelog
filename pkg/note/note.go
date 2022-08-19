package note

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultNotesDir = "notes"

	defaultFileBaseName = "main."

	// DotNotelogDir is a dir where notelog store its files
	DotNotelogDir = ".notelog"

	archiveNoteDir = ".archive"

	defaultDirPerms = 0700

	defaultFilePerms = 0644
)

var (
	notesRootPath = defineNotesRootPath()

	nowFn = func() time.Time { return time.Now() }
)

func curDirContainsNotes() bool {
	_, err := os.Stat(filepath.Join(".", DotNotelogDir))
	return err == nil
}

func defineNotesRootPath() string {
	defaultPath := filepath.Join(os.Getenv("HOME"), defaultNotesDir)
	if curDirContainsNotes() {
		path, err := os.Getwd()
		if err != nil {
			return defaultPath
		}
		return path
	}
	path := strings.TrimSpace(os.Getenv("NOTELOG_HOME"))
	if path == "" {
		return defaultPath
	}
	return path
}

// Note represents a note in the system. A directory with the main.org file as note file as of now.
type Note struct {
	typ     Format
	name    string
	homeDir string

	modifiedAt time.Time
}

// newNote creates a new instance of a Note
func newNote(name string, homeDir string) *Note {
	return &Note{
		typ:     Unknown,
		name:    name,
		homeDir: homeDir,
	}
}

// ModifiedAt returns time when this node was modified
func (n *Note) ModifiedAt() time.Time {
	return n.modifiedAt
}

// FullPath returns full path to the notes file
func (n *Note) FullPath() string {
	return filepath.Join(n.homeDir, n.name, defaultFileBaseName+string(n.typ))
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

// Size tells this node size in bytes
func (n *Note) Size() (int64, error) {
	st, err := os.Stat(n.FullPath())
	if err != nil {
		return 0, err
	}
	return st.Size(), nil
}

// RemoveIfEmpty cleans up note resources if the note is empty
func (n *Note) RemoveIfEmpty() error {
	l, err := n.Size()
	if err != nil {
		return err
	}
	if l > 0 {
		return nil
	}

	return os.RemoveAll(n.dir())
}

// Init initializes this note
func (n *Note) Init() error {
	return os.Mkdir(n.dir(), defaultDirPerms)
}

// Reader returns a reader to read this note content
func (n *Note) Reader() (io.ReadCloser, error) {
	return os.Open(n.FullPath())
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

// Archived determines if this note is archived
func (n *Note) Archived() bool {
	return strings.HasSuffix(n.homeDir, archiveNoteDir)
}

func (n *Note) writer() (io.WriteCloser, error) {
	return os.OpenFile(n.FullPath(), os.O_WRONLY, defaultFilePerms)
}

func (n *Note) dir() string {
	return filepath.Join(n.homeDir, n.name)
}

// NameFromFilename returns a note name from a given filename
func NameFromFilename(path string) string {
	return filepath.Base(filepath.Dir(path))
}
