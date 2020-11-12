package searcher

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
	"github.com/nchern/notelog/pkg/env"
)

const (
	defaultGrep     = "egrep"
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

	cmd, err := s.makeCmd(s.notes, req)
	if err != nil {
		return err
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

// GetLastNthResult returns nth result from last saved search results
func GetLastNthResult(notes Notes, n int) (string, error) {
	f, err := os.Open(notes.MetadataFilename(lastResultsFile))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // just an empty result same if we asked for non-existing item
		}
		return "", err
	}
	i := 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if i == n {
			return scanner.Text(), nil
		}
		i++
	}
	return "", scanner.Err()
}

func (s *Searcher) makeCmd(notes Notes, req *request) (*exec.Cmd, error) {
	cmd, extraArgs, err := parseToCmdAndExtraArgs(s.grepCmd)
	if err != nil {
		return nil, err
	}

	args := []string{defaultGrepArgs}
	if len(extraArgs) > 0 {
		args = append(extraArgs, args...)
	}

	if len(req.excludeTerms) > 0 {
		return searchCmdWithExcludeTerms(cmd, args, req, notes.HomeDir()), nil
	}

	args = append(args, regexOr(req.terms), notes.HomeDir())
	return exec.Command(cmd, args...), nil
}

func searchCmdWithExcludeTerms(cmd string, args []string, req *request, homeDir string) *exec.Cmd {
	args = append(args, quote(regexOr(req.terms)), homeDir)
	findCmd := c(append([]string{cmd}, args...)...)

	excludeCmd := c(cmd, "-vi", quote(regexOr(req.excludeTerms)))

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
