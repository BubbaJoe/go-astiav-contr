package main

import (
	"fmt"

	"github.com/asticode/go-astiav"
	"gocv.io/x/gocv"
)

// var (
// 	input  = flag.String("i", "", "the input path")
// 	output = flag.String("o", "", "the output path")
// )

// var (
// 	c                   = astikit.NewCloser()
// 	inputFormatContext  *astiav.FormatContext
// 	outputFormatContext *astiav.FormatContext
// 	streams             = make(map[int]*stream) // Indexed by input stream index
// )

// type stream struct {
// 	buffersinkContext *astiav.FilterContext
// 	buffersrcContext  *astiav.FilterContext
// 	decCodec          *astiav.Codec
// 	decCodecContext   *astiav.CodecContext
// 	decFrame          *astiav.Frame
// 	encCodec          *astiav.Codec
// 	encCodecContext   *astiav.CodecContext
// 	encPkt            *astiav.Packet
// 	filterFrame       *astiav.Frame
// 	filterGraph       *astiav.FilterGraph
// 	inputStream       *astiav.Stream
// 	outputStream      *astiav.Stream
// }

// func main() {
// 	// Handle ffmpeg logs
// 	astiav.SetLogLevel(astiav.LogLevelInfo)
// 	astiav.SetLogCallback(func(l astiav.LogLevel, fmt, msg, parent string) {
// 		log.Printf("ffmpeg log: %s (level: %d)\n", strings.TrimSpace(msg), l)
// 	})

// 	// Parse flags
// 	flag.Parse()

// 	// Usage
// 	if *input == "" || *output == "" {
// 		log.Println("Usage: <binary path> -i <input path> -o <output path>")
// 		return
// 	}

// 	// We use an astikit.Closer to free all resources properly
// 	defer c.Close()

// 	// Open input file
// 	if err := openInputFile(); err != nil {
// 		log.Fatal(fmt.Errorf("main: opening input file failed: %w", err))
// 	}

// 	// Open output file
// 	if err := openOutputFile(); err != nil {
// 		log.Fatal(fmt.Errorf("main: opening output file failed: %w", err))
// 	}

// 	// Init filters
// 	if err := initFilters(); err != nil {
// 		log.Fatal(fmt.Errorf("main: initializing filters failed: %w", err))
// 	}

// 	// Alloc packet
// 	pkt := astiav.AllocPacket()
// 	c.Add(pkt.Free)

// 	// Loop through packets
// 	for {
// 		// Read frame
// 		if err := inputFormatContext.ReadFrame(pkt); err != nil {
// 			if errors.Is(err, astiav.ErrEof) {
// 				break
// 			}
// 			log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
// 		}

// 		// Get stream
// 		s, ok := streams[pkt.StreamIndex()]
// 		if !ok {
// 			continue
// 		}

// 		// Update packet
// 		pkt.RescaleTs(s.inputStream.TimeBase(), s.decCodecContext.TimeBase())

// 		// Send packet
// 		if err := s.decCodecContext.SendPacket(pkt); err != nil {
// 			log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
// 		}

// 		// Loop
// 		for {
// 			// Receive frame
// 			if err := s.decCodecContext.ReceiveFrame(s.decFrame); err != nil {
// 				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
// 					break
// 				}
// 				log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
// 			}

// 			// Filter, encode and write frame
// 			if err := filterEncodeWriteFrame(s.decFrame, s); err != nil {
// 				log.Fatal(fmt.Errorf("main: filtering, encoding and writing frame failed: %w", err))
// 			}
// 		}
// 	}

// 	// Loop through streams
// 	for _, s := range streams {
// 		// Flush filter
// 		if err := filterEncodeWriteFrame(nil, s); err != nil {
// 			log.Fatal(fmt.Errorf("main: filtering, encoding and writing frame failed: %w", err))
// 		}

// 		// Flush encoder
// 		if err := encodeWriteFrame(nil, s); err != nil {
// 			log.Fatal(fmt.Errorf("main: encoding and writing frame failed: %w", err))
// 		}
// 	}

// 	// Write trailer
// 	if err := outputFormatContext.WriteTrailer(); err != nil {
// 		log.Fatal(fmt.Errorf("main: writing trailer failed: %w", err))
// 	}

// 	// Success
// 	log.Println("success")
// }

// func openInputFile() (err error) {
// 	// Alloc input format context
// 	if inputFormatContext = astiav.AllocFormatContext(); inputFormatContext == nil {
// 		err = errors.New("main: input format context is nil")
// 		return
// 	}
// 	c.Add(inputFormatContext.Free)

// 	// Open input
// 	if err = inputFormatContext.OpenInput(*input, nil, nil); err != nil {
// 		err = fmt.Errorf("main: opening input failed: %w", err)
// 		return
// 	}
// 	c.Add(inputFormatContext.CloseInput)

// 	// Find stream info
// 	if err = inputFormatContext.FindStreamInfo(nil); err != nil {
// 		err = fmt.Errorf("main: finding stream info failed: %w", err)
// 		return
// 	}

// 	// Loop through streams
// 	for _, is := range inputFormatContext.Streams() {
// 		// Only process audio or video
// 		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
// 			is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
// 			continue
// 		}

// 		// Create stream
// 		s := &stream{inputStream: is}

// 		// Find decoder
// 		if s.decCodec = astiav.FindDecoder(is.CodecParameters().CodecID()); s.decCodec == nil {
// 			err = errors.New("main: codec is nil")
// 			return
// 		}

// 		// Alloc codec context
// 		if s.decCodecContext = astiav.AllocCodecContext(s.decCodec); s.decCodecContext == nil {
// 			err = errors.New("main: codec context is nil")
// 			return
// 		}
// 		c.Add(s.decCodecContext.Free)

// 		// Update codec context
// 		if err = is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
// 			err = fmt.Errorf("main: updating codec context failed: %w", err)
// 			return
// 		}

// 		// Set framerate
// 		if is.CodecParameters().MediaType() == astiav.MediaTypeVideo {
// 			s.decCodecContext.SetFramerate(inputFormatContext.GuessFrameRate(is, nil))
// 		}

// 		// Open codec context
// 		if err = s.decCodecContext.Open(s.decCodec, nil); err != nil {
// 			err = fmt.Errorf("main: opening codec context failed: %w", err)
// 			return
// 		}

// 		// Alloc frame
// 		s.decFrame = astiav.AllocFrame()
// 		c.Add(s.decFrame.Free)

// 		// Store stream
// 		streams[is.Index()] = s
// 	}
// 	return
// }

// func openOutputFile() (err error) {
// 	// Alloc output format context
// 	if outputFormatContext, err = astiav.AllocOutputFormatContext(nil, "", *output); err != nil {
// 		err = fmt.Errorf("main: allocating output format context failed: %w", err)
// 		return
// 	} else if outputFormatContext == nil {
// 		err = errors.New("main: output format context is nil")
// 		return
// 	}
// 	c.Add(outputFormatContext.Free)

// 	// Loop through streams
// 	for _, is := range inputFormatContext.Streams() {
// 		// Get stream
// 		s, ok := streams[is.Index()]
// 		if !ok {
// 			continue
// 		}

// 		// Create output stream
// 		if s.outputStream = outputFormatContext.NewStream(nil); s.outputStream == nil {
// 			err = errors.New("main: output stream is nil")
// 			return
// 		}

// 		// Get codec id
// 		codecID := astiav.CodecIDH264
// 		if s.decCodecContext.MediaType() == astiav.MediaTypeAudio {
// 			codecID = astiav.CodecIDAac
// 		}

// 		// Find encoder
// 		if s.encCodec = astiav.FindEncoder(codecID); s.encCodec == nil {
// 			err = errors.New("main: codec is nil")
// 			return
// 		}

// 		// Alloc codec context
// 		if s.encCodecContext = astiav.AllocCodecContext(s.encCodec); s.encCodecContext == nil {
// 			err = errors.New("main: codec context is nil")
// 			return
// 		}
// 		c.Add(s.encCodecContext.Free)

// 		// Update codec context
// 		if s.decCodecContext.MediaType() == astiav.MediaTypeAudio {
// 			if v := s.encCodec.ChannelLayouts(); len(v) > 0 {
// 				s.encCodecContext.SetChannelLayout(v[0])
// 			} else {
// 				s.encCodecContext.SetChannelLayout(s.decCodecContext.ChannelLayout())
// 			}
// 			s.encCodecContext.SetChannels(s.decCodecContext.Channels())
// 			s.encCodecContext.SetSampleRate(s.decCodecContext.SampleRate())
// 			if v := s.encCodec.SampleFormats(); len(v) > 0 {
// 				s.encCodecContext.SetSampleFormat(v[0])
// 			} else {
// 				s.encCodecContext.SetSampleFormat(s.decCodecContext.SampleFormat())
// 			}
// 			s.encCodecContext.SetTimeBase(s.decCodecContext.TimeBase())
// 		} else {
// 			s.encCodecContext.SetHeight(s.decCodecContext.Height())
// 			if v := s.encCodec.PixelFormats(); len(v) > 0 {
// 				s.encCodecContext.SetPixelFormat(v[0])
// 			} else {
// 				s.encCodecContext.SetPixelFormat(s.decCodecContext.PixelFormat())
// 			}
// 			// fmt.Println("Previous pixel format: ", s.encCodecContext.PixelFormat())
// 			// fmt.Println("Previous pixel format: ", s.decCodecContext.PixelFormat())
// 			s.encCodecContext.SetSampleAspectRatio(s.decCodecContext.SampleAspectRatio())
// 			s.encCodecContext.SetTimeBase(s.decCodecContext.TimeBase())
// 			s.encCodecContext.SetWidth(s.decCodecContext.Width())
// 		}

// 		// Update flags
// 		if s.decCodecContext.Flags().Has(astiav.CodecContextFlagGlobalHeader) {
// 			s.encCodecContext.SetFlags(s.encCodecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader))
// 		}

// 		// Open codec context
// 		if err = s.encCodecContext.Open(s.encCodec, nil); err != nil {
// 			err = fmt.Errorf("main: opening codec context failed: %w", err)
// 			return
// 		}

// 		// Update codec parameters
// 		if err = s.outputStream.CodecParameters().FromCodecContext(s.encCodecContext); err != nil {
// 			err = fmt.Errorf("main: updating codec parameters failed: %w", err)
// 			return
// 		}

// 		// Update stream
// 		s.outputStream.SetTimeBase(s.encCodecContext.TimeBase())
// 	}

// 	// If this is a file, we need to use an io context
// 	if !outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
// 		// Create io context
// 		ioContext := astiav.NewIOContext()

// 		// Open io context
// 		if err = ioContext.Open(*output, astiav.NewIOContextFlags(astiav.IOContextFlagWrite)); err != nil {
// 			err = fmt.Errorf("main: opening io context failed: %w", err)
// 			return
// 		}
// 		c.AddWithError(ioContext.Closep)

// 		// Update output format context
// 		outputFormatContext.SetPb(ioContext)
// 	}

// 	// Write header
// 	if err = outputFormatContext.WriteHeader(nil); err != nil {
// 		err = fmt.Errorf("main: writing header failed: %w", err)
// 		return
// 	}
// 	return
// }

// func initFilters() (err error) {
// 	// Loop through output streams
// 	for _, s := range streams {
// 		// Alloc graph
// 		if s.filterGraph = astiav.AllocFilterGraph(); s.filterGraph == nil {
// 			err = errors.New("main: graph is nil")
// 			return
// 		}
// 		c.Add(s.filterGraph.Free)

// 		// Alloc outputs
// 		outputs := astiav.AllocFilterInOut()
// 		if outputs == nil {
// 			err = errors.New("main: outputs is nil")
// 			return
// 		}
// 		c.Add(outputs.Free)

// 		// Alloc inputs
// 		inputs := astiav.AllocFilterInOut()
// 		if inputs == nil {
// 			err = errors.New("main: inputs is nil")
// 			return
// 		}
// 		c.Add(inputs.Free)

// 		// Switch on media type
// 		var args astiav.FilterArgs
// 		var buffersrc, buffersink *astiav.Filter
// 		var content string
// 		switch s.decCodecContext.MediaType() {
// 		case astiav.MediaTypeAudio:
// 			args = astiav.FilterArgs{
// 				"channel_layout": s.decCodecContext.ChannelLayout().String(),
// 				"sample_fmt":     s.decCodecContext.SampleFormat().Name(),
// 				"sample_rate":    strconv.Itoa(s.decCodecContext.SampleRate()),
// 				"time_base":      s.decCodecContext.TimeBase().String(),
// 			}
// 			buffersrc = astiav.FindFilterByName("abuffer")
// 			buffersink = astiav.FindFilterByName("abuffersink")
// 			content = fmt.Sprintf("aformat=sample_fmts=%s:channel_layouts=%s", s.encCodecContext.SampleFormat().Name(), s.encCodecContext.ChannelLayout().String())
// 		default:
// 			args = astiav.FilterArgs{
// 				"pix_fmt":      strconv.Itoa(int(s.decCodecContext.PixelFormat())),
// 				"pixel_aspect": s.decCodecContext.SampleAspectRatio().String(),
// 				"time_base":    s.decCodecContext.TimeBase().String(),
// 				"video_size":   strconv.Itoa(s.decCodecContext.Width()) + "x" + strconv.Itoa(s.decCodecContext.Height()),
// 			}
// 			buffersrc = astiav.FindFilterByName("buffer")
// 			buffersink = astiav.FindFilterByName("buffersink")
// 			content = fmt.Sprintf("format=pix_fmts=%s", s.encCodecContext.PixelFormat().Name())
// 		}

// 		// Check filters
// 		if buffersrc == nil {
// 			err = errors.New("main: buffersrc is nil")
// 			return
// 		}
// 		if buffersink == nil {
// 			err = errors.New("main: buffersink is nil")
// 			return
// 		}

// 		// Create filter contexts
// 		if s.buffersrcContext, err = s.filterGraph.NewFilterContext(buffersrc, "in", args); err != nil {
// 			err = fmt.Errorf("main: creating buffersrc context failed: %w", err)
// 			return
// 		}
// 		if s.buffersinkContext, err = s.filterGraph.NewFilterContext(buffersink, "in", nil); err != nil {
// 			err = fmt.Errorf("main: creating buffersink context failed: %w", err)
// 			return
// 		}

// 		// Update outputs
// 		outputs.SetName("in")
// 		outputs.SetFilterContext(s.buffersrcContext)
// 		outputs.SetPadIdx(0)
// 		outputs.SetNext(nil)

// 		// Update inputs
// 		inputs.SetName("out")
// 		inputs.SetFilterContext(s.buffersinkContext)
// 		inputs.SetPadIdx(0)
// 		inputs.SetNext(nil)

// 		// Parse
// 		if err = s.filterGraph.Parse(content, inputs, outputs); err != nil {
// 			err = fmt.Errorf("main: parsing filter failed: %w", err)
// 			return
// 		}

// 		// Configure
// 		if err = s.filterGraph.Configure(); err != nil {
// 			err = fmt.Errorf("main: configuring filter failed: %w", err)
// 			return
// 		}

// 		// Alloc frame
// 		s.filterFrame = astiav.AllocFrame()
// 		c.Add(s.filterFrame.Free)

// 		// Alloc packet
// 		s.encPkt = astiav.AllocPacket()
// 		c.Add(s.encPkt.Free)
// 	}
// 	return
// }

// func filterEncodeWriteFrame(f *astiav.Frame, s *stream) (err error) {
// 	// Add frame
// 	if err = s.buffersrcContext.BuffersrcAddFrame(f, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
// 		err = fmt.Errorf("main: adding frame failed: %w", err)
// 		return
// 	}

// 	// Loop
// 	for {
// 		// Unref frame
// 		s.filterFrame.Unref()

// 		// Get frame
// 		if err = s.buffersinkContext.BuffersinkGetFrame(s.filterFrame, astiav.NewBuffersinkFlags()); err != nil {
// 			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
// 				err = nil
// 				break
// 			}
// 			err = fmt.Errorf("main: getting frame failed: %w", err)
// 			return
// 		}

// 		// Reset picture type
// 		s.filterFrame.SetPictureType(astiav.PictureTypeNone)
// 		// Print previous pixel format
// 		// fmt.Println("Previous pixel format: ", s.filterFrame.PixelFormat())
// 		s.filterFrame.SetPixelFormat(astiav.PixelFormatYuv420P)

// 		// Encode and write frame
// 		if err = encodeWriteFrame(s.filterFrame, s); err != nil {
// 			err = fmt.Errorf("main: encoding and writing frame failed: %w", err)
// 			return
// 		}
// 	}
// 	return
// }

// func encodeWriteFrame(f *astiav.Frame, s *stream) (err error) {
// 	// Unref packet
// 	s.encPkt.Unref()

// 	// Send frame
// 	if err = s.encCodecContext.SendFrame(f); err != nil {
// 		err = fmt.Errorf("main: sending frame failed: %w", err)
// 		return
// 	}

// 	// Loop
// 	for {
// 		// Receive packet
// 		if err = s.encCodecContext.ReceivePacket(s.encPkt); err != nil {
// 			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
// 				err = nil
// 				break
// 			}
// 			err = fmt.Errorf("main: receiving packet failed: %w", err)
// 			return
// 		}

// 		// Update pkt
// 		s.encPkt.SetStreamIndex(s.outputStream.Index())
// 		s.encPkt.RescaleTs(s.encCodecContext.TimeBase(), s.outputStream.TimeBase())

// 		// Write frame
// 		if err = outputFormatContext.WriteInterleavedFrame(s.encPkt); err != nil {
// 			err = fmt.Errorf("main: writing frame failed: %w", err)
// 			return
// 		}
// 	}
// 	return
// }

// func writeFrame(cc *astiav.CodecContext, fc *astiav.FormatContext, frame *astiav.Frame) error {
// 	pkt := astiav.AllocPacket()

// 	err := cc.SendFrame(frame)
// 	if err != nil {
// 		return err
// 	}
// 	cc.ReceivePacket(pkt)
// 	if err != nil {
// 		return err
// 	}
// 	fc.WriteInterleavedFrame(pkt)
// 	defer pkt.Unref()
// 	return nil
// }

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

	s, err := cc.GetBufferSize(width, height, 1)
	if err != nil {
		panic(err)
	}

	err = astiav.ImageFillFrameArrays(
		frame, s, cc.PixelFormat(),
		width, height, 1,
	)
	if err != nil {
		panic(err)
	}
	frame.SetWidth(width)
	frame.SetHeight(height)
	frame.SetPixelFormat(cc.PixelFormat())

	_, err = frame.AllocImage(32)
	if err != nil {
		panic(err)
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

func streamVideo(width, height, fps, cameraId int) {
	streamKey := "rfBd56ti2SMtYvSgD5xAV0YU99zampta7Z7S575KLkIZ9PYk"
	output := fmt.Sprintf("rtmp://localhost:1935/live/%s", streamKey)
	// output := "./test.flv"
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

	if fmtCtx.Pb() == nil {
		fmtCtx.SetPb(astiav.NewIOContext())
	}
	// defer fmtCtx.Free()

	ioCtx := astiav.NewIOContext()
	defer ioCtx.Closep()
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

	frame := allocFrameBuffer(outCodecCtx, width, height)

	err = fmtCtx.WriteHeader(nil)
	if err != nil {
		panic(err)
	}
	numReads := 0
	retry := 0
	for {
		if cam.Read(&mat) {
			numReads += 1
			// fmt.Println("Read frame", mat.Step())
			matPtrs, err := mat.DataPtrUint8()
			if err != nil {
				panic(err)
			}
			// tmp := make([]byte, len(matPtrs))
			// copy(tmp, matPtrs)
			newMatBytes := [8][]byte{}
			subSize := len(matPtrs) / 8
			for i := 0; i < 8; i++ {
				newMatBytes[i] = make([]byte, subSize)
				for j := 0; j < subSize; j++ {
					newMatBytes[i] = append(newMatBytes[i], matPtrs[i*subSize+j])
				}
			}
			stride := [8]int{}
			for i := range stride {
				stride[i] = mat.Step()
			}

			srcFrame := astiav.AllocFrame()
			srcFrame.SetPixelFormat(astiav.PixelFormatBgr24)
			srcFrame.SetWidth(width)
			srcFrame.SetHeight(height)
			srcFrame.SetData(newMatBytes)
			srcFrame.SetLinesize(stride)
			srcFrame.AllocImage(1)
			err = SwsCtx.ScaleFrames(srcFrame, frame)

			// err = SwsCtx.ScaleDstFrame(
			// 	tmp, stride,
			// 	0, height, frame,
			// )
			if err != nil {
				fmt.Println("Closing, error: ", err)
				break
			}
			// window.IMShow(mat)
			// window.WaitKey(1)
			frame.SetPts(frame.Pts() + int64(outCodecCtx.
				TimeBase().Rescale(1, outStream.TimeBase())))
			writeFrame(outCodecCtx, fmtCtx, frame)

		} else {
			if retry > 100 {
				fmt.Println("Closing, numReads: ", numReads)
				break
			} else {
				retry++
			}
		}
	}

	fmtCtx.WriteTrailer()
	fmtCtx.Pb().Closep()
	// fmtCtx.Free()
}

func main() {
	streamVideo(1920, 1080, 30, 0)
}
