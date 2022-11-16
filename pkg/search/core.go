package search

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/nchern/notelog/pkg/note"
)

const (
	fgRed     = 31
	fgGreen   = 32
	fgMagenta = 35

	fgHiRed     = 91
	fgHiMagenta = 95

	bold      = "1"
	underline = "4"
)

// Result represents one search result
type Result struct {
	lineNum  int
	text     string
	name     string
	archived bool

	matches []string
}

// String returns stringified representation of Result
func (r *Result) String() string {
	archive := ""
	if r.archived {
		archive = "a"
	}
	return fmt.Sprintf("%s:%d:%s", r.name, r.lineNum, archive)
}

// Display returns display friendly name of a result
func (r *Result) Display(colored bool) string {
	// TODO: elaborate better name
	text := r.text
	name := r.name
	lineNum := fmt.Sprintf("%d", r.lineNum)
	if colored {
		for _, m := range r.matches {
			s := colorize(m, fgRed, bold)
			text = strings.Replace(text, m, s, -1)

			lineNum = colorize(lineNum, fgGreen, bold)

			if strings.Index(r.name, m) > -1 {
				toks := strings.Split(name, m)
				for i, tok := range toks {
					toks[i] = colorize(tok, fgMagenta, bold)
				}
				name = colorize(m, fgRed, bold, underline)
				name = strings.Join(toks, name)
			} else {
				name = colorize(name, fgMagenta, bold)
			}
		}
	}
	return fmt.Sprintf("%s:%s:%s", name, lineNum, text)
}

func colorize(s string, color int, attrs ...string) string {
	return fmt.Sprintf("\033[%d;%sm%s\033[0m", color, strings.Join(attrs, ";"), s)
}

func searchNote(notes Notes, nt *note.Note, matcher matcherFunc) ([]*Result, error) {
	var buf bytes.Buffer
	err := nt.Dump(&buf)
	if err != nil {
		return nil, err
	}

	lnum := 0
	res := []*Result{}
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		lnum++
		l := scanner.Text()
		matches := matcher(l)
		found := len(matches) != 0
		if found {
			res = append(res, &Result{
				lineNum:  lnum,
				text:     l,
				name:     nt.Name(),
				matches:  matches,
				archived: nt.Archived(),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func searchInNotes(notes Notes, matcher matcherFunc, onlyNames bool) ([]*Result, error) {

	l, err := notes.All()
	if err != nil {
		return nil, err
	}

	resChan := make(chan []*Result, len(l))
	errChan := make(chan error, len(l))

	for _, nt := range l {
		go func(nt *note.Note) {
			results, err := searchNote(notes, nt, matcher)
			if err != nil {
				errChan <- err
				return
			}
			matches := matcher(nt.Name())
			if len(matches) != 0 {
				results = append(results, &Result{
					lineNum:  1,
					text:     " ",
					name:     nt.Name(),
					matches:  matches,
					archived: nt.Archived(),
				})
			}
			resChan <- results
		}(nt)
	}

	names := map[string]bool{}
	results := []*Result{}
	for i := 0; i < len(l); i++ {
		select {
		case found := <-resChan:
			if onlyNames {
				for _, res := range found {
					if names[res.name] {
						continue
					}
					names[res.name] = true
					results = append(results, res)
				}
			} else {
				results = append(results, found...)
			}
		case err := <-errChan:
			log.Printf("ERROR search: %s", err) // TODO: find better way of error handling
		}
	}

	sort.Sort(byName(results))

	return results, nil
}
