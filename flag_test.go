package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlag(t *testing.T) {
	f1 := flags(2 | 4 | 16)
	require.Equal(t, flags(0b10110), f1)
	f2 := flags(f1.add(1))
	require.Equal(t, flags(0b10111), f2)
	f3 := flags(f2.del(2))
	require.Equal(t, flags(0b10101), f3)

	require.True(t, f3.has(1))
	require.False(t, f3.has(2))
	require.True(t, f3.has(4))
	require.False(t, f3.has(8))
	require.True(t, f3.has(16))
	require.False(t, f3.has(32))
}
