package remote

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nchern/notelog/pkg/env"
)

const (
	ConfigName = "remote"
)

var (
	errUnknownScheme = errors.New("Unknown scheme")

	schemeToCmd = map[string][]string{
		"rsync": []string{"rsync", "-r"},
	}
)

type action func(*entry, string) ([]string, error)

type entry struct {
	Addr   string
	Scheme string
}

func Push() error {
	return execute(push)
}

func Pull() error {
	return execute(pull)
}

func push(e *entry, name string) ([]string, error) {
	res := schemeToCmd[e.Scheme]
	if res == nil {
		return nil, errUnknownScheme
	}
	res = append(res, withTrailingSlash(name), withTrailingSlash(e.Addr))
	return res, nil
}

func pull(e *entry, name string) ([]string, error) {
	res := schemeToCmd[e.Scheme]
	if res == nil {
		return nil, errUnknownScheme
	}
	res = append(res, withTrailingSlash(e.Addr), withTrailingSlash(name))
	return res, nil
}

func execute(pushOrPull action) error {
	f, err := os.Open(env.NotesMetadataPath(ConfigName))
	if err != nil {
		return err
	}
	remotes, err := parse(f)
	if err != nil {
		return err
	}
	// TODO: support multiple remotes
	args, err := pushOrPull(remotes[0], env.NotesRootPath())
	if err != nil {
		return err
	}
	return run(args)
}

func run(args []string) error {
	log.Println("running> ", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func withTrailingSlash(s string) string {
	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	return s
}
