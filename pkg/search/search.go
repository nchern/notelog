package search

import (
	"regexp"
	"strings"

	"github.com/nchern/notelog/pkg/note"
)

const (
	lastResultsFile = "last_search"
)

// Notes abstracts note collection to search in
type Notes interface {
	All() ([]*note.Note, error)
	MetadataFilename(name string) string
}

type matcherFunc func(string) []string

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
	return func(s string) []string {
		matches := terms.FindAllString(s, -1)
		if len(matches) == 0 {
			return matches
		}
		if excludeTerms != nil && len(excludeTerms.FindAllString(s, -1)) != 0 {
			// filter out excludeTerms if provided
			return []string{}
		}
		return matches
	}, nil
}

func (r *request) compileRegexp(terms []string) (*regexp.Regexp, error) {
	for i := 0; i < len(terms); i++ {
		terms[i] = regexp.QuoteMeta(terms[i])
	}
	return compileRx(regexOr(terms), !r.CaseSensitive)
}

func compileRx(expr string, ignoreCase bool) (*regexp.Regexp, error) {
	opts := ""
	if ignoreCase {
		opts = "(?i)"
	}
	return regexp.Compile(opts + expr)
}

// Engine represents a simple search engine over notes
type Engine struct {
	OnlyNames     bool // TODO: consolidate with Renderers
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
	req.CaseSensitive = s.CaseSensitive

	match, err := req.matcher()
	if err != nil {
		return nil, err
	}

	return searchInNotes(s.notes, match, s.OnlyNames)
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

type byNameAndLine []*Result

func (a byNameAndLine) Len() int      { return len(a) }
func (a byNameAndLine) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byNameAndLine) Less(i, j int) bool {
	if a[i].name == a[j].name {
		return a[i].lineNum < a[j].lineNum
	}
	return a[i].name < a[j].name
}
