package searcher

import (
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

// Notes abstracts note collection to search in
type Notes interface {
	HomeDir() string
	All() ([]*note.Note, error)
	MetadataFilename(name string) string
}

type matcherFunc func(string) bool

type request struct {
	CaseSensitive bool

	terms        []string
	excludeTerms []string
}

func (r *request) termsToRegexp() (terms *regexp.Regexp, excludeTerms *regexp.Regexp, err error) {
	terms, err = r.compileRegexp(r.terms)
	if err != nil {
		return
	}

	if len(r.excludeTerms) < 1 {
		return
	}

	excludeTerms, err = r.compileRegexp(r.excludeTerms)
	if err != nil {
		return
	}
	return
}

func (r *request) matcher() (matcherFunc, error) {
	terms, excludeTerms, err := r.termsToRegexp()
	if err != nil {
		return nil, err
	}
	return func(s string) bool {
		if !terms.MatchString(s) {
			return false
		}
		if excludeTerms != nil && excludeTerms.MatchString(s) {
			// filter out excludeTerms if provided
			return false
		}
		return true
	}, nil
}

func (r *request) compileRegexp(terms []string) (*regexp.Regexp, error) {
	opts := ""
	if !r.CaseSensitive {
		opts = "(?i)"
	}
	return regexp.Compile(opts + regexOr(terms))
}

// Searcher represents a search engine over notes
type Searcher struct {
	OnlyNames bool

	SaveResults bool

	CaseSensitive bool

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
func (s *Searcher) Search(terms ...string) (int, error) {
	req := parseRequest(terms...)
	req.CaseSensitive = s.CaseSensitive

	var resF = ioutil.Discard

	if s.SaveResults {
		f, err := os.Create(s.notes.MetadataFilename(lastResultsFile))
		if err != nil {
			return 0, err
		}
		defer f.Close()
		resF = f
	}

	l, err := s.notes.All()
	if err != nil {
		return 0, err
	}

	res, err := searchInNotes(l, req, s.OnlyNames)
	if err != nil {
		return 0, err
	}

	return len(res), s.outputResults(res, resF)
}

func (s *Searcher) outputResults(results []*result, persistentOut io.Writer) error {
	if s.OnlyNames {
		sort.Sort(byPath(results))
	}
	for _, res := range results {
		orig := *res
		if s.OnlyNames {
			res.text = " "
			res.lineNum = 1
		}

		if _, err := fmt.Fprintln(s.out, res.Display()); err != nil {
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
