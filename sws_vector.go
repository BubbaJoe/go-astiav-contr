package astiav

/*
#cgo pkg-config: libswscale
#include "libswscale/swscale.h"
*/
import "C"
import (
	"unsafe"
)

type SwsVector struct {
	c *C.struct_SwsVector
}

func NewSwsVector(size int) *SwsVector {
	return &SwsVector{
		c: C.sws_allocVec(C.int(size)),
	}
}

func (v *SwsVector) Free() {
	C.sws_freeVec(v.c)
}

func NewGaussianVector(variance, quality float64) *SwsVector {
	vec := C.sws_getGaussianVec(C.double(variance), C.double(quality))
	if vec == nil {
		return nil
	}
	return &SwsVector{
		c: vec,
	}
}

func (v *SwsVector) Scale(scalar float64) {
	C.sws_scaleVec(v.c, C.double(scalar))
}

func (v *SwsVector) Normalize(height float64) {
	C.sws_normalizeVec(v.c, C.double(height))
}

func (v *SwsVector) Coefficients() []float64 {
	result := make([]float64, int(v.c.length))
	for _, d := range unsafe.Slice(v.c.coeff, v.c.length) {
		result = append(result, float64(d))
	}
	return result
}
