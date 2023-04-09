package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/rational.h>
//#include <libavutil/mathematics.h>
import "C"
import (
	"strconv"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/rational.h#L58
type Rational struct {
	c C.struct_AVRational
}

func newRationalFromC(c C.struct_AVRational) Rational {
	return Rational{c: c}
}

func NewRational(num, den int) Rational {
	var r Rational
	r.SetNum(num)
	r.SetDen(den)
	return r
}

func (r Rational) Rescale(a int64, r2 Rational) int {
	return int(C.av_rescale_q(C.int64_t(a), r.c, r2.c))
}

func (r Rational) Compare(r2 Rational) int {
	return int(C.av_cmp_q(r.c, r2.c))
}

func (r Rational) Multiply(r2 Rational) Rational {
	return newRationalFromC(C.av_mul_q(r.c, r2.c))
}

func (r Rational) Div(r2 Rational) Rational {
	return newRationalFromC(C.av_div_q(r.c, r2.c))
}

func (r Rational) Add(r2 Rational) Rational {
	return newRationalFromC(C.av_add_q(r.c, r2.c))
}

func (r Rational) Subtract(r2 Rational) Rational {
	return newRationalFromC(C.av_sub_q(r.c, r2.c))
}

func (r Rational) Invert() Rational {
	return newRationalFromC(C.av_inv_q(r.c))
}

func (r Rational) Nearer(q1, q2 Rational) int {
	return int(C.av_nearer_q(r.c, q1.c, q2.c))
}

func (r Rational) NearestIndex(q []Rational) int {
	return int(C.av_find_nearest_q_idx(r.c, (*C.struct_AVRational)(unsafe.Pointer(&q[0]))))
}

func (r Rational) Num() int {
	return int(r.c.num)
}

func (r *Rational) SetNum(num int) {
	r.c.num = C.int(num)
}

func (r Rational) Den() int {
	return int(r.c.den)
}

func (r *Rational) SetDen(den int) {
	r.c.den = C.int(den)
}

func (r Rational) ToDouble() float64 {
	if r.Num() == 0 || r.Den() == 0 {
		return 0
	}
	return float64(r.Num()) / float64(r.Den())
}

func (r Rational) String() string {
	if r.Num() == 0 || r.Den() == 0 {
		return "0"
	}
	return strconv.Itoa(r.Num()) + "/" + strconv.Itoa(r.Den())
}
