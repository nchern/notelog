package repo

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/nchern/notelog/pkg/note"
)

// Init initialises new git repo in $NOTELOG_HOME
func Init(notes note.List, errStream io.Writer) error {
	cmd := git(notes, errStream, "init")

	if err := cmd.Run(); err != nil {
		return err
	}
	return createGitIgnore(notes)

}

func createGitIgnore(notes note.List) error {
	path := filepath.Join(notes.HomeDir(), gitIgnoreFile)
	return ioutil.WriteFile(path, []byte(note.DotNotelogDir), defaultFilePerms)
}
