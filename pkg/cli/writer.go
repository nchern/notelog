package cli

import (
	"fmt"
	"io"
)

type nlWriter struct {
	inner         io.Writer
	ln            int
	numberWritten bool
}

func (w *nlWriter) Write(buf []byte) (int, error) {
	// not optimal implementation
	written := 0
	for _, b := range buf {
		n, err := w.write(b)
		written += n
		if err != nil {
			return written, err
		}
		if rune(b) == '\n' {
			w.numberWritten = false
		}
	}
	return written, nil
}

func (w *nlWriter) write(b byte) (int, error) {
	if !w.numberWritten {
		w.ln++
		_, err := fmt.Fprintf(w.inner, "%d. ", w.ln)
		if err != nil {
			return 0, err
		}
		w.numberWritten = true
	}
	n, err := w.inner.Write([]byte{b})
	return n, err
}
