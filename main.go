package main

import (
	"flag"
	"log"

	"github.com/nchern/notelog/pkg/cli"
	"github.com/nchern/notelog/pkg/note"
)

func init() {
	log.SetFlags(0)
	flag.Parse()

	must(note.NewList().Init())
}

func main() {
	must(cli.Execute())
}

func must(err error) {
	dieIf(err)
}

func dieIf(err error) {
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}
}
