package remote

import (
	"os"

	"github.com/nchern/notelog/pkg/env"
)

func Pull() error {
	f, err := os.Open(env.NotesMetadataPath(ConfigName))
	if err != nil {
		return err
	}
	remotes, err := parse(f)
	if err != nil {
		return err
	}
	args, err := remotes[0].Pull(env.NotesRootPath())
	if err != nil {
		return err
	}
	return run(args)
}
