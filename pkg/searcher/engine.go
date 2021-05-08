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
}

func (r *result) String() string {
	return fmt.Sprintf("%s:%d:%s", r.path, r.lineNum, r.text)
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
			res = append(res, &result{path: nt.FullPath(), lineNum: lnum, text: l})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func searchInNotes(notes []*note.Note, req *request) ([]*result, error) {
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
			resChan <- results
		}(nt)
	}

	res := []*result{}
	for i := 0; i < len(notes); i++ {
		select {
		case found := <-resChan:
			res = append(res, found...)
		case err := <-errChan:
			log.Println(err) // TODO: find better way of error handling
		}
	}
	return res, nil
}

func searchInNames(notes []*note.Note, req *request) ([]*result, error) {
	match, err := req.matcher()
	if err != nil {
		return nil, err
	}

	res := []*result{}
	for _, it := range notes {
		if match(it.Name()) {
			res = append(res, &result{path: it.FullPath(), lineNum: 1, text: " "})
		}
	}

	return res, nil
}
