package log

import (
	"io"
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
	require.Equal(t, f.Name(), index.Name())
	entries := []struct {
		Off uint32
		Pos uint64
	}{
		{Off: 0, Pos: 0},
		{Off: 1, Pos: 10},
	}
	for _, want := range entries {
		err = index.Write(want.Off, want.Pos)
		require.NoError(t, err)

		_, pos, err := index.Read(int64(want.Off))
		require.NoError(t, err)
		require.Equal(t, want.Pos, pos)
	}
}

func TestWriteWithoutSpace(t *testing.T) {
	f, err := os.CreateTemp("", "write_without_space_index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	c := Config{}
	c.Segment.MaxIndexBytes = 0
	i, err := newIndex(f, c)
	require.NoError(t, err)
	err = i.Write(uint32(2), uint64(2))
	require.ErrorIs(t, err, io.EOF)
}

func TestWriteWithSpace(t *testing.T) {
	f, err := os.CreateTemp("", "write_with_space_index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	i, err := newIndex(f, c)
	require.NoError(t, err)
	err = i.Write(0, 1)
	require.NoError(t, err)
	out, pos, err := i.Read(0)
	require.NoError(t, err)
	require.Equal(t, out, uint32(0))
	require.Equal(t, pos, uint64(1))
}
