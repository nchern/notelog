package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

var (
	// \x1b(or \x1B)	is the escape special character (sed does not support alternatives \e and \033)
	// \[				is the second character of the escape sequence
	// [0-9;]*			is the color value(s) regex
	// m				is the last character of the escape sequence
	termEscapeSequence = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

func stripColor(in io.Reader, out io.Writer) error {
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

func main() {
	must(stripColor(os.Stdin, os.Stdout))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
