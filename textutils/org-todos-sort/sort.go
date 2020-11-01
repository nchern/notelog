package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

const (
	normal = iota
	done
	undone
	nonTodoHeader
)

type state int

var (
	todoItemUndone = regexp.MustCompile(`\*\s+?TODO `)
	todoItemDone   = regexp.MustCompile(`\*\s+?DONE `)

	todoNot = regexp.MustCompile(`\*\s+?\w+`)
)

// Sort sorts org todo itmes from in and outputs sorted to out
// As of now it supports only top level todos
func Sort(in io.Reader, out io.Writer) error {
	i := 0
	p := &printer{w: out}
	scanner := bufio.NewScanner(in)

	doneBuf := []string{}
	var current state = normal

	for scanner.Scan() {
		if i > 0 && current != done {
			p.Println()
		}
		l := scanner.Text()
		newState := getState(l)

		isInTodo := current == done || current == undone
		if !(newState == normal && isInTodo) {
			current = newState
		}

		if current == done {
			doneBuf = append(doneBuf, l)
			continue
		}

		if current == nonTodoHeader {
			if err := flush(p, doneBuf); err != nil {
				return err
			}
			doneBuf = []string{}
		}

		p.Print(l)
		if p.err != nil {
			return p.err
		}
		i++
	}
	p.Println() // for external text editors integration
	if err := flush(p, doneBuf); err != nil {
		return err
	}

	return scanner.Err()
}

func flush(p *printer, lines []string) error {
	for _, l := range lines {
		p.Print(l)
		p.Println()
	}
	return p.err
}

func getState(s string) state {
	if todoItemUndone.MatchString(s) {
		return undone
	}
	if todoItemDone.MatchString(s) {
		return done
	}
	if todoNot.MatchString(s) {
		return nonTodoHeader
	}
	return normal
}

type printer struct {
	err error
	w   io.Writer
}

func (p *printer) Printf(format string, a ...interface{}) {
	if p.err != nil {
		return
	}
	_, p.err = fmt.Fprintf(p.w, format, a...)
}

func (p *printer) Println() { p.Printf("\n") }

func (p *printer) Print(s string) { p.Printf("%s", s) }
