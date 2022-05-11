package repo

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/nchern/notelog/pkg/note"
)

// Sync syncs git repo in current $NOTELOG_HOME if the repo exists
func Sync(notes note.List, customMsg string) error {
	msg := createMessage(customMsg)

	logName := notes.MetadataFilename(gitErrorLog)
	logFile, err := openErrorLog(logName)
	if err != nil {
		return err
	}
	defer logFile.Close()

	doCommit := fmt.Sprintf("%s diff-index --quiet HEAD || %s commit -q -m '%s'", gitCmd, gitCmd, msg)
	cmds := []*exec.Cmd{
		git(notes, logFile, "add", "."),
		sh(doCommit, notes.HomeDir(), logFile),
		git(notes, logFile, "pull", "-q", "--rebase"),
		git(notes, logFile, "push", "-q", "origin", "master"),
	}

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
