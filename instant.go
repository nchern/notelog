package main

import (
	"fmt"
	"io"
	"os"
)

func writeInstantRecord(filename string, instantRecord string) error {
	dstFileName := filename + ".t"
	dstFile, err := os.Create(dstFileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := fmt.Fprintf(dstFile, " - %s\n\n", instantRecord); err != nil {
		return err
	}

	srcFile, err := os.Open(filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		srcFile.Close()
		return err
	}
	srcFile.Close()

	return os.Rename(dstFileName, filename)
}
