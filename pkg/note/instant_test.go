package note

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const sample = "instant"

func mkTestNotes(dir string) *Note {
	n := &Note{homeDir: dir, name: testDir}
	if err := n.Init(); err != nil && !errors.Is(err, os.ErrExist) {
		panic(err)
	}
	return n
}

func TestWriteInstantRecord(t *testing.T) {
	underTest := mkTestNotes(t.TempDir())

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
		{"should write at the beginning if asked to skip more lines that a file contains",
			text(
				sample,
				"",
				"foo",
				"bar",
				"buzz"),
			20,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			must(underTest.Overwrite(initial))
			defer os.Remove(underTest.FullPath())

			assert.NoError(t,
				underTest.WriteInstantRecord(sample, tt.givenSkipLines, nil))

			actual, err := os.ReadFile(underTest.FullPath())

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(actual))
		})
	}
}

func TestWriteInstantRecordWithSkipLinesByRegex(t *testing.T) {
	underTest := mkTestNotes(t.TempDir())

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
		{"should write at the beginning if no mateches",
			text(
				sample,
				"",
				"foo",
				"fuu",
				"bar",
				"buzz"),
			regexp.MustCompile("abc"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			must(underTest.Overwrite(initial))
			defer os.Remove(underTest.FullPath())

			assert.NoError(t,
				underTest.WriteInstantRecord(sample, 0, tt.givenSkipLinesRegexp))

			actual, err := os.ReadFile(underTest.FullPath())

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(actual))
		})
	}
}

func TestWriteInstantRecordShouldExpand(t *testing.T) {
	underTest := mkTestNotes(t.TempDir())

	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"date macro",
			fmt.Sprintf("%s foobar\n", testToday), "$d foobar"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(underTest.FullPath())

			assert.NoError(t,
				underTest.WriteInstantRecord(tt.given, 0, nil))

			actual, err := os.ReadFile(underTest.FullPath())

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
