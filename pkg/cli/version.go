package cli

import "fmt"

var version string

func printVersion() error {
	_, err := fmt.Printf("notelog version %s\n", version)
	return err
}
