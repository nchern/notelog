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
	// \x1b(or \x1B)	is the escape special character (sed does not support alternatives \e and \033)
	// \[				is the second character of the escape sequence
	// [0-9;]*			is the color value(s) regex
	// m				is the last character of the escape sequence
	termEscapeSequence = regexp.MustCompile(`\x1b\[[0-9;]*m`)

	flagPrint = flag.Bool("p", false, "prints a given row. Row must be given as a 1st positional argument")

	delimeter = ":"
)

// srb - Search Results (Console) Browser
func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) < 1 {
		must(save())
		return
	}

	row, err := parseArg(0)
	dieIf(err)

	if *flagPrint {
		col, _ := parseArg(1)
		must(printRowOrField(row, col))
	} else {
		must(edit(row))
	}
}

func parseArg(i int) (int, error) {
	x, err := strconv.Atoi(flag.Arg(i))
	if err != nil {
		return -1, err
	}
	return x, nil
}

func printRowOrField(row int, col int) error {
	e, err := getNthEntry(row)
	if err != nil {
		return err
	}

	if col > -1 {
		_, err = fmt.Println(e.getField(col))
	} else {
		_, err = fmt.Println(e.Row)
	}
	return err
}

func edit(n int) error {
	e, err := getNthEntry(n)
	if err != nil {
		return err
	}

	return openEditor(getEditor(), e.fullPath(), e.RowNum)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [number]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "	when no arguments provided, reads stdin line by line, outputs numbered lines and stores them in cache\n")
	fmt.Fprintf(os.Stderr, "	when [number] is given, fetches the line with this ordinal number, treats it as '<filename>:<lineno>: ...'\n")
	fmt.Fprintf(os.Stderr, "		and tries to open <filename> in your $EDITOR at the line #lineno\n")
	fmt.Fprintf(os.Stderr, "\nAdditional params:\n")

	flag.PrintDefaults()
}

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

	tokenized []string
}

func (e *entry) getField(col int) string {
	if col < 0 || col >= len(e.tokenized) {
		return ""
	}
	return e.tokenized[col]
}

func (e *entry) fullPath() string {
	if filepath.IsAbs(e.Filename) {
		return e.Filename
	}
	return filepath.Join(e.Root, e.Filename)
}

func (e *entry) parse(s string) {
	// example: ../pkg/checklist/sort.go:24:type kind int
	e.Row = s

	e.tokenized = strings.Split(s, delimeter)
	if len(e.tokenized) < 1 {
		return
	}

	e.Filename = e.tokenized[0]
	if len(e.tokenized) > 1 {
		e.RowNum, _ = strconv.ParseInt(e.tokenized[1], 0, 64)
	}
}

func getNthEntry(n int) (entry, error) {
	i := 0
	var res entry

	f, err := os.Open(dotFileFullPath())
	if err != nil {
		return res, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
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
