package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// Taken from here: https://github.com/acarl005/stripansi/blob/master/stripansi.go
const ansiTermColors = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var (
	termColorsRegex = regexp.MustCompile(ansiTermColors)
)

func stripColor(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		l := scanner.Text()

		l = termColorsRegex.ReplaceAllString(l, "")
		if _, err := fmt.Fprintln(out, l); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func main() {
	must(stripColor(os.Stdin, os.Stdout))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
