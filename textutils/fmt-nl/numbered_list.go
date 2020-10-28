package main // numberedlist

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	indent     = regexp.MustCompile(`^\s+?`)
	listNumber = regexp.MustCompile(`(^\s*?\d+?\.|^\s*?\d+?\s+?)`)
)

// Format formats numbered list provided on a given reader
func Format(src io.Reader, dst io.Writer) error {
	items, err := readItems(src)
	if err != nil {
		return err
	}

	return printOutNumbered(dst, items)
}

func printOutNumbered(w io.Writer, items []string) error {
	n := 0
	p := &printer{w: w}

	for i, l := range items {
		if i > 0 {
			p.Println()
		}

		l = cleanOldListNumbers(l)
		if isEmptyOrSubitem(l) {
			p.Print(l)
			continue
		}

		n++
		p.Printf("%d. %s", n, l)
		if p.err != nil {
			return p.err
		}
	}
	return p.err
}

func cleanOldListNumbers(s string) string {
	if listNumber.MatchString(s) {
		return strings.TrimSpace(listNumber.ReplaceAllString(s, ""))
	}
	return s
}

func isEmptyOrSubitem(s string) bool {
	return strings.TrimSpace(s) == "" || indent.MatchString(s)
}

func readItems(r io.Reader) ([]string, error) {
	items := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		items = append(items, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type printer struct {
	err error
	w   io.Writer
}

func (p *printer) Printf(format string, a ...interface{}) {
	if p.err != nil {
		return
	}
	_, p.err = fmt.Fprintf(p.w, format, a...)
}

func (p *printer) Println() { p.Printf("\n") }

func (p *printer) Print(s string) { p.Printf("%s", s) }
