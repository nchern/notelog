package cli

import (
	"errors"
	"regexp"
	"strings"
)

const (
	validateNoteNameRegex = "^[a-zA-Z0-9-_]+?$"
)

var (
	validName = regexp.MustCompile(validateNoteNameRegex)

	errNameStartsWithDot = errors.New("Note name can not start with '.'")
	errEmptyName         = errors.New("Empty note name. Specify the real name")
	errNameRegexNoMatch  = errors.New("Note name must comply the following regex: " + validateNoteNameRegex)
)

func validateNoteName(name string) error {
	if name == "" {
		return errEmptyName
	}
	if strings.HasPrefix(name, ".") {
		return errNameStartsWithDot
	}
	if !validName.MatchString(name) {
		return errNameRegexNoMatch
	}
	return nil
}
