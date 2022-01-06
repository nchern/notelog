package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	indent     = regexp.MustCompile(`^\s+?`)
	listNumber = regexp.MustCompile(`(^\s*?\d+?\.\s*|^\s*?\d+?\s+?|^\s*?(\*|#)+\s*)`)
)

// Format formats identation list
func Format(src io.Reader, dst io.Writer) error {
	i := 0
	itemIndent := ""
	isListItem := false

	p := &printer{w: dst}
	scanner := bufio.NewScanner(src)

	for scanner.Scan() {
		if i > 0 {
			p.Println()
		}

		l := scanner.Text()
		l = strings.TrimSpace(l)

		if listNumber.MatchString(l) {
			isListItem = true
			itemIndent = strings.Repeat(" ", len(listNumber.FindString(l)))
		} else {
			if isListItem {
				if l != "" {
					l = itemIndent + indent.ReplaceAllString(l, "")
				}
			}
		}

		p.Print(l)

		if p.err != nil {
			return p.err
		}
		i++
	}
	return scanner.Err()
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
