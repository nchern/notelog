package editor

import (
	"bufio"
	"fmt"
	"os"
)

const (
	recordTemplate = "%s"
)

// WriteInstantRecord directly writes an `instant` string to a given note
func WriteInstantRecord(note Note, record string, skipLines uint) error {
	filename := note.FullPath()
	srcFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, DefaultFilePerms)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFileName := filename + ".t"
	dstFile, err := os.Create(dstFileName)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	i := uint(0)
	scanner := bufio.NewScanner(srcFile)
	for scanner.Scan() {
		if i > 0 {
			if _, err := fmt.Fprintln(dstFile); err != nil {
				return err
			}
		}
		if i == skipLines {
			if _, err := fmt.Fprintf(dstFile, recordTemplate+"\n\n", record); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(dstFile, scanner.Text()); err != nil {
			return err
		}
		i++
	}
	if i <= skipLines {
		if _, err := fmt.Fprintf(dstFile, "\n"+recordTemplate+"\n", record); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return os.Rename(dstFileName, filename)
}
