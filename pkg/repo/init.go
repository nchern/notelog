package repo

import (
	"os"

	"github.com/nchern/notelog/pkg/note"
)

// Init initialises new git repo in $NOTELOG_HOME
func Init(notes note.List) error {
	cmd := git(notes, os.Stderr, "init")

	return cmd.Run()
}
