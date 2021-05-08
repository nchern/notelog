package searcher

import (
	"bufio"
	"fmt"
	"log"
	"regexp"

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

func searchNote(nt *note.Note, terms *regexp.Regexp, excludeTerms *regexp.Regexp) ([]*result, error) {
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
		if terms.MatchString(l) {
			if excludeTerms != nil && excludeTerms.MatchString(l) {
				// filter out excludeTerms if provided
				continue
			}
			res = append(res, &result{path: nt.FullPath(), lineNum: lnum, text: l})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func searchInNotes(notes Notes, req *request) ([]*result, error) {
	terms, excludeTerms, err := req.termsToRegexp()

	l, err := notes.All()
	if err != nil {
		return nil, err
	}

	resChan := make(chan []*result, len(l))
	errChan := make(chan error, len(l))

	for _, nt := range l {
		go func(nt *note.Note) {
			results, err := searchNote(nt, terms, excludeTerms)
			if err != nil {
				errChan <- err
				return
			}
			resChan <- results
		}(nt)
	}

	res := []*result{}
	for i := 0; i < len(l); i++ {
		select {
		case found := <-resChan:
			res = append(res, found...)
		case err := <-errChan:
			log.Println(err) // TODO: find better way of error handling
		}
	}
	return res, nil
}

func searchInNames(notes Notes, req *request) ([]*result, error) {
	terms, excludeTerms, err := req.termsToRegexp()
	if err != nil {
		return nil, err
	}

	res := []*result{}
	items, err := notes.All()
	if err != nil {
		return nil, err
	}

	for _, it := range items {
		if terms.MatchString(it.Name()) {
			if excludeTerms != nil && excludeTerms.MatchString(it.Name()) {
				// filter out excludeTerms if provided
				continue
			}
			res = append(res, &result{path: it.FullPath(), lineNum: 1, text: " "})
		}
	}

	return res, nil
}
