package astiav

/*
#cgo pkg-config: libavutil
#include "libavutil/opt.h"
#include "libavutil/log.h"
*/
import "C"
import "unsafe"

type Class struct {
	c *C.struct_AVClass
}

func (c Class) Name() string {
	return C.GoString(c.c.class_name)
}

type Classer interface {
	Class() Class
}

func SetOptionsFromDictionary(o Classer, d *Dictionary) error {
	cptr := unsafe.Pointer(o.Class().c)
	return newError(C.av_opt_set_dict(cptr, &d.c))
}
