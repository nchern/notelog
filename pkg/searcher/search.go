package searcher

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/nchern/notelog/pkg/note"
)

const (
	lastResultsFile = "last_search"
)

var ErrNoResults = errors.New("search: nothing found")

// Notes abstracts note collection to search in
type Notes interface {
	HomeDir() string
	All() ([]*note.Note, error)
	MetadataFilename(name string) string
}

type request struct {
	terms        []string
	excludeTerms []string
}

func (r *request) termsToRegexp() (terms *regexp.Regexp, excludeTerms *regexp.Regexp, err error) {
	terms, err = regexp.Compile("(?i)" + regexOr(r.terms))
	if err != nil {
		return
	}

	if len(r.excludeTerms) < 1 {
		return
	}

	excludeTerms, err = regexp.Compile("(?i)" + regexOr(r.excludeTerms))
	if err != nil {
		return
	}
	return
}

// Searcher represents a search engine over notes
type Searcher struct {
	OnlyNames bool

	SaveResults bool

	notes Notes

	out io.Writer
}

// NewSearcher returns a new Searcher instance
func NewSearcher(notes Notes, out io.Writer) *Searcher {
	return &Searcher{
		out:   out,
		notes: notes,
	}
}

// Search runs the search over all notes in notes home and prints results to stdout
// Terms grammar looks like: "foo bar -buzz -fuzz" where -xxx means exclude xxx matches from the output
func (s *Searcher) Search(terms ...string) error {
	req := parseRequest(terms...)
	var resF = ioutil.Discard

	if s.SaveResults {
		f, err := os.Create(s.notes.MetadataFilename(lastResultsFile))
		if err != nil {
			return err
		}
		defer f.Close()
		resF = f
	}

	matchedNames := []*result{}
	matchedNamesErr := make(chan error)
	go func() {
		var err error
		matchedNames, err = searchInNames(s.notes, req)
		matchedNamesErr <- err
	}()

	res, err := searchInNotes(s.notes, req)
	if err != nil {
		return err
	}

	if err := <-matchedNamesErr; err != nil {
		return err
	}

	res = append(res, matchedNames...)
	if len(res) == 0 {
		return ErrNoResults
	}

	return s.outputResults(res, resF)
}

func (s *Searcher) outputResults(results []*result, persistentOut io.Writer) error {
	names := map[string]bool{}
	if s.OnlyNames {
		sort.Sort(byPath(results))
	}
	for _, res := range results {
		orig := *res
		if s.OnlyNames {
			if names[res.path] {
				continue
			}
			names[res.path] = true
			res.text = " "
			res.lineNum = 1
		}

		if _, err := fmt.Fprintln(s.out, res); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(persistentOut, &orig); err != nil {
			return err
		}
	}
	return nil
}

func regexOr(terms []string) string {
	return "(" + strings.Join(terms, "|") + ")"
}

func parseRequest(args ...string) *request {
	res := &request{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			res.excludeTerms = append(res.excludeTerms, strings.TrimPrefix(arg, "-"))
			continue
		}
		res.terms = append(res.terms, arg)
	}
	return res
}

type byPath []*result

func (a byPath) Len() int           { return len(a) }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPath) Less(i, j int) bool { return a[i].path < a[j].path }
