package remote

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const actualConfig = `
# comment

rsync://user@example.com:bak/name/

git://git@github.com:nchern/go-codegen.git
`

const malformedConfig = `
# comment

foo bar
`

func TestParseConfig(t *testing.T) {
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

func TestParseMalformedConfig(t *testing.T) {
	actualRemotes, err := parse(bytes.NewBufferString(malformedConfig))
	assert.Error(t, err)
	assert.Nil(t, actualRemotes)
}
