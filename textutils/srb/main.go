package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	dotFile       = ".srb"
	defaultEditor = "vi"
)

var (
	nth = flag.Int("n", 1, "Results ordinal number")

	// \x1b(or \x1B)	is the escape special character (sed does not support alternatives \e and \033)
	// \[				is the second character of the escape sequence
	// [0-9;]*			is the color value(s) regex
	// m				is the last character of the escape sequence
	termEscapeSequence = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

func dotFileFullPath() string {
	home, err := os.UserHomeDir()
	dieIf(err)
	return filepath.Join(home, dotFile)
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor != "" {
		return editor
	}
	return defaultEditor
}

func save() error {
	f, err := os.Create(dotFileFullPath())
	if err != nil {
		return err
	}
	defer f.Close()

	return readAndSave(os.Stdin, os.Stdout, f)
}

func edit() error {
	f, err := os.Open(dotFileFullPath())
	if err != nil {
		return err
	}
	defer f.Close()

	e, err := getNthFilename(*nth, f)
	if err != nil {
		return err
	}

	return openEditor(getEditor(), e.fullPath(), e.RowNum)
}

// srb - Search Results (Console) Browser
// option: srp stands for Search results page
func main() {
	if len(os.Args) < 2 {
		must(save())
		return
	}

	flag.Parse()
	must(edit())
}

func readAndSave(src io.Reader, dst io.Writer, file io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(file, "# %s\n\n", cwd); err != nil {
		return err
	}

	i := 0
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		i++
		l := scanner.Text()
		if _, err := fmt.Fprintf(dst, "%6d %s\n", i, l); err != nil {
			return err
		}
		l = termEscapeSequence.ReplaceAllString(l, "")
		if _, err := fmt.Fprintln(file, l); err != nil {
			return err
		}
	}
	return scanner.Err()
}

type entry struct {
	RowNum int64

	Row      string
	Root     string
	Filename string
}

func (e *entry) fullPath() string {
	if filepath.IsAbs(e.Filename) {
		return e.Filename
	}
	return filepath.Join(e.Root, e.Filename)
}

func (e *entry) parse(s string) {
	// example: ../pkg/checklist/sort.go:24:type kind int

	toks := strings.Split(s, ":")
	e.Filename = toks[0]
	if len(toks) > 1 {
		e.RowNum, _ = strconv.ParseInt(toks[1], 0, 64)
	}
}

func getNthFilename(n int, file io.Reader) (entry, error) {
	i := 0
	var res entry

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
			res.Root = strings.TrimSpace(strings.TrimPrefix(l, "#"))
			continue
		} else if strings.TrimSpace(l) == "" {
			continue
		}
		i++
		res.parse(l)
		if i == n {
			return res, nil
		}
	}
	return res, scanner.Err()
}

func openEditor(editor string, filename string, lnum int64) error {
	args := []string{filename}
	if lnum > 0 {
		args = append(args, fmt.Sprintf("+%d", lnum))
	}
	cmd := exec.Command(editor, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func dieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func must(err error) {
	dieIf(err)
}
