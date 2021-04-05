package note

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveIfEmpty(t *testing.T) {
	withNotes(t, func(notes List) {

		underTest := NewNote("empty", notes.HomeDir())
		found, err := underTest.Exists()
		require.NoError(t, err)
		require.False(t, found)

		assert.NoError(t, underTest.Init())

		_, err = os.Stat(underTest.dir())
		require.NoError(t, err)

		assert.NoError(t, underTest.RemoveIfEmpty())

		_, err = os.Stat(underTest.dir())
		require.True(t, os.IsNotExist(err))
	})
}
