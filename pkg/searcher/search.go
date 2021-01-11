package searcher

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultGrep     = "grep -E"
	defaultGrepArgs = "-rni"

	lastResultsFile = "last_search"
)

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

	notes   Notes
	grepCmd string

	out io.Writer
}

// NewSearcher returns a new Searcher instance
func NewSearcher(notes Notes, out io.Writer) *Searcher {
	return &Searcher{
		out:     out,
		notes:   notes,
		grepCmd: env.Get("NOTELOG_GREP", defaultGrep),
	}
}

// Search runs the search over all notes in notes home and prints results to stdout
// Terms grammar looks like: "foo bar -buzz -fuzz" where -xxx means exclude xxx matches from the output
func (s *Searcher) Search(terms ...string) error {
	req := parseRequest(terms...)

	cmd, err := buildSearchCmd(s.grepCmd, s.notes, req)
	if err != nil {
		return err
	}

	cmd.Stdout = s.out
	cmd.Stderr = os.Stderr

	if s.SaveResults {
		return s.runSearchAndSaveResults(cmd)
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (s *Searcher) runSearchAndSaveResults(cmd *exec.Cmd) error {
	buf := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(s.out, buf)

	err := cmd.Run()
	if err != nil {
		return err
	}

	f, err := os.Create(s.notes.MetadataFilename(lastResultsFile))
	if err != nil {
		return err
	}
	defer f.Close()

	return stripTermColors(buf, f)
}

func buildSearchCmd(grepCmd string, notes Notes, req *request) (*exec.Cmd, error) {
	cmdName, extraArgs, err := parseToCmdAndExtraArgs(grepCmd)
	if err != nil {
		return nil, err
	}

	if len(req.excludeTerms) > 0 {
		return searchCmdWithExcludeTerms(cmdName, extraArgs, req, notes.HomeDir()), nil
	}

	args := append(extraArgs, defaultGrepArgs, regexOr(req.terms), notes.HomeDir())
	return exec.Command(cmdName, args...), nil
}

func searchCmdWithExcludeTerms(cmd string, args []string, req *request, homeDir string) *exec.Cmd {
	findArgs := append(args, defaultGrepArgs, quote(regexOr(req.terms)), homeDir)
	findCmd := c(append([]string{cmd}, findArgs...)...)

	excludeCmd := c(cmd, strings.Join(args, " "), "-vi", quote(regexOr(req.excludeTerms)))

	return exec.Command("sh", "-c", pipe(findCmd, excludeCmd))
}

func parseToCmdAndExtraArgs(s string) (cmd string, args []string, err error) {
	toks, err := shlex.Split(strings.TrimSpace(s))
	if err != nil {
		return
	}

	if len(toks) > 0 {
		cmd = toks[0]
	}
	if len(toks) > 1 {
		args = toks[1:]
	}
	return
}

func pipe(s ...string) string {
	return strings.Join(s, " | ")
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
