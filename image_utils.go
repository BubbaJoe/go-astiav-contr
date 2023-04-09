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
	frame *Frame, src []byte,
	width, height, align int,
) error {
	if frame == nil || frame.c == nil {
		panic("frame is nil")
	}
	return newError(C.av_image_fill_arrays(
		(**C.uchar)(unsafe.Pointer(&frame.c.data)),
		(*C.int)(unsafe.Pointer(&frame.c.linesize)),
		(*C.uchar)(unsafe.Pointer(&src[0])),
		int32(frame.PixelFormat()),
		C.int(width),
		C.int(height),
		C.int(align),
	))
}

func ImageCopyFrameToBuffer(
	buffer_size int, frame *Frame, src []byte,
	pixFmt PixelFormat, width, height, align int,
) ([]byte, error) {
	return nil, nil
}

func ImageGetBufferSize(
	pixFmt PixelFormat, width, height, align int,
) (int, error) {
	ret := C.av_image_get_buffer_size(
		C.enum_AVPixelFormat(pixFmt),
		C.int(width), C.int(height), C.int(align))
	if ret < 0 {
		return 0, newError(ret)
	}
	return int(ret), nil
}

func ImageFillFrameBlack(
	frame *Frame, width, height, align int,
) error {
	return nil
}

func ImageGetLinesize(
	pxlFmt PixelFormat,
	width, plane int,
) (int, error) {
	ret := C.av_image_get_linesize(
		C.enum_AVPixelFormat(pxlFmt),
		C.int(width), C.int(plane),
	)
	return int(ret), newError(ret)
}
