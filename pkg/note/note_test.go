package note

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testToday = "2022-07-30"

func init() {
	now, err := time.Parse(dateFormat, testToday)
	if err != nil {
		panic(err)
	}

	nowFn = func() time.Time { return now }
}

func TestRemoveIfEmpty(t *testing.T) {
	withNotes(t, func(notes List) {

		underTest, err := notes.GetOrCreate("empty", Org)
		require.NoError(t, err)

		found, err := underTest.Exists()
		require.NoError(t, err)
		require.True(t, found)

		assert.NoError(t, underTest.RemoveIfEmpty())

		_, err = os.Stat(underTest.dir())
		require.True(t, os.IsNotExist(err))
	})
}
