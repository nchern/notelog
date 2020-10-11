package repo

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/nchern/notelog/pkg/note"
)

const gitErrorLog = "git-errors.log"

func Sync(notes note.List) error {
	msg := "notelog: pre-sync update"

	logName := notes.MetadataFilename(gitErrorLog)
	logFile, err := openErrorLog(logName)
	if err != nil {
		return err
	}
	defer logFile.Close()

	cmds := []*exec.Cmd{
		git(notes, logFile, "add", "."),
		git(notes, logFile, "commit", "-q", "-m", msg),
		git(notes, logFile, "pull", "-q", "--rebase"),
		git(notes, logFile, "push", "-q", "origin", "master"),
	}

	for _, cmd := range cmds {
		fmt.Printf("notelog: calling %s\n", cmd)
		err = cmd.Run()

		switch v := err.(type) {
		case *exec.ExitError:
			fmt.Println(v.ExitCode())
		}
	}

	return err
}

func git(notes note.List, logFile io.Writer, args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = logFile
	cmd.Dir = notes.HomeDir()

	return cmd
}

func openErrorLog(logName string) (*os.File, error) {
	if _, err := os.Stat(logName); err != nil {
		if os.IsNotExist(err) {
			return os.Create(logName)
		}
		return nil, err
	}

	// file exists
	return os.OpenFile(logName, os.O_WRONLY|os.O_APPEND, 0666)
}
