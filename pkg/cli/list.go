package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/nchern/notelog/pkg/note"
)

func listNotes() error {
	dirs, err := ioutil.ReadDir(note.NewList().HomeDir())
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if validateNoteName(dir.Name()) != nil {
			continue
		}
		fmt.Println(dir.Name())
	}
	return nil
}
