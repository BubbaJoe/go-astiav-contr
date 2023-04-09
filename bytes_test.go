package astiav

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytes(t *testing.T) {
	// TODO Test stringFromC
	// TODO Test bytesFromC
	// TODO Test bytesToC
}

// func Test_bytesToC(t *testing.T) {
// 	wasRan := false
// 	bb := []byte{1, 2, 3, 4, 5, 6, 7, 8}
// 	bytesToC(bb, func(p unsafe.Pointer, size int) error {
// 		for i, b := range bb {
// 			require.Equal(t, b, *(*byte)(
// 				offsetPtr(unsafe.Pointer(p), i)))
// 		}
// 		wasRan = true
// 		return nil
// 	})
// 	require.True(t, wasRan)
// }

func Test_copyGoBytes(t *testing.T) {
	ptr8, free := mallocBytes(8)
	defer free()
	bb := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	err := copyGoBytes(ptr8, bb)
	require.NoError(t, err)

	for i, b := range bb {
		require.Equal(t, b, *(*byte)(offsetPtr(ptr8, i)))
	}
}

func Test_copyGoBytes2(t *testing.T) {
	ptr8, free := mallocBytes(8)
	defer free()
	buf := new(bytes.Buffer)

	ints := []int32{1, 2, 3, 4, 5, 6, 7, 8}
	for _, i := range ints {
		binary.Write(buf, binary.LittleEndian, i)
	}
	err := copyGoBytes(ptr8, buf.Bytes())
	require.NoError(t, err)

	for i, b := range ints {
		require.Equal(t, b, *(*int32)(offsetPtr(ptr8, i*4)))
	}
}
