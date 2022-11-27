package astiav

/*
#cgo pkg-config: libavformat
#include <libavformat/avformat.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type ioContextCbs struct {
	readCb  IOContextReadFunc
	writeCb IOContextWriteFunc
	seekCb  IOContextSeekFunc
}

//export go_ioctx_proxy_read
func go_ioctx_proxy_read(opaque unsafe.Pointer, buf *C.uint8_t, buf_size C.int) C.int {
	id := int(*(*C.int)(opaque))
	if ctx, ok := fetchIOCallback(id); ok {
		gobuf := make([]byte, int(buf_size))
		n := ctx.readCb(gobuf)
		cn := C.int(n)
		if n < 0 {
			return C.int(cn)
		} else if n == 0 {
			// returning 0 throws error in ffmpeg
			return C.int(ErrUnknown)
		}
		for i := 0; i < n; i++ {
			curBuf := (*C.uint8_t)(offsetPtr(unsafe.Pointer(buf), i))
			*curBuf = C.uint8_t(gobuf[i])
		}
		// cbuf := C.CBytes(gobuf[:n])
		// defer C.free(cbuf)
		// fmt.Println("GG1:", buf, &buf, cbuf, &cbuf)
		// C.memcpy(unsafe.Pointer(&buf), cbuf, C.size_t(n))
		// C.memcpy(unsafe.Pointer(buf), cbuf, C.size_t(n))
		// fmt.Println("GG2:", buf, &buf, cbuf, &cbuf)
		return cn
	}
	fmt.Println("RETURNING ERR:", ErrEio)
	return C.int(ErrEio)
}

//export go_ioctx_proxy_write
func go_ioctx_proxy_write(opaque unsafe.Pointer, buf *C.uint8_t, buf_size C.int) C.int {
	id := int(*(*C.int)(opaque))
	// fmt.Println("SOURCE WRITE:", string(*buf), buf_size, payload.buf_size)
	if ctx, ok := fetchIOCallback(id); ok {
		return C.int(ctx.writeCb(C.GoBytes(unsafe.Pointer(buf), buf_size)))
	}
	return C.int(ErrEio)
}

//export go_ioctx_proxy_seek
func go_ioctx_proxy_seek(opaque unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	id := int(*(*C.int)(opaque))
	if ctx, ok := fetchIOCallback(id); ok {
		return C.int64_t(ctx.seekCb(int64(offset), int(whence)))
	}
	return C.int64_t(-1)
}

// func readCallback(buf []byte) int {}
