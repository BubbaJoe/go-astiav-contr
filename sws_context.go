package astiav

/*
#cgo pkg-config: libswscale
#include "libswscale/swscale.h"
*/
import "C"
import (
	"unsafe"
)

var (
	SWS_FAST_BILINEAR int = C.SWS_FAST_BILINEAR
	SWS_BILINEAR      int = C.SWS_BILINEAR
	SWS_BICUBIC       int = C.SWS_BICUBIC
	SWS_X             int = C.SWS_X
	SWS_POINT         int = C.SWS_POINT
	SWS_AREA          int = C.SWS_AREA
	SWS_BICUBLIN      int = C.SWS_BICUBLIN
	SWS_GAUSS         int = C.SWS_GAUSS
	SWS_SINC          int = C.SWS_SINC
	SWS_LANCZOS       int = C.SWS_LANCZOS
	SWS_SPLINE        int = C.SWS_SPLINE
)

type SwsContext struct {
	c          *C.struct_SwsContext
	srcW, srcH int
	dstW, dstH int
	srcFormat  int
	dstFormat  int
	srcFilter  *SwsFilter
	dstFilter  *SwsFilter
}

func (ctx *SwsContext) Class() *Class {
	return &Class{c: C.sws_get_class()}
}

func NewSwsContext(
	srcW, srcH int, srcPixFmt PixelFormat,
	dstW, dstH int, dstPixFmt PixelFormat,
	flags int, srcFilter, dstFilter *SwsFilter,
	params ...float64,
) *SwsContext {
	// float64 array to *C.double
	var cparams *C.double = nil
	if params != nil && len(params) > 0 {
		floatSize := 8
		cparams := C.malloc(C.size_t(len(params) * floatSize))
		for i, p := range params {
			curPtr := offsetPtr(cparams, i*floatSize)
			*(*float64)(curPtr) = p
			// *(*C.double)(curPtr) = C.double(p)
		}
	}
	var cSrcFilter, cDstFilter *C.struct_SwsFilter
	if srcFilter != nil {
		cSrcFilter = srcFilter.c
	}
	if dstFilter != nil {
		cDstFilter = dstFilter.c
	}
	sws := C.sws_getContext(
		C.int(srcW),
		C.int(srcH),
		int32(srcPixFmt),
		C.int(dstW),
		C.int(dstH),
		int32(dstPixFmt),
		C.int(flags),
		cSrcFilter,
		cDstFilter,
		cparams,
	)

	if sws == nil {
		return nil
	}

	return &SwsContext{
		c:         sws,
		dstW:      dstW,
		dstH:      dstH,
		srcW:      srcW,
		srcH:      srcH,
		srcFormat: int(srcPixFmt),
		dstFormat: int(dstPixFmt),
		srcFilter: srcFilter,
		dstFilter: dstFilter,
	}
}

func (ctx *SwsContext) InitContext(srcFilter, dstFilter *SwsFilter) error {
	err := newError(C.sws_init_context(ctx.c, srcFilter.c, dstFilter.c))
	if err != nil {
		return err
	}
	ctx.srcFilter = srcFilter
	ctx.dstFilter = dstFilter
	return nil
}

func (ctx *SwsContext) CachedContext(
	srcW, srcH int, srcPixFmt PixelFormat,
	dstW, dstH int, dstPixFmt PixelFormat,
	flag int,
) {
	ctx.c = C.sws_getCachedContext(
		ctx.c,
		C.int(srcW),
		C.int(srcH),
		int32(srcPixFmt),
		C.int(dstW),
		C.int(dstH),
		int32(dstPixFmt),
		C.int(flag),
		ctx.srcFilter.c,
		ctx.dstFilter.c,
		nil,
	)

	ctx.dstH = dstH
	ctx.dstW = dstW
	ctx.srcH = srcH
	ctx.srcW = srcW
	ctx.dstFormat = int(dstPixFmt)
	ctx.srcFormat = int(srcPixFmt)
}

type SwsColorspaceDetails struct {
	InvTable   [4]int
	SrcRange   int
	Table      [4]int
	DstRange   int
	Brightness int
	Contrast   int
	Saturation int
}

// Set color space details
func (ctx *SwsContext) SetColorspaceDetails(
	invTable [4]int, srcRange int,
	table [4]int, dstRange, brightness, contrast, saturation int,
) {
	C.sws_setColorspaceDetails(ctx.c, (*C.int)(unsafe.Pointer(&invTable)),
		C.int(srcRange), (*C.int)(unsafe.Pointer(&table)), C.int(dstRange),
		C.int(brightness), C.int(contrast), C.int(saturation))
}

func (ctx *SwsContext) ColorspaceDetails() *SwsColorspaceDetails {
	var (
		invTable   [4]int
		table      [4]int
		srcRange   int
		dstRange   int
		brightness int
		contrast   int
		saturation int
	)

	C.sws_getColorspaceDetails(ctx.c, (**C.int)(unsafe.Pointer(&invTable)),
		(*C.int)(unsafe.Pointer(&srcRange)), (**C.int)(unsafe.Pointer(&table)),
		(*C.int)(unsafe.Pointer(&dstRange)), (*C.int)(unsafe.Pointer(&brightness)),
		(*C.int)(unsafe.Pointer(&contrast)), (*C.int)(unsafe.Pointer(&saturation)))

	return &SwsColorspaceDetails{
		InvTable:   invTable,
		SrcRange:   srcRange,
		Table:      table,
		DstRange:   dstRange,
		Brightness: brightness,
		Contrast:   contrast,
		Saturation: saturation,
	}
}

func (ctx *SwsContext) SrcWidth() int {
	return int(ctx.srcW)
}

func (ctx *SwsContext) SrcHeight() int {
	return int(ctx.srcH)
}

func (ctx *SwsContext) SrcPixelFormat() PixelFormat {
	return PixelFormat(int(ctx.srcFormat))
}

func (ctx *SwsContext) Width() int {
	return int(ctx.dstW)
}

func (ctx *SwsContext) Height() int {
	return int(ctx.dstH)
}

func (ctx *SwsContext) PixelFormat() PixelFormat {
	return PixelFormat(int(ctx.dstFormat))
}

// Setters for the src and dst width and height
func (ctx *SwsContext) SetSrcWidth(w int) {
	ctx.srcW = w
}

func (ctx *SwsContext) SetSrcHeight(h int) {
	ctx.srcH = h
}

func (ctx *SwsContext) SetWidth(w int) {
	ctx.dstW = w
}

func (ctx *SwsContext) SetHeight(h int) {
	ctx.dstH = h
}

func (ctx *SwsContext) SetSrcPixelFormat(pixFmt PixelFormat) {
	ctx.srcFormat = int(pixFmt)
}

func (ctx *SwsContext) SetPixelFormat(pixFmt PixelFormat) {
	ctx.dstFormat = int(pixFmt)
}

func (ctx *SwsContext) ScaleFrames(src *Frame, dst *Frame) error {
	return newError(C.sws_scale(
		ctx.c,
		(**C.uint8_t)(unsafe.Pointer(&src.c.data)),
		(*C.int)(unsafe.Pointer(&src.c.linesize)),
		0, src.c.height,
		(**C.uint8_t)(unsafe.Pointer(&dst.c.data)),
		(*C.int)(unsafe.Pointer(&dst.c.linesize))))
}

func (ctx *SwsContext) Scale(
	srcSlice []byte, srcStride []int,
	srcSliceY, srcSliceH int,
	dstSlice []byte, dstStride []int,
) error {
	return newError(C.sws_scale(
		ctx.c,
		(**C.uint8_t)(unsafe.Pointer(&srcSlice[0])),
		(*C.int)(unsafe.Pointer(&srcStride[0])),
		C.int(srcSliceY),
		C.int(srcSliceH),
		(**C.uint8_t)(unsafe.Pointer(&dstSlice[0])),
		(*C.int)(unsafe.Pointer(&dstStride[0]))),
	)
}

func (ctx *SwsContext) ScaleDstFrame(
	srcSlice []byte, srcStride []int,
	srcSliceY, srcSliceH int, dstFrame *Frame,
) error {
	return newError(C.sws_scale(
		ctx.c, (**C.uint8_t)(unsafe.Pointer(&srcSlice[0])),
		(*C.int)(unsafe.Pointer(&srcStride[0])),
		C.int(srcSliceY), C.int(srcSliceH),
		(**C.uint8_t)(unsafe.Pointer(&dstFrame.c.data[0])),
		(*C.int)(unsafe.Pointer(&dstFrame.c.linesize[0])),
	))
}

func (ctx *SwsContext) ScaleSrcFrame(
	srcFrame *Frame, dstSlice []byte, dstStride []int,
) error {
	return newError(C.sws_scale(
		ctx.c,
		(**C.uint8_t)(unsafe.Pointer(&srcFrame.c.data)),
		(*C.int)(unsafe.Pointer(&srcFrame.c.linesize)),
		0, srcFrame.c.height,
		(**C.uint8_t)(unsafe.Pointer(&dstSlice[0])),
		(*C.int)(unsafe.Pointer(&dstStride[0]))),
	)
}

func (ctx *SwsContext) Free() {
	C.sws_freeContext(ctx.c)
}
