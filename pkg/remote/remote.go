package remote

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/nchern/notelog/pkg/env"
)

const (
	remoteConfigName = "remote"
)

var (
	schemeToCmd = map[string][]string{
		"rsync": []string{"rsync", "-r"},
	}
)

var errUnknownScheme = errors.New("Unknown scheme")

type entry struct {
	Scheme string
	Addr   string
}

func (e *entry) Push(name string) ([]string, error) {
	res := schemeToCmd[e.Scheme]
	if res == nil {
		return nil, errUnknownScheme
	}
	res = append(res, name, e.Addr)
	return res, nil
}

func Push() error {
	f, err := os.Open(env.NotesMetadataPath(remoteConfigName))
	if err != nil {
		return err
	}
	remotes, err := parse(f)
	if err != nil {
		return err
	}
	if len(remotes) < 1 {
		return fmt.Errorf("No remotes configured")
	}
	args, err := remotes[0].Push(env.NotesRootPath())
	if err != nil {
		return err
	}
	return run(args)
}

func run(args []string) error {
	log.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}
