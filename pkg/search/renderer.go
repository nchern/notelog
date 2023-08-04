package search

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Renderer renders given search results
type Renderer interface {
	Render(*Result) error
	Close() error
}

// Render renders a results collection using given renderer
func Render(renderer Renderer, results []*Result, onlyNames bool) error {
	defer renderer.Close()
	for _, res := range results {
		if err := renderer.Render(res); err != nil {
			return err
		}
	}
	return nil
}

// StreamRenderer renders search result to the underline stream
type StreamRenderer struct {
	OnlyNames bool
	Colorize  bool
	W         io.Writer
}

// Render renders search results
func (sr *StreamRenderer) Render(res *Result) (err error) {
	rs := *res
	if sr.OnlyNames {
		rs.text = " "
		rs.lineNum = 1
	}
	_, err = fmt.Fprintln(sr.W, rs.Display(sr.Colorize))
	return
}

// Close closes the underline stream
func (sr *StreamRenderer) Close() error { return nil }

type persistentRenderer struct {
	rndr Renderer
	w    io.WriteCloser
}

// NewPersistentRenderer returns a new instance of a Renderer
// that wraps a given renderer and also dumps the results into a lastResults file
func NewPersistentRenderer(notes Notes, r Renderer) (Renderer, error) {
	f, err := os.Create(notes.MetadataFilename(lastResultsFile))
	if err != nil {
		return nil, err
	}
	return &persistentRenderer{w: f, rndr: r}, nil
}

// Render renders search results
func (pr *persistentRenderer) Render(res *Result) (err error) {
	if err = pr.rndr.Render(res); err != nil {
		return
	}
	_, err = fmt.Fprintln(pr.w, res)
	return
}

// Close closes all underlying streams; best effort
func (pr *persistentRenderer) Close() error {
	if err := pr.rndr.Close(); err != nil {
		log.Printf("ERROR persistentRenderer: %s", err)
	}
	return pr.w.Close()
}
