package editor

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteInstantRecordShouldWork(t *testing.T) {
	const sample = "instant"
	initial := text(
		"foo",
		"bar",
		"buzz")

	expected := text(
		sample,
		"",
		"foo",
		"bar",
		"buzz")

	n := &noteMock{}

	must(ioutil.WriteFile(n.FullPath(), []byte(initial), DefaultFilePerms))
	defer os.Remove(n.FullPath())

	assert.NoError(t, WriteInstantRecord(n, sample))

	actual, err := ioutil.ReadFile(n.FullPath())

	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func text(lines ...string) string { return strings.Join(lines, "\n") }

type noteMock struct{}

func (n *noteMock) Dir() string {
	return "/tmp/"
}

func (n *noteMock) FullPath() string {
	return "/tmp/test-note.txt"
}
