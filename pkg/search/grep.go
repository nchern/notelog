package search

// GrepEngine uses plain regular expressions to search through the notes,
// resembles simple grep
type GrepEngine struct {
	Engine
}

// NewGrepEngine returns a new GrepEngine instance
func NewGrepEngine(notes Notes) *GrepEngine {
	return &GrepEngine{
		Engine{notes: notes},
	}
}

// Search runs grep search
func (s *GrepEngine) Search(terms ...string) ([]*Result, error) {
	expr := ""
	if len(terms) > 0 {
		expr = terms[0]
	}

	rx, err := compileRx(expr, !s.CaseSensitive)
	if err != nil {
		return nil, err
	}

	matches := func(s string) []string {
		return rx.FindAllString(s, -1)
	}

	return searchInNotes(s.notes, matches, s.OnlyNames)
}
