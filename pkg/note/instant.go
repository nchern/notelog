package note

import (
	"bufio"
	"fmt"
	"os"
)

const (
	recordTemplate = "%s"
)

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
func (n *Note) WriteInstantRecord(record string, skipLines uint) error {
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

	for scanner.Scan() {
		if i > 0 {
			dst.println()
		}
		if i == skipLines {
			dst.writeRecord(record)
			dst.println()
		}
		if err := dst.print(scanner.Text()); err != nil {
			return err
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if i <= skipLines {
		dst.println()
		dst.writeRecord(record)
	}
	if dst.err != nil {
		return dst.err
	}
	return os.Rename(dstFile.Name(), filename)
}
