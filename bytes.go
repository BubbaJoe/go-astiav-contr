package astiav

/*
#include <stdlib.h>
#include <stdio.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

func stringFromC(len int, fn func(buf *C.char, size C.size_t) error) (string, error) {
	size := C.size_t(len)
	buf := (*C.char)(C.malloc(size))
	if buf == nil {
		return "", errors.New("astiav: cannot alloc memory")
	}
	defer C.free(unsafe.Pointer(buf))
	if err := fn(buf, size); err != nil {
		return "", err
	}
	return C.GoString(buf), nil
}

func bytesFromC(fn func(size *C.int) *C.uint8_t) []byte {
	var size int
	r := fn((*C.int)(unsafe.Pointer(&size)))
	return C.GoBytes(unsafe.Pointer(r), C.int(size))
}

func bytesToC(b []byte, fn func(b *C.uint8_t, size C.int) error) error {
	var ptr *C.uint8_t
	if b != nil {
		cb := C.CBytes(b)
		if cb == nil {
			return errors.New("astiav: cannot alloc memory")
		}
		ptr = (*C.uint8_t)(cb)
	}
	return fn(ptr, C.int(len(b)))
}

func offsetPtr(ptr unsafe.Pointer, i int) unsafe.Pointer {
	if ptr == nil {
		panic("astiav: ptr is nil")
	}
	return unsafe.Pointer(uintptr(ptr) + uintptr(i))
}

func copyGoBytes(dst unsafe.Pointer, src []byte) error {
	if len(src) == 0 {
		return errors.New("astiav: src is empty")
	}
	ft := len(src)
	offsetPtr(dst, ft)
	memCopy(dst, unsafe.Pointer(&src[0]), C.size_t(len(src)))
	return nil
}

func copyCBytes(dst []byte, src unsafe.Pointer, size C.size_t) error {
	if len(dst) < int(size) {
		return errors.New("astiav: dst is too small")
	}
	memCopy(unsafe.Pointer(&dst[0]), src, size)
	return nil
}

func memCopy(dst, src unsafe.Pointer, size C.size_t) {
	memCopyInc(dst, src, size, 1)
}

func memCopyInc(dst, src unsafe.Pointer, size C.size_t, inc int) {
	for i := 0; i < int(size); i += inc {
		ptr := (*C.uint8_t)(offsetPtr(dst, i))
		if ptr == nil {
			panic("astiav: dst is nil")
		}
		if (*C.uint8_t)(offsetPtr(src, i)) == nil {
			panic("astiav: src is nil")
		}
		*ptr = *(*C.uint8_t)(offsetPtr(src, i))
	}
}

func mallocBytes(size int) (unsafe.Pointer, func()) {
	ptr := C.malloc(C.size_t(size))
	if ptr == nil {
		panic("astiav: cannot alloc memory")
	}
	return ptr, func() {
		C.free(ptr)
	}
}
