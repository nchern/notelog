package repo

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func init() {
	// use mock cmd instead of real git
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gitCmd = filepath.Join(cwd, "./testdata/fake-git.sh")
}

func TestInit(t *testing.T) {
	testutil.WithNotes(map[string]string{}, func(notes note.List) {
		out := bytes.Buffer{}
		assert.NoError(t, Init(notes, &out))
		assert.Equal(t, "init\n", out.String())

		// check if gitignore is created and contains proper files to ignore
		gitIgnoreContent, err := ioutil.ReadFile(filepath.Join(notes.HomeDir(), gitIgnoreFile))
		assert.NoError(t, err)
		assert.Equal(t, note.DotNotelogDir, string(gitIgnoreContent))
	})
}
