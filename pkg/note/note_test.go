package note

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveIfEmpty(t *testing.T) {
	withNotes(t, func(notes List) {

		underTest, err := notes.GetOrCreate("empty")
		require.NoError(t, err)

		found, err := underTest.Exists()
		require.NoError(t, err)
		require.True(t, found)

		assert.NoError(t, underTest.RemoveIfEmpty())

		_, err = os.Stat(underTest.dir())
		require.True(t, os.IsNotExist(err))
	})
}
