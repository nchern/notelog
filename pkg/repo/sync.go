package repo

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/nchern/notelog/pkg/note"
)

// Sync syncs git repo in current $NOTELOG_HOME if the repo exists
func Sync(notes note.List, customMsg string, stderr io.Writer) error {
	msg := createMessage(customMsg)

	doCommit := fmt.Sprintf("%s diff-index --quiet HEAD || %s commit -q -m '%s'", gitCmd, gitCmd, msg)
	cmds := []*exec.Cmd{
		git(notes, stderr, "add", "."),
		sh(doCommit, notes.HomeDir(), stderr),
		git(notes, stderr, "pull", "-q", "--rebase"),
		git(notes, stderr, "push", "-q", "origin", "master"),
	}

	var err error
	for _, cmd := range cmds {
		fmt.Printf("notelog: calling %s\n", cmd)
		err = cmd.Run()

		switch v := err.(type) {
		case *exec.ExitError:
			fmt.Fprintf(os.Stderr, "notelog: [%s] returned code %d\n", cmd, v.ExitCode())
		}
	}

	return err
}

func createMessage(msg string) string {
	msg = strings.TrimSpace(msg)

	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}
	hostname := "unknown"
	if name, err := os.Hostname(); err == nil {
		hostname = name
	}

	res := fmt.Sprintf("notelog: sync called by %s@%s", username, hostname)

	if msg != "" {
		res += "; " + msg
	}
	return res
}
