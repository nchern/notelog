package searcher

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

// Taken from here: https://github.com/acarl005/stripansi/blob/master/stripansi.go
const ansiTermColors = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var (
	termEscapeSequence = regexp.MustCompile(ansiTermColors)
)

func stripTermColors(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		l := scanner.Text()

		l = termEscapeSequence.ReplaceAllString(l, "")
		if _, err := fmt.Fprintln(out, l); err != nil {
			return err
		}
	}

	return scanner.Err()
}
