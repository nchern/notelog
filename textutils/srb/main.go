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

	cmdEdit  = "edit"
	cmdPrint = "print"
)

var (
	// \x1b(or \x1B)	is the escape special character (sed does not support alternatives \e and \033)
	// \[				is the second character of the escape sequence
	// [0-9;]*			is the color value(s) regex
	// m				is the last character of the escape sequence
	termEscapeSequence = regexp.MustCompile(`\x1b\[[0-9;]*m`)

	cmd = flag.String("c", cmdEdit, "command to perform")
)

// srb - Search Results (Console) Browser
func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) < 1 {
		must(save())
		return
	}

	n, err := strconv.Atoi(flag.Args()[0])
	dieIf(err)

	switch *cmd {
	case cmdEdit:
		must(edit(n))
	case cmdPrint:
		must(printLine(n))
	default:
		fmt.Fprintf(os.Stderr, "wrong command: %s\n\n", *cmd)
		usage()
		os.Exit(2)
	}
}

func printLine(n int) error {
	f, err := os.Open(dotFileFullPath())
	if err != nil {
		return err
	}
	defer f.Close()

	e, err := getNthFilename(n, f)
	if err != nil {
		return err
	}

	_, err = fmt.Println(e.Row)
	return err
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [number]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "	when no arguments provided, reads stdin line by line, outputs numbered lines and stores it in cache\n")
	fmt.Fprintf(os.Stderr, "	when [number] is given, fetches the line with this number, treats it as '<filename>:<lineno>: ...'\n")
	fmt.Fprintf(os.Stderr, "		and tries to open <filename> in your $EDITOR at line #lineno\n")
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

func edit(n int) error {
	f, err := os.Open(dotFileFullPath())
	if err != nil {
		return err
	}
	defer f.Close()

	e, err := getNthFilename(n, f)
	if err != nil {
		return err
	}

	return openEditor(getEditor(), e.fullPath(), e.RowNum)
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
	e.Row = s

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
