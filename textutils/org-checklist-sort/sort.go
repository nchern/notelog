package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

const (
	other kind = iota
	done
	blank
	notDone
)

var (
	whitespaceRx = regexp.MustCompile(`^\s*$`)

	todoItemDoneRx    = regexp.MustCompile(`(?i)^\s*?-\s*?\[x\]\s+?.+$`)
	todoItemNotDoneRx = regexp.MustCompile(`(?i)^\s*?-\s*?\[\s?\]\s+?.+$`)
)

type kind int

// Sort sorts checkbox lists from r and outputs sorted to w
func Sort(r io.Reader, w io.Writer) error {
	emptyLinesCount := 0
	prevIndent, curIndent := 0, 0

	doneBuffer := []string{}
	prevDone := false

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		lineKind := getLineKind(line)
		curIndent = getIndentLen(line)

		if prevDone && curIndent > prevIndent {
			doneBuffer = append(doneBuffer, line)
			continue
		}
		prevDone = false
		prevIndent = curIndent

		switch lineKind {
		case done:
			writeEmptyLines(w, emptyLinesCount)
			emptyLinesCount = 0

			doneBuffer = append(doneBuffer, line)

			if !prevDone {
				prevIndent = curIndent
				prevDone = true
			}
		case notDone:
			writeEmptyLines(w, emptyLinesCount)
			emptyLinesCount = 0

			fmt.Fprintln(w, line)
		case blank:
			emptyLinesCount++
			if emptyLinesCount > 1 {
				writeLines(w, doneBuffer)
				writeEmptyLines(w, emptyLinesCount)

				emptyLinesCount = 0
				doneBuffer = []string{}
			}
		default:
			writeLines(w, doneBuffer)
			writeEmptyLines(w, emptyLinesCount)

			emptyLinesCount = 0
			doneBuffer = []string{}

			fmt.Fprintln(w, line)
		}
	}
	writeLines(w, doneBuffer)

	return scanner.Err()
}

func getLineKind(line string) kind {
	if todoItemDoneRx.MatchString(line) {
		return done
	}
	if todoItemNotDoneRx.MatchString(line) {
		return notDone
	}
	if whitespaceRx.MatchString(line) {
		return blank
	}
	return other
}

func getIndentLen(s string) int {
	n := 0
	for _, c := range s {
		if c == ' ' || c == '\t' {
			n++
			continue
		}
		break
	}

	return n
}

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
