package repo

import (
	"io"
	"os"
	"os/exec"

	"github.com/nchern/notelog/pkg/note"
)

const (
	defaultFilePerms = 0644

	gitIgnoreFile = ".gitignore"
	gitErrorLog   = "git-errors.log"
)

var gitCmd = "git"

func sh(expr string, cwd string, logFile io.Writer) *exec.Cmd {
	cmd := exec.Command("sh", "-c", expr)

	cmd.Stdout = os.Stdout
	cmd.Stderr = logFile
	cmd.Dir = cwd

	return cmd
}

func git(notes note.List, logFile io.Writer, args ...string) *exec.Cmd {
	cmd := exec.Command(gitCmd, args...)

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
