package note

import (
	"bytes"
	"os"
	"regexp"
	"strings"
)

const (
	recordTemplate = "%s"

	dateFormat = "2006-01-02"
)

// ShouldWriteFunc defines a function that determines if instant record should be written after a given line
type ShouldWriteFunc func(i uint, s string, prev string) bool

func expandVars(s string) string {
	mapping := func(s string) string {
		if s == "d" {
			return nowFn().Format(dateFormat)
		}
		return s
	}
	return os.Expand(s, mapping)
}

// WriteInstantRecord directly writes an `instant` string to a given note
func (n *Note) WriteInstantRecord(record string, skipLines uint, skipLinesAfter *regexp.Regexp) error {
	buf := bytes.Buffer{}
	err := n.Dump(&buf)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := []string{}
	if len(buf.Bytes()) > 0 {
		// strings.Split returns non empty array on empty input string
		lines = strings.Split(buf.String(), "\n")
	}

	out := []string{}

	if skipLinesAfter != nil {
		for i, l := range lines {
			if skipLinesAfter.MatchString(l) {
				skipLines = uint(i) + 1
				break
			}
		}
	}

	if skipLines >= uint(len(lines)) || skipLines == 0 {
		out = append([]string{expandVars(record), ""}, lines...)
	} else {
		tail := append([]string{expandVars(record), ""}, lines[skipLines:]...)
		out = append(lines[:skipLines], tail...)
	}

	body := strings.Join(out, "\n")
	return n.Write(body)
}
