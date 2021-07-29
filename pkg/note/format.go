package note

import "fmt"

// Format is the type of the Note
type Format string

const (
	// Unknown represents unknown note types
	Unknown Format = ""
	// Md represents markdown note types
	Md Format = "md"
	// Org represents org note types
	Org Format = "org"
)

var (
	supportedFormats = map[Format]bool{
		Md:  true,
		Org: true,
	}
)

// ParseFormat parses note format from string. If format is unsupported, returns error
func ParseFormat(s string) (Format, error) {
	if s == "" {
		return Unknown, nil
	}
	t := Format(s)
	if supportedFormats[t] {
		return t, nil
	}
	return Unknown, fmt.Errorf("Unsupported note type: %s", s)
}
