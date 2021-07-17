package main // numberedlist

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errRead  = errors.New("boom")
	errWrite = errors.New("boom boom")
)

func TestShouldFormat(t *testing.T) {

	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"empty",
			text(""),
			text(""),
		},
		{"simple not-yet numbered list with empty items",
			text(
				"1. foo",
				"2. bar",
				"",
				"\t",
				"3. fuzz buzz",
				"",
			),
			text(
				"foo",
				"bar",
				"",
				"\t",
				"fuzz buzz",
				"\n",
			),
		},
		{"list items have sub content",
			text(
				"1. foo",
				"   bazz abc",
				"2. bar",
				"   sub-item 1",
				"   sub-item 2",
				"   sub-item 3",
				"\t  ",
				"3. fuzz buzz",
				"   lala 1",
				"   lala 2",
				"4. foobar",
			),
			text(
				"foo",
				"   bazz abc",
				"bar",
				"   sub-item 1",
				"   sub-item 2",
				"   sub-item 3",
				"\t  ",
				"fuzz buzz",
				"   lala 1",
				"   lala 2",
				"foobar",
			),
		},
		{"already numbered list",
			text(
				"",
				"1. foo",
				"2. bar",
				" .  sub-item 1",
				"3. fuzz buzz",
				"4. foobar",
			),
			text(
				"",
				"1. foo",
				"3.bar",
				" .  sub-item 1",
				" 5. fuzz buzz",
				"\t6. foobar",
			),
		},
		{"numbered items with no text",
			text(
				"1.",
				"2. foo",
				" foobar",
				"3. bar",
				"\tbuzz",
			),
			text(
				"2.",
				"foo",
				" foobar",
				"4. bar",
				"\tbuzz",
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := &bytes.Buffer{}

			err := Format(bytes.NewBufferString(tt.given), actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected+"\n", actual.String())
		})
	}
}

func TestShouldFailOnReaderError(t *testing.T) {
	actual := &bytes.Buffer{}
	err := Format(&brokenStream{}, actual)
	assert.Equal(t, errRead, err)
}

func TestShouldFailOnWriterError(t *testing.T) {
	err := Format(bytes.NewBufferString("foo"), &brokenStream{})
	assert.Equal(t, errWrite, err)
}

func text(lines ...string) string {
	return strings.Join(lines, "\n")
}

type brokenStream struct{}

func (t *brokenStream) Read(b []byte) (int, error) { return 0, errRead }

func (t *brokenStream) Write(b []byte) (int, error) { return 0, errWrite }
