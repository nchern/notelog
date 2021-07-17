package searcher

import (
	"bufio"
	"fmt"
	"log"

	"github.com/nchern/notelog/pkg/note"
)

type result struct {
	path    string
	lineNum int
	text    string
	name    string
}

func (r *result) String() string {
	return fmt.Sprintf("%s:%d:%s", r.path, r.lineNum, r.text)
}

func (r *result) Display() string {
	// TODO: elaborate better name
	return fmt.Sprintf("%s:%d:%s", r.name, r.lineNum, r.text)
}

func searchNote(nt *note.Note, match matcherFunc) ([]*result, error) {
	r, err := nt.Reader()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	lnum := 0
	res := []*result{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lnum++
		l := scanner.Text()
		if match(l) {
			res = append(res, &result{
				path:    nt.FullPath(),
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

func searchInNotes(notes []*note.Note, req *request, onlyNames bool) ([]*result, error) {
	match, err := req.matcher()
	if err != nil {
		return nil, err
	}

	resChan := make(chan []*result, len(notes))
	errChan := make(chan error, len(notes))

	for _, nt := range notes {
		go func(nt *note.Note) {
			results, err := searchNote(nt, match)
			if err != nil {
				errChan <- err
				return
			}
			if match(nt.Name()) {
				results = append(results, &result{
					path:    nt.FullPath(),
					lineNum: 1,
					text:    " ",
					name:    nt.Name(),
				})
			}
			resChan <- results
		}(nt)
	}

	names := map[string]bool{}
	results := []*result{}
	for i := 0; i < len(notes); i++ {
		select {
		case found := <-resChan:
			if onlyNames {
				for _, res := range found {
					if names[res.path] {
						continue
					}
					names[res.path] = true
					results = append(results, res)
				}
			} else {
				results = append(results, found...)
			}
		case err := <-errChan:
			log.Println(err) // TODO: find better way of error handling
		}
	}
	return results, nil
}
