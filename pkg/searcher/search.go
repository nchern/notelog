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
// Terms grammar looks like: "foo bar -buzz -fuzz" where -xxx means exclude xxx matches from the output
func Search(notes Notes, out io.Writer, terms ...string) error {
	var cmd *exec.Cmd
	req := parseRequest(terms...)

	if len(req.excludeTerms) > 0 {
		findCmd := strings.Join([]string{grepCmd, defaultGrepArgs, quote(regexOr(req.terms)), notes.HomeDir()}, " ")
		excludeCmd := strings.Join([]string{grepCmd, "-vi", quote(regexOr(req.excludeTerms))}, " ")
		cmd = exec.Command("sh", "-c", fmt.Sprintf("%s | %s", findCmd, excludeCmd))
	} else {
		cmd = exec.Command(grepCmd, defaultGrepArgs, regexOr(req.terms), notes.HomeDir())
	}

	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	return cmd.Run()
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
