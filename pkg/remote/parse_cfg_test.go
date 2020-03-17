package remote

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const actualConfig = `
# comment

rsync://user@example.com:bak/name/

git://git@github.com:nchern/go-codegen.git
`

const (
	malformedConfig = `# comment

foo bar
`

	noRemotesConfig = `# comment

`
)

func TestShouldParseConfig(t *testing.T) {
	actual, err := parse(bytes.NewBufferString(actualConfig))
	assert.NoError(t, err)
	assert.Len(t, actual, 2)

	assert.Equal(t,
		[]*entry{
			&entry{Scheme: "rsync", Addr: "user@example.com:bak/name/"},
			&entry{Scheme: "git", Addr: "git@github.com:nchern/go-codegen.git"},
		},
		actual)
}

func TestShouldFailtToParseMalformedConfig(t *testing.T) {
	var tests = []struct {
		name     string
		expected error
		given    string
	}{
		{"malformed line", errConfigMalformed, malformedConfig},
		{"no remotes", ErrConfigEmpty, noRemotesConfig},
		{"empty", ErrConfigEmpty, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parse(bytes.NewBufferString(tt.given))
			assert.Nil(t, actual)
			assert.True(t, errors.Is(err, tt.expected))
		})
	}
}
