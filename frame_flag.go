package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/frame.h>
import "C"

type FrameFlag int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L1519
const (
	FrameFlagCorrupt = FormatEventFlag(C.AV_FRAME_FLAG_CORRUPT)
	FrameFlagDescard = FormatEventFlag(C.AV_FRAME_FLAG_DISCARD)
)
