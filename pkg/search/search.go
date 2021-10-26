package search

import (
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
	OnlyNames     bool
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

// Engine represents a search engine over notes
type Engine struct {
	OnlyNames bool

	CaseSensitive bool

	notes Notes
}

// NewEngine returns a new Engine instance
func NewEngine(notes Notes) *Engine {
	return &Engine{
		notes: notes,
	}
}

// Search runs the search over all notes in notes home and prints results to stdout
// Terms grammar looks like: "foo bar -buzz -fuzz" where -xxx means exclude xxx matches from the output
func (s *Engine) Search(terms ...string) ([]*Result, error) {
	req := parseRequest(terms...)
	req.OnlyNames = s.OnlyNames
	req.CaseSensitive = s.CaseSensitive

	l, err := s.notes.All()
	if err != nil {
		return nil, err
	}

	res, err := searchInNotes(l, req)
	if err != nil {
		return nil, err
	}

	sort.Sort(byName(res))

	return res, nil
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

type byName []*Result

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].name < a[j].name }
