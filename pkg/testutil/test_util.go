package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nchern/notelog/pkg/note"
)

const (
	mode    = 0644
	dirMode = 0755
)

// WithNotes - test helper function
func WithNotes(files map[string]string, fn func(notes note.List)) {
	homeDir, err := ioutil.TempDir("", "test_notes")
	if err != nil {
		panic(err)
	}

	must(os.MkdirAll(homeDir, dirMode))
	defer os.RemoveAll(homeDir)

	must(os.MkdirAll(filepath.Join(homeDir, note.DotNotelogDir), dirMode))

	for name, body := range files {
		fullName := filepath.Join(homeDir, name, "main.org")
		dir, _ := filepath.Split(fullName)
		must(os.MkdirAll(dir, dirMode))
		must(ioutil.WriteFile(fullName, []byte(body), mode))
	}

	fn(note.List(homeDir))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
