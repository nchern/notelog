package search

import (
	"bufio"
	"fmt"
	"log"

	"github.com/nchern/notelog/pkg/note"
)

// Result represents one search result
type Result struct {
	lineNum int
	text    string
	name    string
}

// String returns stringified representation of Result
func (r *Result) String() string {
	return fmt.Sprintf("%s:%d", r.name, r.lineNum)
}

// Display returns display friendly name of a result
func (r *Result) Display() string {
	// TODO: elaborate better name
	return fmt.Sprintf("%s:%d:%s", r.name, r.lineNum, r.text)
}

func searchNote(nt *note.Note, match matcherFunc) ([]*Result, error) {
	r, err := nt.Reader()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	lnum := 0
	res := []*Result{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lnum++
		l := scanner.Text()
		if match(l) {
			res = append(res, &Result{
				lineNum: lnum,
				text:    l,
				name:    nt.Name(),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func searchInNotes(notes []*note.Note, req *request, onlyNames bool) ([]*Result, error) {
	match, err := req.matcher()
	if err != nil {
		return nil, err
	}

	resChan := make(chan []*Result, len(notes))
	errChan := make(chan error, len(notes))

	for _, nt := range notes {
		go func(nt *note.Note) {
			results, err := searchNote(nt, match)
			if err != nil {
				errChan <- err
				return
			}
			if match(nt.Name()) {
				results = append(results, &Result{
					lineNum: 1,
					text:    " ",
					name:    nt.Name(),
				})
			}
			resChan <- results
		}(nt)
	}

	names := map[string]bool{}
	results := []*Result{}
	for i := 0; i < len(notes); i++ {
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
	return results, nil
}
