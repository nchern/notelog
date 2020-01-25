package remote

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func parse(r io.Reader) ([]*entry, error) {
	res := []*entry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) < 2 {
			return nil, fmt.Errorf("remote: bad line: '%s'", line)
		}
		log.Println(line)
		addr := strings.TrimPrefix(tokens[1], "//")

		res = append(res, &entry{Scheme: tokens[0], Addr: addr})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
