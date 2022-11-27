package astiav

/*
#cgo pkg-config: libswscale
#include "libswscale/swscale.h"
*/
import "C"

type SwsFilter struct {
	c *C.struct_SwsFilter
}

func SwsDefaultFilter(lumaGBlur, chromaGBlur, lumaSharpen, chromaSharpen, chromaHShift, chromaVShift float64, verbose int) *SwsFilter {
	return &SwsFilter{
		c: C.sws_getDefaultFilter(
			C.float(lumaGBlur),
			C.float(chromaGBlur),
			C.float(lumaSharpen),
			C.float(chromaSharpen),
			C.float(chromaHShift),
			C.float(chromaVShift),
			C.int(verbose),
		),
	}
}

func (f *SwsFilter) Free() {
	if f.c != nil {
		C.sws_freeFilter(f.c)
	}
}
