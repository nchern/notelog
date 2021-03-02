package note

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteInstantRecord(t *testing.T) {
	const sample = "instant"
	n := &Note{homeDir: "/tmp", name: "test-note"}

	initial := text(
		"foo",
		"bar",
		"buzz")

	var tests = []struct {
		name           string
		expected       string
		givenSkipLines uint
	}{
		{"should write to the top",
			text(
				sample,
				"",
				"foo",
				"bar",
				"buzz"),
			0,
		},
		{"should skip one line",
			text(
				"foo",
				sample,
				"",
				"bar",
				"buzz"),
			1,
		},
		{"should skip two lines",
			text(
				"foo",
				"bar",
				sample,
				"",
				"buzz"),
			2,
		},
		{"should write to the end if asked to skip more lines that a file contains",
			text(
				"foo",
				"bar",
				"buzz",
				sample,
				""),
			20,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			must(n.Init())
			must(ioutil.WriteFile(n.FullPath(), []byte(initial), defaultFilePerms))
			defer os.Remove(n.FullPath())

			assert.NoError(t, n.WriteInstantRecord(sample, tt.givenSkipLines))

			actual, err := ioutil.ReadFile(n.FullPath())

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(actual))
		})
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func text(lines ...string) string { return strings.Join(lines, "\n") }