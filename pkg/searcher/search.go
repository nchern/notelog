package searcher

import (
	"os"
	"os/exec"
	"strings"

	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultGrep     = "egrep"
	defaultGrepArgs = "-rni"
)

var grepCmd = env.Get("NOTELOG_GREP", defaultGrep)

// Notes abstracts note collection to search in
type Notes interface {
	HomeDir() string
}

type request struct {
	terms        []string
	excludeTerms []string
}

// Search runs the search over all notes in notes home and prints results to stdout
// Terms grammar looks like: "foo bar -buzz -fuzz" where -buzz means exclude buzz matches from the output
func Search(notes Notes, terms ...string) error {
	req := requestFromStrings(terms...)
	s := regexOr(req.terms)
	cmd := exec.Command(grepCmd, defaultGrepArgs, s, notes.HomeDir())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func regexOr(terms []string) string {
	return "(" + strings.Join(terms, "|") + ")"
}

func requestFromStrings(args ...string) *request {
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
