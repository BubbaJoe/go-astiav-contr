package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/asticode/go-astiav"
	"gocv.io/x/gocv"
)

func setCodecParams(fc *astiav.FormatContext, cc *astiav.CodecContext, width, height, fps int) error {
	cc.SetCodecTag(0)
	cc.SetCodecID(astiav.CodecIDH264)
	cc.SetWidth(width)
	cc.SetHeight(height)
	cc.SetTimeBase(astiav.NewRational(1, 25))
	cc.SetFramerate(astiav.NewRational(25, 1))
	cc.SetGopSize(10)
	cc.SetPixelFormat(astiav.PixelFormatYuv420P)
	dstFps := astiav.NewRational(fps, 1)
	cc.SetFramerate(dstFps)
	cc.SetTimeBase(dstFps.Invert())
	if fc.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
		cc.SetFlags(cc.Flags().Add(astiav.CodecContextFlagGlobalHeader))
	}
	return nil
}

func initCodecStream(stream *astiav.Stream, cc *astiav.CodecContext, codec *astiav.Codec) {
	err := stream.CodecParameters().FromCodecContext(cc)
	if err != nil {
		panic(err)
	}
	codecOpts := astiav.NewDictionary()
	codecOpts.Set("preset", "ultrafast", 0)
	codecOpts.Set("tune", "zerolatency", 0)
	codecOpts.Set("profile", "high", 0)

	err = cc.Open(codec, codecOpts)
	if err != nil {
		panic(err)
	}
}

// SwsContext *initialize_sample_scaler(AVCodecContext *codec_ctx, double width, double height)
func initSampleScalear(codecCtx *astiav.CodecContext, width, height int) *astiav.SwsContext {
	swsCtx := astiav.NewSwsContext(
		width, height, astiav.PixelFormatBgr24,
		width, height, codecCtx.PixelFormat(),
		astiav.SWS_BICUBIC,
		// astiav.SWS_BILINEAR,
		nil, nil, nil,
	)
	if swsCtx == nil {
		panic("swsCtx is nil")
	}
	return swsCtx
}

func allocFrameBuffer(cc *astiav.CodecContext, width, height int) *astiav.Frame {
	frame := astiav.AllocFrame()
	frame.SetHeight(height)
	frame.SetWidth(width)
	frame.SetPixelFormat(cc.PixelFormat())

	// _, err := cc.GetBufferSize(
	// 	width, height, 3)
	// _, err := astiav.ImageGetBufferSize(
	// 	cc.PixelFormat(), width, height, 3)
	// if err != nil {
	// 	panic(err)
	// }

	// err = astiav.ImageFillFrameArrays(
	// 	frame, make([]byte, s),
	// 	width, height, 1,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	n, err := frame.AllocImage(4)
	if err != nil {
		panic(err)
	}
	if n <= 0 {
		panic("whoops")
	}
	return frame
}

func writeFrame(cc *astiav.CodecContext, fc *astiav.FormatContext, frame *astiav.Frame) {
	pkt := astiav.AllocPacket()
	defer pkt.Unref()

	err := cc.SendFrame(frame)
	if err != nil {
		panic(err)
	}
	cc.ReceivePacket(pkt)
	if err != nil {
		panic(err)
	}
	fc.WriteInterleavedFrame(pkt)
}

func matToFrame(mat *gocv.Mat) *astiav.Frame {
	frame := astiav.AllocFrame()
	frame.SetHeight(mat.Rows())
	frame.SetWidth(mat.Cols())
	frame.SetPixelFormat(astiav.PixelFormatBgr24)
	_, err := frame.AllocImage(4)
	if err != nil {
		panic(err)
	}
	frame.SetDataFull(mat.ToBytes())
	// frame.SetLinesize([8]int{
	// 	1920, 960, 960,
	// })
	// frame.SetDataFull(mat.ToBytes())

	// err = astiav.ImageFillFrameArrays(
	// 	frame, mat.ToBytes(),
	// 	mat.Cols(), mat.Rows(), 1,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	return frame
}

func frameToMat(frame *astiav.Frame) *gocv.Mat {
	if frame.PixelFormat() != astiav.PixelFormatBgr24 {
		panic(frame.PixelFormat())
	}
	mat, err := gocv.NewMatFromBytes(
		frame.Height(), frame.Width(),
		gocv.MatTypeCV8UC3, frame.DataFull())
	if err != nil {
		panic(err)
	}
	// d := frame.Data()
	// fmt.Println(d[0])

	return &mat
}

func streamVideo(width, height, fps, cameraId int) {
	output := fmt.Sprintf("rtmp://localhost:1935/live/%s",
		"rfBd56ti2SMtYvSgD5xAV0YU99zampta7Z7S575KLkIZ9PYk")

	cam, err := gocv.OpenVideoCapture(cameraId)
	if err != nil {
		panic(err)
	}
	// window := gocv.NewWindow("Hello")
	mat := gocv.NewMatWithSize(width, height, gocv.MatTypeCV8UC3)
	defer mat.Close()

	fmtCtx, err := astiav.AllocOutputFormatContext(nil, "flv", "")
	if err != nil {
		panic(err)
	}
	// defer fmtCtx.Free()

	ioCtx := astiav.NewIOContext()
	fmtCtx.SetPb(ioCtx)
	defer ioCtx.Close()
	err = ioCtx.Open(output, astiav.NewIOContextFlags(astiav.IOContextFlagReadWrite))
	// err = ioCtx.Open2(fmtCtx.Pb(), output,
	// 	astiav.IOContextFlags(astiav.IOContextFlagReadWrite), nil)
	if err != nil {
		fmt.Println("open error", err, output)
		panic(err)
	}
	fmtCtx.SetPb(ioCtx)

	outCodec := astiav.FindEncoder(astiav.CodecIDH264)
	outStream := fmtCtx.NewStream(outCodec)
	outCodecCtx := astiav.AllocCodecContext(outCodec)
	outCodecCtx.SetPixelFormat(astiav.PixelFormatYuv420P)
	defer outCodecCtx.Free()

	setCodecParams(fmtCtx, outCodecCtx, width, height, fps)
	initCodecStream(outStream, outCodecCtx, outCodec)

	outStream.CodecParameters().SetExtradata(outCodecCtx.Extradata())
	outStream.CodecParameters().SetExtradataSize(outCodecCtx.ExtradataSize())

	fmtCtx.DumpFormat(0, output, true)

	SwsCtx := initSampleScalear(outCodecCtx, width, height)
	defer SwsCtx.Free()

	err = fmtCtx.WriteHeader(nil)
	if err != nil {
		panic(err)
	}
	numReads := 0
	retry := 0

	frame := allocFrameBuffer(
		outCodecCtx, width, height)
	frame.SetPixelFormat(astiav.PixelFormatYuv420P)

	defer frame.Free()

	for {
		if cam.Read(&mat) {
			numReads += 1
			// matPtrs, err := mat.DataPtrUint8()
			if err != nil {
				panic(err)
			}

			// newMatBytes := [8][]byte{}
			// subSize := len(matPtrs) / 8
			// for i := 0; i < 8; i++ {
			// 	newMatBytes[i] = make([]byte, subSize)
			// 	for j := 0; j < subSize; j++ {
			// 		newMatBytes[i] = append(newMatBytes[i], matPtrs[i*subSize+j])
			// 	}
			// }
			// stride := [8]int{}
			// for i := 0; i < mat.Channels(); i++ {
			// 	stride[i] = mat.Channels()
			// }

			srcFrame := matToFrame(&mat)
			// srcFrame.SetWidth(width)
			// srcFrame.SetHeight(height)
			// srcFrame.SetLinesize(stride)

			err = SwsCtx.ScaleFrames(srcFrame, frame)
			// err = SwsCtx.ScaleDstFrame(
			// 	make([]byte, 1), stride[:],
			// 	0, height, frame,
			// )
			if err != nil {
				fmt.Println("Closing, error: ", err)
				break
			}
			// fmt.Println(len(mat.ToBytes()), mat.Size(), mat.Channels(), mat.Type())
			// matx := *frameToMat(matToFrame(&mat))
			// fmt.Println(len(matx.ToBytes()), matx.Size(), matx.Channels(), matx.Type())
			// fmt.Println(mat.ToBytes()[:50], mat.ToBytes()[len(matx.ToBytes())-50:])
			// fmt.Println(matx.ToBytes()[:50], matx.ToBytes()[len(matx.ToBytes())-50:])
			// for i:=0;i<len(matx.ToBytes()); i++ {
			// 	if
			// }
			// fmt.Println(len(frame.DataFull()), frame.Linesize(), frame.Channels())
			// fmt.Println(len(srcFrame.DataFull()), srcFrame.Linesize(), frame.Channels())
			// yuv, _ := mat.ToImageYUV()
			// fmt.Println(
			// 	len(yuv.Y)+len(yuv.Cb)+len(yuv.Cr), "|",
			// 	len(yuv.Y), len(yuv.Cb), len(yuv.Cr),
			// 	yuv.YStride, yuv.CStride, yuv.CStride,
			// )

			frame.SetPts(frame.Pts() + int64(outCodecCtx.
				TimeBase().Rescale(1, outStream.TimeBase())))
			writeFrame(outCodecCtx, fmtCtx, frame)
			// frame.Free()
		} else {
			if retry > 100 {
				fmt.Println("Closing, numReads: ", numReads)
				break
			} else {
				time.Sleep(1)
				fmt.Println("Error getting picture: ", numReads)
				retry++
			}
		}
	}

	fmtCtx.WriteTrailer()
	fmtCtx.Pb().Closep()
}

func main() {
	astiav.SetLogLevel(astiav.LogLevelInfo)
	astiav.SetLogCallback(func(l astiav.LogLevel, format, msg, parent string) {
		if strings.TrimSpace(msg) == "" {
			return
		}
		fmt.Printf("FFMPEG: %s\n", msg)
	})

	streamVideo(1920, 1080, 30, 0)
}
