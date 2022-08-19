package note

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const (
	recordTemplate = "%s"
)

// ShouldWriteFunc defines a function that determines if instant record should be written after a given line
type ShouldWriteFunc func(i uint, s string, prev string) bool

type writer struct {
	err error
	f   *os.File
}

func (w *writer) println() error {
	if w.err != nil {
		return w.err
	}
	_, w.err = fmt.Fprintln(w.f)
	return w.err
}

func (w *writer) print(s string) error {
	if w.err != nil {
		return w.err
	}
	_, w.err = fmt.Fprint(w.f, s)
	return w.err
}

func (w *writer) writeRecord(record string) error {
	if w.err != nil {
		return w.err
	}
	_, w.err = fmt.Fprintf(w.f, recordTemplate, record)
	w.println()
	return w.err
}

// WriteInstantRecord directly writes an `instant` string to a given note
func (n *Note) WriteInstantRecord(record string, shouldWrite ShouldWriteFunc) error {
	filename := n.FullPath()

	srcFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, defaultFilePerms)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filename + ".t")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	i := uint(0)
	dst := &writer{f: dstFile}
	scanner := bufio.NewScanner(srcFile)

	prev := ""
	written := false
	for scanner.Scan() {
		s := scanner.Text()
		if i > 0 {
			dst.println()
		}
		if shouldWrite(i, s, prev) {
			dst.writeRecord(record)
			dst.println()
			written = true
		}
		dst.print(s)
		prev = s
		i++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !written {
		dst.println()
		dst.writeRecord(record)
	}
	if dst.err != nil {
		return dst.err
	}
	return os.Rename(dstFile.Name(), filename)
}

// SkipLines returns ShouldWriteFunc that tells to write instant record after N lines
func SkipLines(n uint) ShouldWriteFunc {
	return func(i uint, s string, prev string) bool {
		return i == n
	}
}

// SkipLinesByRegex returns ShouldWriteFunc that tells to write instant record after the line matching regexp
func SkipLinesByRegex(rx *regexp.Regexp) ShouldWriteFunc {
	return func(i uint, s string, prev string) bool {
		return rx.MatchString(prev)
	}
}
