package astiav_test

import (
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestIOContext_Open_ReadWriteSeek(t *testing.T) {
	c := astiav.NewIOContext()
	path := filepath.Join(t.TempDir(), "iocontext.txt")

	// Write Test
	err := c.Open(path, astiav.NewIOContextFlags(
		astiav.IOContextFlagReadWrite))
	require.NoError(t, err)

	err = c.Write(nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), c.Size())

	err = c.Write([]byte("testtest"))
	c.Flush()
	require.Equal(t, int64(8), c.Size())
	require.NoError(t, err)

	err = c.Closep()
	require.NoError(t, err)

	// Read Test
	c = astiav.NewIOContext()
	// using IOContextFlagReadWrite, ...
	err = c.Open(path, astiav.NewIOContextFlags(astiav.IOContextFlagRead))
	require.NoError(t, err)

	// Seek Test
	// require.True(t, c.Seekable())
	// i, err := c.Seek(4, io.SeekStart)
	// require.Equal(t, int64(4), i)
	// require.NoError(t, err)

	d := make([]byte, 32768)
	j, err := c.Read(d)
	require.NoError(t, err)
	require.Equal(t, 8, j)
	require.Equal(t, "testtest", string(d[:j]))

	// Cleanup
	err = c.Closep()
	require.NoError(t, err)

	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "testtest", string(b))

	err = os.Remove(path)
	require.NoError(t, err)
}

func TestIOContext_OpenWith_Write(t *testing.T) {
	c := astiav.NewIOContext()
	path := filepath.Join(t.TempDir(), "iocontext.txt")

	// Write Test
	dict := astiav.NewDictionary()
	defer dict.Free()
	dict.Set("test", "test", 0)
	err := c.OpenWith(path, astiav.NewIOContextFlags(
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

func TestIOContext_Reader(t *testing.T) {
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
	buffer := randomBytes(1024 * 1024)

	c := astiav.AllocIOContextBufferReader(buffer)
	defer c.Closep()

	err := c.Write([]byte("testtest"))
	require.NoError(t, err)

	// fmt.Println("write done", c.Size())
	// c.Flush()
	// fmt.Println("flush done", c.Size())

	// require.True(t, c.Seekable())
	// i, err := c.Seek(0, io.SeekStart)
	// require.Equal(t, int64(0), i)

	buf := make([]byte, 256)
	n, err := c.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 256, n)
	// require.Equal(t, "testtest", string(buf))
	// j, err := c.Seek(0, io.SeekStart)
	// require.NoError(t, err)
	// require.Equal(t, int64(0), j)

	// data := make([]byte, 8)
	// n, err := c.Read(data)

	// require.NoError(t, err)
	// require.Equal(t, 8, n)

}

func TestIOContext_Callbacks_(t *testing.T) {
	readCount := 0
	writeCount := 0
	seekCount := 0

	b16 := randomBytes(16)
	b32 := randomBytes(32)
	b64 := randomBytes(64)

	c := astiav.AllocIOContextCallback(
		func(buf []byte) int {
			readCount += 1
			if readCount == 1 {
				copy(buf, b16)
				return len(b16)
			}
			return int(astiav.ErrEof)
		}, func(buf []byte) int {
			writeCount += 1
			if writeCount == 1 {
				require.Equal(t, b16, buf)
			} else if writeCount == 2 {
				require.Equal(t, b32, buf)
			} else if writeCount == 3 {
				require.Equal(t, b64, buf)
			}
			return len(buf)
		}, func(offset int64, whence int) int64 {
			seekCount += 1
			return offset
		},
	)
	defer c.Closep()

	err := c.Write(b16)
	require.NoError(t, err)
	require.Equal(t, 1, writeCount)

	err = c.Write(b32)
	require.NoError(t, err)
	require.Equal(t, 2, writeCount)

	err = c.Write(b64)
	require.NoError(t, err)
	require.Equal(t, 3, writeCount)

	require.True(t, c.Seekable())
	i, err := c.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), i)

	buf := make([]byte, 64)
	n, err := c.Read(buf)
	require.Equal(t, buf[:n], b16)
	require.NoError(t, err)
	require.Equal(t, 16, int(n))

	buf = make([]byte, 64)
	n, err = c.Read(buf)
	require.ErrorIs(t, err,
		(astiav.Error)(astiav.ErrEof))
	require.Equal(t, int(astiav.ErrEof), int(n))
}

func TestIOContext_Callbacks(t *testing.T) {
	byteArr := make([]byte, 64)
	size := 0
	pos := 0
	c := astiav.AllocIOContextCallback(
		func(buf []byte) int {
			min := len(buf)
			if pos >= size {
				return int(astiav.ErrEof)
			}
			if size < min {
				min = size
			}
			for i := 0; i < min; i++ {
				buf[i] = byteArr[pos+i]
			}
			// fmt.Println("READ BA:", pos, len(byteArr), byteArr)
			// fmt.Println("READ B:", pos, len(buf), buf)
			pos += min
			return min
		}, func(buf []byte) int {
			// fmt.Println("WR CB:", len(buf), buf)
			// fmt.Println("WR CB X:", pos, size, len(byteArr), byteArr)
			bufSize := len(buf)

			if pos >= len(byteArr) {
				return 0
			}
			if (pos + bufSize) > len(byteArr) {
				bufSize = (pos + bufSize) - len(byteArr)
			}
			for i := 0; i < bufSize; i++ {
				byteArr[pos+i] = buf[i]
			}
			pos += bufSize
			size += bufSize
			return bufSize
		}, func(offset int64, whence int) int64 {
			pos = int(offset)
			return offset
		})
	defer c.Closep()

	original := randomBytes(128)
	err := c.Write(original)
	require.NoError(t, err)

	c.Flush()

	require.True(t, c.Seekable())
	i, err := c.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), i)

	buf := make([]byte, 64)
	n, err := c.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 64, n)

	buf = make([]byte, 64)
	n, err = c.Read(buf)
	require.Equal(t, astiav.ErrEof, (astiav.Error)(n))
}

func randomBytes(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func BenchmarkIOContext_OpenMisc(b *testing.B) {
	astiav.SetLogLevel(astiav.LogLevelQuiet)
	file, err := os.Open("testdata/video.mp4")
	if err != nil {
		b.Fatal(err)
	}
	// b.Logf("file size: %d\n", len(file))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		fc := astiav.AllocFormatContext()
		defer fc.Free()
		ioCtx := astiav.AllocIOContextBufferReader(file)
		if ioCtx == nil {
			b.Fatal("ioCtx is nil")
		}
		defer ioCtx.Closep()
		fc.SetPb(ioCtx)
		dict1 := astiav.NewDictionary()
		if dict1 == nil {
			b.Fatal("dict is nil")
		}
		err = fc.OpenInput("", nil, dict1)
		if err != nil {
			b.Fatal(err)
		}
		dict2 := astiav.NewDictionary()
		if dict2 == nil {
			b.Fatal("dict is nil")
		}
		err := fc.FindStreamInfo(dict2)
		if err != nil {
			b.Fatal(err)
		}
		// if dict.Len() > 0 {
		// 	buf := make([]byte, 0, 1024*1024)
		// 	dict.Unpack(buf)
		// 	b.Logf("Dictionay Data: %s\n", string(buf))
		// }
		// decCc := astiav.AllocCodecContext(fc.Streams()[0].Codecpar())
		for _, is := range fc.Streams() {
			// Only process audio or video
			// b.Logf("found media type (%s) for stream%d", is.CodecParameters().MediaType().String(), is.Index())
			if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
				is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
				continue
			}
			// Find decoder
			// 		if dec := astiav.FindDecoder(is.CodecParameters().CodecID()); dec == nil {
			// 			b.Fatal("main: codec is nil")
			// 		} else {
			// 			b.Logf("main: found decoder for %s: %s\n", is.CodecParameters().CodecID().Name(), s.decCodec.Name())
			// 		}

			// 		// Find encode
			// 		if enc := astiav.FindEncoder(is.CodecParameters().CodecID()); enc == nil {
			// 			b.Fatal(errors.New("main: codec is nil"))
			// 		}

			// 		// Find pixel format
			// 		if is.CodecParameters().MediaType() == astiav.MediaTypeVideo {
			// 			b.Logf("pixel format: %s\n", is.CodecParameters().PixelFormat().Name())
			// 			b.Logf("color range id: %d\n", is.CodecParameters().ColorRange())
			// 			b.Logf("color space id: %d\n", is.CodecParameters().ColorSpace())
			// 			b.Logf("color primaries id: %d\n", is.CodecParameters().ColorPrimaries())
			// 		}

			// 		// Alloc decoder codec context
			// 		if cc = astiav.AllocCodecContext(dec); cc == nil {
			// 			b.Fatal("codec context is nil")
			// 		}
			// 		defer s.decCodecContext.Free()

			// 		// Update codec context
			// 		if err := is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
			// 			b.Fatal(fmt.Errorf("updating codec context failed: %w", err))
			// 		}

			// 		// Open codec context
			// 		if err := s.decCodecContext.Open(s.decCodec, nil); err != nil {
			// 			b.Fatal(fmt.Errorf("opening codec context failed: %w", err))
			// 		}

		}
		b.StopTimer()
	}
}
