package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberedWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	underTest := &nlWriter{inner: buf}

	given := text(
		"foo",
		"bar",
		"buzz")

	actualN, err := fmt.Fprint(underTest, given)
	assert.NoError(t, err)

	assert.Equal(t, len(given), actualN)
	assert.Equal(t,
		text("1. foo",
			"2. bar",
			"3. buzz"),
		buf.String())
}
