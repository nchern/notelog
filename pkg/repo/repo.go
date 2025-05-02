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
)

var gitCmd = "git"

func sh(expr string, cwd string, stderr io.Writer) *exec.Cmd {
	cmd := exec.Command("sh", "-c", expr)

	cmd.Stdout = os.Stdout
	cmd.Stderr = stderr
	cmd.Dir = cwd

	return cmd
}

func git(notes note.List, stderr io.Writer, args ...string) *exec.Cmd {
	cmd := exec.Command(gitCmd, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = stderr
	cmd.Dir = notes.HomeDir()

	return cmd
}
