package astiav_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestIOContext_ReadWriteSeek(t *testing.T) {
	c := astiav.NewIOContext()
	path := filepath.Join(t.TempDir(), "iocontext.txt")

	// Write Test
	err := c.Open(path, astiav.NewIOContextFlags(
		astiav.IOContextFlagReadWrite))
	require.NoError(t, err)
	fmt.Println("got here", c.Size())

	err = c.Write(nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), c.Size())
	fmt.Println("got here", c.Size())

	err = c.Write([]byte("testtest"))
	c.Flush()
	require.Equal(t, int64(8), c.Size())
	require.NoError(t, err)
	fmt.Println("got here", c.Size())

	err = c.Closep()
	require.NoError(t, err)

	// Read Test
	c = astiav.NewIOContext()
	// using IOContextFlagReadWrite doesn't seem to work here.
	err = c.Open(path, astiav.NewIOContextFlags(astiav.IOContextFlagRead))
	require.NoError(t, err)
	// fmt.Println("got here 1", c.Size())

	// Seek Test
	// require.True(t, c.Seekable())
	// i, err := c.Seek(4, io.SeekStart)
	// require.Equal(t, int64(4), i)
	// require.NoError(t, err)
	// fmt.Println("got here", path, c.Size())

	d := make([]byte, 4)
	j, err := c.Read(d)
	require.NoError(t, err)
	require.Equal(t, 4, j)
	require.Equal(t, "test", string(d[:j]))
	fmt.Println("got here", c.Size())

	// Cleanup
	err = c.Closep()
	require.NoError(t, err)

	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "testtest", string(b))

	err = os.Remove(path)
	require.NoError(t, err)
}

func TestIOContext_Open2Write(t *testing.T) {
	c := astiav.NewIOContext()
	path := filepath.Join(t.TempDir(), "iocontext.txt")

	// Write Test
	dict := astiav.NewDictionary()
	defer dict.Free()
	dict.Set("test", "test", 0)
	err := c.Open2(c, path, astiav.NewIOContextFlags(
		astiav.IOContextFlagReadWrite), dict)
	require.NoError(t, err)

	err = c.Write(nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), c.Size())

	err = c.Write([]byte("testtest"))
	c.Flush()
	require.Equal(t, int64(8), c.Size())
	require.NoError(t, err)

	// Cleanup
	err = c.Closep()
	require.NoError(t, err)

	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "testtest", string(b))

	err = os.Remove(path)
	require.NoError(t, err)
}

// func TestIOContext_ReadWriteCopy_(t *testing.T) {
// 	cx := astiav.NewIOContext().Wrapper()
// 	path := filepath.Join(t.TempDir(), "iocontext.txt")

// 	// Write Test
// 	err := cx.Open(path, astiav.NewIOContextFlags(
// 		astiav.IOContextFlagReadWrite))
// 	require.NoError(t, err)
// 	fmt.Println("got x", path, cx.Size())

// 	n, err := cx.Write(nil)
// 	require.NoError(t, err)
// 	require.Equal(t, int64(0), n)
// 	fmt.Println("got x", cx.Size())

// 	l, err := cx.Write([]byte("testtest"))
// 	require.Equal(t, int64(8), l)
// 	require.NoError(t, err)
// 	fmt.Println("got x", cx.Size())

// 	err = cx.Close()
// 	require.NoError(t, err)

// 	// Read Test
// 	c := astiav.NewIOContext().Wrapper()

// 	// using IOContextFlagReadWrite doesn't seem to work here.
// 	err = c.Open(path, astiav.NewIOContextFlags(astiav.IOContextFlagRead))
// 	require.NoError(t, err)
// 	fmt.Println("got here 1", c.Size())

// 	// Seek Test
// 	require.True(t, c.Seekable())
// 	i, err := c.Seek(0, io.SeekStart)
// 	require.Equal(t, int64(0), i)
// 	require.NoError(t, err)
// 	fmt.Println("got x", c.Size())

// 	buf := bytes.NewBuffer(make([]byte, 8))
// 	j, err := c.Read(buf.Bytes())
// 	require.NoError(t, err)
// 	require.Equal(t, 8, j)
// 	require.Equal(t, "testtest", buf.String())
// 	fmt.Println("got x", c.Size())

// 	// Cleanup
// 	err = c.Close()
// 	require.NoError(t, err)
// 	fmt.Println("got y", c.Size())

// 	b, err := os.ReadFile(path)
// 	require.NoError(t, err)
// 	require.Equal(t, "testtest", string(b))

// 	err = os.Remove(path)
// 	require.NoError(t, err)
// }

func TestIOContext_ReaderWriter(t *testing.T) {
	// buf := bytes.NewBuffer(make([]byte, 1024))
	// c := astiav.AllocIOContext(
	// 	buf.Bytes(), true, func(i1 *interface{}, b []byte, i2 int) int {
	// 		fmt.Println("read", i1, b, i2)
	// 		return 0
	// 	}, func(i1 *interface{}, b []byte, i2 int) int {
	// 		fmt.Println("write", i1, b, i2)
	// 		return 0
	// 	}, func(i1 *interface{}, i2 int64, i3 int) int64 {
	// 		fmt.Println("seek", i1, i2, i3)
	// 		return 0
	// 	},
	// )
	c := astiav.NewIOContext()
	defer c.Closep()

	// c.Open("", astiav.NewIOContextFlags(astiav.IOContextFlagReadWrite))

	err := c.Write([]byte("testtest"))
	require.NoError(t, err)

	require.True(t, c.Seekable())
	// j, err := c.Seek(0, io.SeekStart)
	// require.NoError(t, err)
	// require.Equal(t, int64(0), j)

	// data := make([]byte, 8)
	// n, err := c.Read(data)

	// require.NoError(t, err)
	// require.Equal(t, 8, n)

}
