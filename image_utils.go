package astiav

/*
#cgo pkg-config: libavutil
#include <libavutil/frame.h>
#include <libavutil/imgutils.h>



*/
import "C"
import "unsafe"

// func ImageFillArrays(
// 	dstData [][]byte, dstLinesize []int, src []byte,
// 	pixFmt PixelFormat, width, height, align int,
// ) error {
// 	dstDataC := (*C.uint8_t)(unsafe.Pointer(&dstData[0]))
// 	dstLinesizeC := (*C.int)(unsafe.Pointer(&dstLinesize[0]))
// 	var srcC = (*C.uint8_t)(unsafe.Pointer(&src[0]))
// 	return newError(C.av_image_fill_arrays(
// 		&dstDataC, dstLinesizeC, srcC,
// 		C.enum_AVPixelFormat(pixFmt),
// 		C.int(width), C.int(height), C.int(align)))
// }

func ImageFillFrameArrays(
	frame *Frame, srcSize int, pixFmt PixelFormat, width, height, align int,
) error {
	if frame == nil || frame.c == nil {
		panic("frame is nil")
	}
	return newError(C.av_image_fill_arrays(
		(**C.uchar)(unsafe.Pointer(&frame.c.data)),
		(*C.int)(unsafe.Pointer(&frame.c.linesize)),
		(*C.uchar)(C.malloc(C.size_t(srcSize))),
		C.enum_AVPixelFormat(pixFmt),
		C.int(width),
		C.int(height),
		C.int(align),
	))
}
