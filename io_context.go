package astiav

/*
#cgo pkg-config: libavformat
#include <libavformat/avformat.h>
*/
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avio.h#L161
type IOContext struct {
	c *C.struct_AVIOContext
}

func NewIOContext() *IOContext {
	return &IOContext{}
}

func newIOContextFromC(c *C.struct_AVIOContext) *IOContext {
	if c == nil {
		return nil
	}
	return &IOContext{c: c}
}

func (ic *IOContext) Free() {
	C.avio_context_free(&ic.c)
}

func (ic *IOContext) Closep() error {
	if ic.c == nil {
		return newError(C.avio_closep(&ic.c))
	}
	return nil
}

// 	var buf *C.uint8_t
// 	var bufSize C.int
// 	if err := newError(C.avio_close_dyn_buf(ic.c, &buf, &bufSize)); err != nil {
// 		return nil, err
// 	}
// 	b := C.GoBytes(unsafe.Pointer(buf), bufSize)
// 	C.av_free(unsafe.Pointer(buf))
// 	return b, nil
// }

func (ic *IOContext) Open(filename string, flags IOContextFlags) error {
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	return newError(C.avio_open(&ic.c, cfi, C.int(flags)))
}

func (ic *IOContext) Accept(client *IOContext) error {
	return newError(C.avio_accept(ic.c, &client.c))
}

func (ic *IOContext) Handshake() error {
	return newError(C.avio_handshake(ic.c))
}

func (ic *IOContext) Open2(pb *IOContext, filename string, flags IOContextFlags, opts *Dictionary) error {
	if pb == nil {
		return newError(-1)
	}
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	var copts **C.struct_AVDictionary
	if opts != nil {
		copts = &opts.c
	}
	return newError(C.avio_open2(&pb.c, cfi, C.int(flags), nil, copts))
}

func (ic *IOContext) EOFReached() bool {
	return int(ic.c.eof_reached) != 0
}

func (ic *IOContext) Write(b []byte) error {
	if b == nil {
		return nil
	}
	C.avio_write(ic.c, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
	return nil
}

func (ic *IOContext) Flush() {
	C.avio_flush(ic.c)
}

func (ic *IOContext) Read(b []byte) (int, error) {
	ret := C.avio_read(ic.c, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
	if ret < 0 {
		return 0, newError(ret)
	}
	return int(ret), nil
}

func (ic *IOContext) Seekable() bool {
	return int(ic.c.seekable) != 0
}

func (ic *IOContext) Seek(offset int64, whence int) (int64, error) {
	ret := C.avio_seek(ic.c, C.int64_t(offset), C.int(whence))
	if ret < 0 {
		return 0, newError(C.int(ret))
	}
	return int64(ret), nil
}

func (ic *IOContext) CurrentPosition() int64 {
	return int64(ic.c.pos)
}

func (ic *IOContext) Size() int64 {
	return int64(C.avio_size(ic.c))
}

// type IOContextWrapper struct {
// 	ic *IOContext
// }

// type IOContextWrapper struct {
// 	ic *IOContext
// }

// func (ic *IOContext) Wrapper() *IOContextWrapper {
// 	return &IOContextWrapper{ic: ic}
// }

// func (w *IOContextWrapper) Close() error {
// 	return w.ic.Closep()
// }

// func (w *IOContextWrapper) Read(b []byte) (int, error) {
// 	l, err := w.ic.Read(b)
// 	if err != nil {
// 		return 0, io.EOF
// 	}
// 	return l, nil
// }

// func (w *IOContextWrapper) Write(b []byte) (int, error) {
// 	err := w.ic.Write(b)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return len(b), nil
// }

// func (w *IOContextWrapper) Open(filename string, flags IOContextFlags) error {
// 	return w.ic.Open(filename, flags)
// }

// func (w *IOContextWrapper) Open2(filename string, flags IOContextFlags, dict *Dictionary) error {
// 	return w.ic.Open2(filename, flags, dict)
// }

// func (w *IOContextWrapper) Seek(offset int64, whence int) (int64, error) {
// 	return w.ic.Seek(offset, whence)
// }

// func (w *IOContextWrapper) Flush() {
// 	w.ic.Flush()
// }

// func (w *IOContextWrapper) Size() int64 {
// 	return w.ic.Size()
// }

// func (w *IOContextWrapper) Seekable() bool {
// 	return w.ic.Seekable()
// }
