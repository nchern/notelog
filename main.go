package main

import (
	"flag"
	"log"

	"github.com/nchern/notelog/pkg/cli"
)

func init() {
	log.SetFlags(0)
	flag.Parse()
}

func main() {
	must(cli.Execute(*cli.Command))
}

func must(err error) {
	dieIf(err)
}

func dieIf(err error) {
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}
}
