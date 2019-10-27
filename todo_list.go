package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var (
	todoItemUndoneRx = regexp.MustCompile(`^\s*?-\s*?\[\s?\]\s+?.+$`)
	todoItemDoneRx   = regexp.MustCompile(`^\s*?-\s*?\[x\]\s+?.+$`)

	whitespaceRx = regexp.MustCompile(`^\s*$`)
)

func writeLines(w io.Writer, lines []string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeEmptyLines(w io.Writer, count int) error {
	for i := 0; i < count; i++ {
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

func sortTODOList(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	emptyLinesCount := 0
	doneBuffer := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		isDoneItem := todoItemDoneRx.MatchString(line)
		isUndoneItem := todoItemUndoneRx.MatchString(line)

		if isDoneItem {
			writeEmptyLines(w, emptyLinesCount)
			emptyLinesCount = 0
			doneBuffer = append(doneBuffer, line)
			continue
		}
		if isUndoneItem {
			writeEmptyLines(w, emptyLinesCount)
			emptyLinesCount = 0
			fmt.Fprintln(w, line)
			continue
		}
		if whitespaceRx.MatchString(line) || line == "" {
			emptyLinesCount++
			if emptyLinesCount > 1 {
				writeLines(w, doneBuffer)
				writeEmptyLines(w, emptyLinesCount)
				doneBuffer = []string{}
				emptyLinesCount = 0
			}
			continue
		}

		writeLines(w, doneBuffer)

		emptyLinesCount = 0
		doneBuffer = []string{}

		fmt.Fprintln(w, line)
	}
	writeLines(w, doneBuffer)

	return scanner.Err()
}
