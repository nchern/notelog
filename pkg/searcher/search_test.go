package searcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	homeDir = "/tmp/test_notes"

	mode = 0644
)

var (
	files = map[string]string{
		"a.txt": "foo bar buzz",
		"b.txt": "foobar bar addd buzz",
		"c.txt": "fuzz bar xx buzz",
	}
)

func init() {
	grepCmd = defaultGrep // make sure we always use defaultGrep in tests
}

func TestShoudSearch(t *testing.T) {
	withFiles(func() {
		var tests = []struct {
			name     string
			expected int
			given    []string
		}{
			{"one term",
				2, []string{"foo"}},
			{"with excluded terms",
				1, []string{"bar", "-fuzz", "-foobar"}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual := &bytes.Buffer{}
				underTest := NewSearcher(&mock{}, actual)
				assert.NoError(t, underTest.Search(tt.given...))
				assert.Equal(t, tt.expected, len(toLines(actual.String())))
			})
		}
	})
}

func TestSearchShoudReturnOneIfFoundNothing(t *testing.T) {
	withFiles(func() {
		actual := &bytes.Buffer{}
		underTest := NewSearcher(&mock{}, actual)
		err := underTest.Search("you will not find me")

		assert.NotNil(t, err)
		assert.Equal(t, 1, (err.(*exec.ExitError)).ExitCode())
	})
}

type mock struct{}

func (m *mock) HomeDir() string { return homeDir }

func (m *mock) MetadataFilename(_ string) string { return homeDir }

func withFiles(fn func()) {
	must(os.MkdirAll(homeDir, 0755))
	defer os.RemoveAll(homeDir)

	for name, body := range files {
		must(ioutil.WriteFile(filepath.Join(homeDir, name), []byte(body), mode))
	}

	fn()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func toLines(s string) []string {
	fmt.Println(s)
	return strings.Split(strings.Trim(s, "\n"), "\n")
}
