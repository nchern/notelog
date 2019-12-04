package todos

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var (
	todoItemUndoneRx = regexp.MustCompile(`(?i)^\s*?-\s*?\[\s?\]\s+?.+$`)
	todoItemDoneRx   = regexp.MustCompile(`(?i)^\s*?-\s*?\[x\]\s+?.+$`)

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

// Sort sorts todo lists from r and outputs sorted to w
func Sort(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	emptyLinesCount := 0
	doneBuffer := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		isDoneItem := todoItemDoneRx.MatchString(line)
		isUndoneItem := todoItemUndoneRx.MatchString(line)
		isBlankLine := whitespaceRx.MatchString(line) || line == ""

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
		if isBlankLine {
			emptyLinesCount++
			if emptyLinesCount > 1 {
				writeLines(w, doneBuffer)
				writeEmptyLines(w, emptyLinesCount)

				emptyLinesCount = 0
				doneBuffer = []string{}
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
