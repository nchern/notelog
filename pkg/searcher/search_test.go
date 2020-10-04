package searcher

import (
	"bytes"
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

func TestSearchShoudWriteLastSearchResults(t *testing.T) {
	withFiles(func() {
		actual := &bytes.Buffer{}

		underTest := NewSearcher(&mock{}, actual)
		underTest.SaveResults = true

		assert.NoError(t, underTest.Search("foobar"))

		resultsFilename := filepath.Join(homeDir, lastResultsFile)
		body, err := ioutil.ReadFile(resultsFilename)

		assert.NoError(t, err) // file must exist
		assert.Equal(t, actual.Bytes(), body)
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

func TestGetLastNthResult(t *testing.T) {
	withFiles(func() {
		m := &mock{}
		buf := &bytes.Buffer{}
		underTest := NewSearcher(m, buf)
		underTest.SaveResults = true

		err := underTest.Search("foo")

		assert.NoError(t, err)
		actual, err := GetLastNthResult(m, 1)

		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(actual, filepath.Join(homeDir, "b.txt")))
	})
}

type mock struct{}

func (m *mock) HomeDir() string { return homeDir }

func (m *mock) MetadataFilename(name string) string { return filepath.Join(homeDir, name) }

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
	return strings.Split(strings.Trim(s, "\n"), "\n")
}
