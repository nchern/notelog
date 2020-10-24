package repo

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/nchern/notelog/pkg/note"
)

// Sync syncs git repo in current $NOTELOG_HOME if the repo exists
func Sync(notes note.List) error {
	msg := createMessage()

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
			fmt.Printf("notelog: [%s] returned code %d\n", cmd, v.ExitCode())
		}
	}

	return err
}

func createMessage() string {
	username := "unknown"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}
	hostname := "unknown"
	if name, err := os.Hostname(); err == nil {
		hostname = name
	}

	return fmt.Sprintf("notelog: sync called by %s@%s", username, hostname)
}
