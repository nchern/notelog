package main

import (
	"flag"
	"log"

	"github.com/nchern/notelog/pkg/cli"
)

func must(err error) {
	dieIf(err)
}

func dieIf(err error) {
	if err != nil {
		log.Fatalf("FATAL: %s", err)
	}
}

func init() {
	log.SetFlags(0)
	flag.Parse()
}

func main() {
	must(cli.Execute(*cli.Command))
}
