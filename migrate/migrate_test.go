package migrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSolve(t *testing.T) {
	t.Run("Down", func(t *testing.T) {
		dir, path, err := solve(5, 0, []int{1, 2, 3, 4, 5, 6})
		require.Nil(t, err)
		require.Equal(t, down, dir)
		require.Equal(t, []int{5, 4, 3, 2, 1}, path)
	})

	t.Run("Up", func(t *testing.T) {
		dir, path, err := solve(0, 5, []int{1, 2, 3, 4, 5})
		require.Nil(t, err)
		require.Equal(t, up, dir)
		require.Equal(t, []int{1, 2, 3, 4, 5}, path)
	})

	t.Run("UpOne", func(t *testing.T) {
		dir, path, err := solve(3, 4, []int{1, 2, 3, 4, 5})
		require.Nil(t, err)
		require.Equal(t, up, dir)
		require.Equal(t, []int{4}, path)
	})

	t.Run("DownOne", func(t *testing.T) {
		dir, path, err := solve(4, 3, []int{1, 2, 3, 4, 5})
		require.Nil(t, err)
		require.Equal(t, down, dir)
		require.Equal(t, []int{4}, path)
	})

	t.Run("Equal", func(t *testing.T) {
		dir, path, err := solve(4, 4, []int{1, 2, 3, 4, 5})
		require.Nil(t, err)
		require.Equal(t, none, dir)
		require.Empty(t, path)
	})

	t.Run("NonExistantDesired", func(t *testing.T) {
		_, _, err := solve(4, 10, []int{1, 2, 3, 4, 5})
		require.NotNil(t, err)
	})

	t.Run("NonExistantCurrent", func(t *testing.T) {
		_, _, err := solve(10, 4, []int{1, 2, 3, 4, 5})
		require.NotNil(t, err)
	})

	t.Run("ZeroVersion", func(t *testing.T) {
		_, _, err := solve(1, 4, []int{0, 1, 2, 3, 4, 5})
		require.NotNil(t, err)
	})
}
