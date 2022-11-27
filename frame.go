package astiav

/*
#cgo pkg-config: libavutil
#include <libavutil/channel_layout.h>
#include <libavutil/frame.h>
#include <libavutil/imgutils.h>
#include <libavutil/samplefmt.h>
*/
import "C"
import (
	"unsafe"
)

const NumDataPointers = uint(C.AV_NUM_DATA_POINTERS)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/frame.h#L317
type Frame struct {
	c *C.struct_AVFrame
}

func newFrameFromC(c *C.struct_AVFrame) *Frame {
	if c == nil {
		return nil
	}
	return &Frame{c: c}
}

func AllocFrame() *Frame {
	return newFrameFromC(C.av_frame_alloc())
}

func (f *Frame) AllocBuffer(align int) error {
	return newError(C.av_frame_get_buffer(f.c, C.int(align)))
}

func (f *Frame) AllocImage(align int) (int, error) {
	n := C.av_image_alloc(
		&f.c.data[0], &f.c.linesize[0],
		f.c.width, f.c.height,
		(C.enum_AVPixelFormat)(f.c.format),
		C.int(align))
	if n < 0 {
		return 0, newError(n)
	}
	return int(n), nil
}

func (f *Frame) AllocSamples(align int) error {
	return newError(C.av_samples_alloc(&f.c.data[0], &f.c.linesize[0], C.av_get_channel_layout_nb_channels(f.c.channel_layout), f.c.nb_samples, (C.enum_AVSampleFormat)(f.c.format), C.int(align)))
}

func (f *Frame) ChannelLayout() ChannelLayout {
	return ChannelLayout(f.c.channel_layout)
}

func (f *Frame) SetChannelLayout(l ChannelLayout) {
	f.c.channel_layout = C.uint64_t(l)
}

func (f *Frame) Data() [NumDataPointers][]byte {
	b := [8][]byte{}
	for i, size := range f.getPlainSizes() {
		if size == 0 {
			b[i] = []byte{}
			continue
		}
		if f.c.data[i] == nil {
			continue
		}
		b[i] = C.GoBytes(unsafe.Pointer(f.c.data[i]), C.int(size))
	}
	for i := 0; i < int(NumDataPointers); i++ {
		size := f.c.linesize[i]
		if f.c.height > 0 {
			size = size * f.c.height
		} else if f.c.channels > 0 {
			size = size * f.c.channels
		}
		b[i] = C.GoBytes(unsafe.Pointer(f.c.data[0]), size)
	}
	return b
}

func (f *Frame) SetData(d [NumDataPointers][]byte) {
	panic("not implemented")
	// for i := 0; i < f.NbSamples(); i++ {
	// 	f.c.data[i] = (*C.uint8_t)(unsafe.Pointer(&d[i]))
	// }
}

func (f *Frame) DataPtr() [NumDataPointers]*byte {
	b := [NumDataPointers]*byte{}
	fData := f.Data()
	for i := 0; i < int(NumDataPointers); i++ {
		b[i] = &fData[i][0]
	}
	return b
}

func (f *Frame) DataFull() []byte {
	totalSize := 0
	sizes := f.getPlainSizes()
	for _, s := range sizes {
		totalSize += s
	}
	fullData := make([]byte, totalSize)
	currentStart := 0
	for _, byteArr := range f.Data() {
		if len(byteArr) == 0 {
			continue
		}
		currentStart += copy(fullData[currentStart:], byteArr)
	}
	return fullData

	// var fullSize int
	// for i := 0; i < int(NumDataPointers); i++ {
	// 	size := int(f.c.linesize[i])
	// 	if size == 0 {
	// 		continue
	// 	}

	// 	ls, err := ImageGetLinesize(f.PixelFormat(),
	// 		f.Width(), i)
	// 	if ls != size || err != nil {
	// 		fmt.Printf("%d != %d: %d\n", ls, size, i)
	// 		panic(err)
	// 	}
	// 	if f.c.height > 0 {
	// 		size *= int(f.c.height)
	// 	} else if f.c.channels > 0 {
	// 		size *= int(f.c.channels)
	// 	}
	// 	fullSize += size
	// }
	// return C.GoBytes(unsafe.Pointer(&f.c.data[0]), C.int(fullSize))
}

func (f *Frame) SetDataFull(b []byte) {
	currentPos := 0
	for i, size := range f.getPlainSizes() {
		sl := C.size_t(size)
		cb := C.av_malloc(sl)
		C.memcpy(cb, C.CBytes(b[currentPos:currentPos+size]), sl)
		f.c.data[i] = (*C.uint8_t)(unsafe.Pointer(cb))
		currentPos += size
	}

}

func (f *Frame) getPlainSize(i int) int {
	if i >= len(f.Linesize()) {
		return 0
	}
	size := f.Linesize()[i]
	if f.c.height > 0 {
		size *= int(f.c.height)
	} else if f.c.channels > 0 {
		size *= int(f.c.channels)
	}
	return size
}

func (f *Frame) getPlainSizes() []int {
	ls := f.Linesize()
	sizes := make([]int, len(ls))
	for i := 0; i < len(sizes); i++ {
		sizes[i] = f.getPlainSize(i)
	}

	return sizes
}

func (f *Frame) Linesize() [NumDataPointers]int {
	lsize := [NumDataPointers]int{}
	for i := 0; i < int(4); i++ {
		lsize[i] = int(f.c.linesize[i])
	}
	return lsize
}

func (f *Frame) Channels() int {
	return int(f.c.channels)
}

func (f *Frame) SetLinesize(l [NumDataPointers]int) {
	for i := 0; i < int(NumDataPointers); i++ {
		f.c.linesize[i] = C.int(l[i])
	}
}

func (f *Frame) Height() int {
	return int(f.c.height)
}

func (f *Frame) SetHeight(h int) {
	f.c.height = C.int(h)
}

func (f *Frame) KeyFrame() bool {
	return int(f.c.key_frame) > 0
}

func (f *Frame) SetKeyFrame(k bool) {
	i := 0
	if k {
		i = 1
	}
	f.c.key_frame = C.int(i)
}

func (f *Frame) NbSamples() int {
	return int(f.c.nb_samples)
}

func (f *Frame) SetNbSamples(n int) {
	f.c.nb_samples = C.int(n)
}

func (f *Frame) PictureType() PictureType {
	return PictureType(f.c.pict_type)
}

func (f *Frame) SetPictureType(t PictureType) {
	f.c.pict_type = C.enum_AVPictureType(t)
}

func (f *Frame) ColorRange() ColorRange {
	return ColorRange(f.c.color_range)
}

func (f *Frame) SetColorRange(cr ColorRange) {
	f.c.color_range = uint32(cr)
}

func (f *Frame) PixelFormat() PixelFormat {
	return PixelFormat(f.c.format)
}

func (f *Frame) SetPixelFormat(pf PixelFormat) {
	f.c.format = C.int(pf)
}

func (f *Frame) PktDts() int64 {
	return int64(f.c.pkt_dts)
}

func (f *Frame) Pts() int64 {
	return int64(f.c.pts)
}

func (f *Frame) SetPktDts(i int64) {
	f.c.pkt_dts = C.int64_t(i)
}

func (f *Frame) SetPts(i int64) {
	f.c.pts = C.int64_t(i)
}

func (f *Frame) SampleFormat() SampleFormat {
	return SampleFormat(f.c.format)
}

func (f *Frame) SetSampleFormat(sf SampleFormat) {
	f.c.format = C.int(sf)
}

func (f *Frame) SampleRate() int {
	return int(f.c.sample_rate)
}

func (f *Frame) SetSampleRate(r int) {
	f.c.sample_rate = C.int(r)
}

func (f *Frame) NewSideData(t FrameSideDataType, size int) *FrameSideData {
	return newFrameSideDataFromC(C.av_frame_new_side_data(f.c,
		(C.enum_AVFrameSideDataType)(t), C.int(size)))
}

func (f *Frame) SideData(t FrameSideDataType) *FrameSideData {
	return newFrameSideDataFromC(C.av_frame_get_side_data(f.c,
		(C.enum_AVFrameSideDataType)(t)))
}

func (f *Frame) Width() int {
	return int(f.c.width)
}

func (f *Frame) SetWidth(w int) {
	f.c.width = C.int(w)
}

func (f *Frame) Free() {
	C.av_frame_free(&f.c)
}

func (f *Frame) Ref(src *Frame) error {
	return newError(C.av_frame_ref(f.c, src.c))
}

func (f *Frame) Clone() *Frame {
	return newFrameFromC(C.av_frame_clone(f.c))
}

func (f *Frame) Unref() {
	C.av_frame_unref(f.c)
}

func (f *Frame) MoveRef(src *Frame) {
	C.av_frame_move_ref(f.c, src.c)
}
