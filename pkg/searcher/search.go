package searcher

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
	"github.com/nchern/notelog/pkg/env"
	"github.com/nchern/notelog/pkg/note"
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
	onlyNames bool

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

func (s *Searcher) OnlyNames() *Searcher {
	s.onlyNames = true
	return s
}

// Search runs the search over all notes in notes home and prints results to stdout
// Terms grammar looks like: "foo bar -buzz -fuzz" where -xxx means exclude xxx matches from the output
func (s *Searcher) Search(terms ...string) error {
	req := parseRequest(terms...)
	var resF io.Writer = ioutil.Discard

	if s.SaveResults {
		f, err := os.Create(s.notes.MetadataFilename(lastResultsFile))
		if err != nil {
			return err
		}
		defer f.Close()
		resF = f
	}

	results, err := s.searchInNotes(req)
	if err != nil {
		return err
	}

	return s.outputResults(results, resF)
}

func (s *Searcher) outputResults(results []string, persistentOut io.Writer) error {
	names := map[string]bool{}
	for _, res := range results {
		if s.onlyNames {
			toks := strings.Split(res, ":")
			if len(toks) < 1 {
				continue
			}
			fileName := toks[0]
			if names[fileName] {
				continue
			}
			names[fileName] = true
			res = fileName + ":1:"
		}

		if _, err := fmt.Fprintln(s.out, res); err != nil {
			return err
		}
		uncolored := termEscapeSequence.ReplaceAllString(res, "")
		if _, err := fmt.Fprintln(persistentOut, uncolored); err != nil {
			return err
		}
	}
	return nil
}

func (s *Searcher) searchInNotes(req *request) ([]string, error) {
	buf := &bytes.Buffer{}
	cmdArgs, err := buildSearchCmdAndArgs(s.grepCmd, s.notes, req)
	if err != nil {
		return nil, err
	}

	// Exclude notelog's dot dir from results
	cmdArgs = append(cmdArgs, fmt.Sprintf("grep -v '/%s/'", note.DotNotelogDir))

	cmd := exec.Command("sh", "-c", pipe(cmdArgs...))
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	results := []string{}

	scn := bufio.NewScanner(buf)
	for scn.Scan() {
		results = append(results, scn.Text())
	}
	if err := scn.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func buildSearchCmdAndArgs(grepCmd string, notes Notes, req *request) ([]string, error) {
	cmdName, extraArgs, err := parseToCmdAndExtraArgs(grepCmd)
	if err != nil {
		return nil, err
	}

	if len(req.excludeTerms) > 0 {
		return searchCmdWithExcludeTerms(cmdName, extraArgs, req, notes.HomeDir()), nil
	}

	args := append(extraArgs, defaultGrepArgs, quote(regexOr(req.terms)), notes.HomeDir())
	findCmd := c(append([]string{cmdName}, args...)...)

	return []string{findCmd}, nil
}

func searchCmdWithExcludeTerms(cmd string, args []string, req *request, homeDir string) []string {
	findArgs := append(args, defaultGrepArgs, quote(regexOr(req.terms)), homeDir)
	findCmd := c(append([]string{cmd}, findArgs...)...)

	excludeCmd := c(cmd, strings.Join(args, " "), "-vi", quote(regexOr(req.excludeTerms)))

	return []string{findCmd, excludeCmd}
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
