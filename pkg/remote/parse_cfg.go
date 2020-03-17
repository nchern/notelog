package remote

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	// ErrConfigEmpty is returned when no remotes were read from the config
	ErrConfigEmpty = errors.New("remote: no remotes configured")

	errConfigMalformed = errors.New("remote: bad line")
)

func parse(r io.Reader) ([]*entry, error) {
	res := []*entry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" ||
			strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) < 2 {
			return nil, fmt.Errorf("%w: '%s'", errConfigMalformed, line)
		}

		addr := strings.TrimPrefix(tokens[1], "//")

		res = append(res, &entry{Scheme: tokens[0], Addr: addr})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, ErrConfigEmpty
	}
	return res, nil
}
