package searcher

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultGrep     = "egrep"
	defaultGrepArgs = "-rni"

	lastResultsFile = "last_search"
)

var grepCmd = env.Get("NOTELOG_GREP", defaultGrep)

// Notes abstracts note collection to search in
type Notes interface {
	HomeDir() string
	MetadataFilename(name string) string
}

type request struct {
	terms        []string
	excludeTerms []string
}

// Searcher represents a search engine over notes
type Searcher struct {
	SaveResults bool
	notes       Notes
	out         io.Writer
}

// NewSearcher returns a new Searcher instance
func NewSearcher(notes Notes, out io.Writer) *Searcher {
	return &Searcher{
		notes: notes,
		out:   out,
	}
}

// Search runs the search over all notes in notes home and prints results to stdout
// Terms grammar looks like: "foo bar -buzz -fuzz" where -xxx means exclude xxx matches from the output
func (s *Searcher) Search(terms ...string) error {
	var cmd *exec.Cmd
	req := parseRequest(terms...)

	if len(req.excludeTerms) > 0 {
		findCmd := c(grepCmd, defaultGrepArgs, quote(regexOr(req.terms)), s.notes.HomeDir())
		excludeCmd := c(grepCmd, "-vi", quote(regexOr(req.excludeTerms)))

		cmd = exec.Command("sh", "-c", fmt.Sprintf("%s | %s", findCmd, excludeCmd))
	} else {
		cmd = exec.Command(grepCmd, defaultGrepArgs, regexOr(req.terms), s.notes.HomeDir())
	}

	cmd.Stderr = os.Stderr

	cmd.Stdout = s.out
	if s.SaveResults {
		f, err := os.Create(s.notes.MetadataFilename(lastResultsFile))
		if err != nil {
			return err
		}
		defer f.Close()
		cmd.Stdout = io.MultiWriter(s.out, f)
	}

	return cmd.Run()
}

func c(s ...string) string {
	return strings.Join(s, " ")
}

func quote(s string) string {
	return "'" + s + "'"
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
