package note

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sample = "instant"

func TestWriteInstantRecord(t *testing.T) {
	n := &Note{homeDir: t.TempDir(), name: testDir}
	if err := n.Init(); err != nil && !errors.Is(err, os.ErrExist) {
		require.NoError(t, err)
	}

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
			must(ioutil.WriteFile(n.FullPath(), []byte(initial), defaultFilePerms))
			defer os.Remove(n.FullPath())

			assert.NoError(t, n.WriteInstantRecord(sample, SkipLines(tt.givenSkipLines)))

			actual, err := ioutil.ReadFile(n.FullPath())

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(actual))
		})
	}
}

func TestWriteInstantRecordWithSkipLinesByRegex(t *testing.T) {
	n := &Note{homeDir: t.TempDir(), name: testDir}
	if err := n.Init(); err != nil && !errors.Is(err, os.ErrExist) {
		require.NoError(t, err)
	}

	initial := text(
		"foo",
		"fuu",
		"bar",
		"buzz")

	var tests = []struct {
		name                 string
		expected             string
		givenSkipLinesRegexp *regexp.Regexp
	}{
		{"should write after regexp match",
			text(
				"foo",
				sample,
				"",
				"fuu",
				"bar",
				"buzz"),
			regexp.MustCompile("foo"),
		},
		{"should write after regexp matching more than one line only after the first match",
			text(
				"foo",
				sample,
				"",
				"fuu",
				"bar",
				"buzz"),
			regexp.MustCompile("f(o|u)"),
		},
		{"should write after 1st occurance of regexp matching near the end",
			text(
				"foo",
				"fuu",
				"bar",
				sample,
				"",
				"buzz"),
			regexp.MustCompile("b.*?"),
		},
		{"should write to the end if no mateches",
			text(
				"foo",
				"fuu",
				"bar",
				"buzz",
				sample,
				""),
			regexp.MustCompile("abc"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			must(ioutil.WriteFile(n.FullPath(), []byte(initial), defaultFilePerms))
			defer os.Remove(n.FullPath())

			assert.NoError(t, n.WriteInstantRecord(sample, SkipLinesByRegex(tt.givenSkipLinesRegexp)))

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
