package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nchern/notelog/pkg/note"
)

const (
	mode = 0644
)

// WithNotes - test helper function
func WithNotes(files map[string]string, fn func(notes note.List)) {
	homeDir, err := ioutil.TempDir("", "test_notes")
	if err != nil {
		panic(err)
	}

	must(os.MkdirAll(homeDir, 0755))
	defer os.RemoveAll(homeDir)

	must(os.MkdirAll(filepath.Join(homeDir, note.DotNotelogDir), 0755))

	for name, body := range files {
		fullName := filepath.Join(homeDir, name, "main.org")
		dir, _ := filepath.Split(fullName)
		must(os.MkdirAll(dir, 0755))
		must(ioutil.WriteFile(fullName, []byte(body), mode))
	}

	fn(note.List(homeDir))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
