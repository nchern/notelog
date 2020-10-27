package editor

import (
	"fmt"
	"io"
	"os"
)

const recordTemplate = "%s"

// WriteInstantRecord directly writes an `instant` string to a given note
func WriteInstantRecord(note Note, record string) error {
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

	if _, err := fmt.Fprintf(dstFile, recordTemplate+"\n\n", record); err != nil {
		return err
	}

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err := dstFile.Sync(); err != nil {
		return err
	}

	return os.Rename(dstFileName, filename)
}
