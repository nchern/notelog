package main

import (
	"log"
	"os"
)

func main() {
	must(Format(os.Stdin, os.Stdout))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
