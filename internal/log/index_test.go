package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewIndex(t *testing.T) {
	f, err := os.CreateTemp("", "new_index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	index, err := newIndex(f, c)
	require.NoError(t, err)
	fi, err := os.Stat(index.Name())
	require.Equal(t, uint64(fi.Size()), c.Segment.MaxIndexBytes)
}
