package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
)

// char args[512];
// int ret = 0;
// const AVFilter *bufferSrc  = avfilter_get_by_name("buffer");
// const AVFilter *bufferOvr  = avfilter_get_by_name("buffer");
// const AVFilter *bufferSink = avfilter_get_by_name("buffersink");
// const AVFilter *ovrFilter  = avfilter_get_by_name("overlay");
// const AVFilter *colorFilter  = avfilter_get_by_name("colorchannelmixer");
// enum AVPixelFormat pix_fmts[] = { AV_PIX_FMT_YUV420P, AV_PIX_FMT_NONE };

// fFilterGraph = avfilter_graph_alloc();
// if (!fFilterGraph) {
// 	ret = AVERROR(ENOMEM);
// 	goto end;
// }

// /* buffer video source: the decoded frames from the decoder will be inserted here. */
// snprintf(args, sizeof(args),
// 		"video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
// 		decCtx->width, decCtx->height, decCtx->pix_fmt,
// 		fTimeBase.num, fTimeBase.den,
// 		decCtx->sample_aspect_ratio.num, decCtx->sample_aspect_ratio.den);
// ret = avfilter_graph_create_filter(&fBufSrc0Ctx, bufferSrc, "in0",
// 					args, NULL, fFilterGraph);
// if (ret < 0)
// 	goto end;

// /* buffer video overlay source: the overlayed frame from the file will be inserted here. */
// snprintf(args, sizeof(args),
// 		"video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
// 		ovrCtx->width, ovrCtx->height, ovrCtx->pix_fmt,
// 		fTimeBase.num, fTimeBase.den,
// 		ovrCtx->sample_aspect_ratio.num, ovrCtx->sample_aspect_ratio.den);
// ret = avfilter_graph_create_filter(&fBufSrc1Ctx, bufferOvr, "in1",
// 					args, NULL, fFilterGraph);
// if (ret < 0)
// 	goto end;

// /* color filter */
// snprintf(args, sizeof(args), "aa=%f", (float)fWatermarkOpacity / 10.0);
// ret = avfilter_graph_create_filter(&fColorFilterCtx, colorFilter, "colorFilter",
// 					args, NULL, fFilterGraph);
// if (ret < 0)
// 	goto end;

// /* overlay filter */
// switch (fWatermarkPos) {
// case 0:
// 	/* Top left */
// 	snprintf(args, sizeof(args), "x=%d:y=%d:repeatlast=1",
// 			fWatermarkOffset, fWatermarkOffset);
// 	break;
// case 1:
// 	/* Top right */
// 	snprintf(args, sizeof(args), "x=W-w-%d:y=%d:repeatlast=1",
// 			fWatermarkOffset, fWatermarkOffset);
// 	break;
// case 3:
// 	/* Bottom left */
// 	snprintf(args, sizeof(args), "x=%d:y=H-h-%d:repeatlast=1",
// 			fWatermarkOffset, fWatermarkOffset);
// 	break;
// case 4:
// 	/* Bottom right */
// 	snprintf(args, sizeof(args), "x=W-w-%d:y=H-h-%d:repeatlast=1",
// 			fWatermarkOffset, fWatermarkOffset);
// 	break;

// case 2:
// default:
// 	/* Center */
// 	snprintf(args, sizeof(args), "x=(W-w)/2:y=(H-h)/2:repeatlast=1");
// 	break;
// }
// ret = avfilter_graph_create_filter(&fOvrFilterCtx, ovrFilter, "overlay",
// 					args, NULL, fFilterGraph);
// if (ret < 0)
// 	goto end;

// /* buffer sink - destination of the final video */
// ret = avfilter_graph_create_filter(&fBufSinkCtx, bufferSink, "out",
// 					NULL, NULL, fFilterGraph);
// if (ret < 0)
// 	goto end;

// ret = av_opt_set_int_list(fBufSinkCtx, "pix_fmts", pix_fmts,
// 				AV_PIX_FMT_NONE, AV_OPT_SEARCH_CHILDREN);
// if (ret < 0)
// 	goto end;

// /*
// 	* Link all filters..
// 	*/
// avfilter_link(fBufSrc0Ctx, 0, fOvrFilterCtx, 0);
// avfilter_link(fBufSrc1Ctx, 0, fColorFilterCtx, 0);
// avfilter_link(fColorFilterCtx, 0, fOvrFilterCtx, 1);
// avfilter_link(fOvrFilterCtx, 0, fBufSinkCtx, 0);
// if ((ret = avfilter_graph_config(fFilterGraph, NULL)) < 0)
// 	goto end;

// end:

var (
	c              = astikit.NewCloser()
	src     string = os.Args[0]
	overlay string = os.Args[1]
	dst     string = os.Args[2]

	srcFormatContext     *astiav.FormatContext
	overlayFormatContext *astiav.FormatContext
	dstFormatContext     *astiav.FormatContext

	bufferSrcContext     *astiav.FilterContext
	bufferOverlayContext *astiav.FilterContext
	bufferSinkContext    *astiav.FilterContext

	decCodec        *astiav.Codec
	decCodecContext *astiav.CodecContext
	decFrame        *astiav.Frame

	filterDesc  = "[in]scale=300:100[scl];[in1][scl]overlay=25:25"
	filterFrame *astiav.Frame
	filterGraph *astiav.FilterGraph

	inputStream *astiav.Stream
	lastPts     int64
)

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelDebug)
	astiav.SetLogCallback(func(l astiav.LogLevel, fmt, msg, parent string) {
		log.Printf("ffmpeg log: %s (level: %s)\n", strings.TrimSpace(msg), l)
	})

	// Usage
	// if inputs[0] == "" {
	// 	log.Println("Usage: <binary path> -i <input path>")
	// 	return
	// }

	// We use an astikit.Closer to free all resources properly
	defer c.Close()

	// Open input file
	if err := openInputFile(); err != nil {
		log.Fatal(fmt.Errorf("main: opening input file failed: %w", err))
	}

	// Init filter
	if err := initFilter(); err != nil {
		log.Fatal(fmt.Errorf("main: initializing filter failed: %w", err))
	}

	// Alloc packet
	pkt := astiav.AllocPacket()
	c.Add(pkt.Free)

	// Loop through packets
	for {
		// Read frame
		if err := srcFormatContext.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
		}

		// Invalid stream
		if pkt.StreamIndex() != inputStream.Index() {
			continue
		}

		// Send packet
		if err := decCodecContext.SendPacket(pkt); err != nil {
			log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
		}

		// Loop
		for {
			// Receive frame
			if err := decCodecContext.ReceiveFrame(decFrame); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					break
				}
				log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
			}

			// Filter frame
			if err := filter(decFrame); err != nil {
				log.Fatal(fmt.Errorf("main: filtering frame failed: %w", err))
			}
		}
	}

	// Flush filter
	if err := filter(nil); err != nil {
		log.Fatal(fmt.Errorf("main: filtering frame failed: %w", err))
	}

	// Success
	log.Println("success")
}

func openInputFile() (err error) {
	// Alloc input format context
	if srcFormatContext = astiav.AllocFormatContext(); srcFormatContext == nil {
		err = errors.New("main: input format context is nil")
		return
	}
	c.Add(srcFormatContext.Free)

	// Open input
	if err = srcFormatContext.OpenInput(src, nil, nil); err != nil {
		err = fmt.Errorf("main: opening input failed: %w", err)
		return
	}
	c.Add(srcFormatContext.CloseInput)

	// Find stream info
	if err = srcFormatContext.FindStreamInfo(nil); err != nil {
		err = fmt.Errorf("main: finding stream info failed: %w", err)
		return
	}

	// Loop through streams
	for _, is := range srcFormatContext.Streams() {
		// Only process video
		if is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}
		inputStream = is
		lastPts = astiav.NoPtsValue

		// Find decoder
		if decCodec = astiav.FindDecoder(is.CodecParameters().CodecID()); decCodec == nil {
			err = errors.New("main: codec is nil")
			return
		}

		// Alloc codec context
		if decCodecContext = astiav.AllocCodecContext(decCodec); decCodecContext == nil {
			err = errors.New("main: codec context is nil")
			return
		}
		c.Add(decCodecContext.Free)

		// Update codec context
		if err = is.CodecParameters().ToCodecContext(decCodecContext); err != nil {
			err = fmt.Errorf("main: updating codec context failed: %w", err)
			return
		}

		// Open codec context
		if err = decCodecContext.Open(decCodec, nil); err != nil {
			err = fmt.Errorf("main: opening codec context failed: %w", err)
			return
		}

		// Alloc frame
		decFrame = astiav.AllocFrame()
		c.Add(decFrame.Free)

		break
	}
	return
}

func initFilter() (err error) {
	// Alloc graph
	if filterGraph = astiav.AllocFilterGraph(); filterGraph == nil {
		err = errors.New("main: graph is nil")
		return
	}
	c.Add(filterGraph.Free)

	// Alloc outputs
	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		err = errors.New("main: outputs is nil")
		return
	}
	c.Add(outputs.Free)

	// Alloc inputs
	inputs := astiav.AllocFilterInOut()
	if inputs == nil {
		err = errors.New("main: inputs is nil")
		return
	}
	c.Add(inputs.Free)

	// Create buffersrc
	buffersrc := astiav.FindFilterByName("buffer")
	if buffersrc == nil {
		err = errors.New("main: buffersrc is nil")
		return
	}

	// Create buffersink
	buffersink := astiav.FindFilterByName("buffersink")
	if buffersink == nil {
		err = errors.New("main: buffersink is nil")
		return
	}

	// Create filter contexts
	if bufferSrcContext, err = filterGraph.NewFilterContext(buffersrc, "in", astiav.FilterArgs{
		"pix_fmt":      strconv.Itoa(decCodecContext.PixelFormat().Int()),
		"pixel_aspect": decCodecContext.SampleAspectRatio().String(),
		"time_base":    inputStream.TimeBase().String(),
		"video_size":   strconv.Itoa(decCodecContext.Width()) + "x" + strconv.Itoa(decCodecContext.Height()),
		// "re":           "",
	}); err != nil {
		err = fmt.Errorf("main: creating buffersrc1 context failed: %w", err)
		return
	}
	if bufferOverlayContext, err = filterGraph.NewFilterContext(buffersrc, "in1", astiav.FilterArgs{
		// "pix_fmt":      strconv.Itoa(decCodecContext.PixelFormat().Int()),
		// "pixel_aspect": decCodecContext.SampleAspectRatio().String(),
		// "time_base":    inputStream.TimeBase().String(),
		// "video_size":   strconv.Itoa(decCodecContext.Width()) + "x" + strconv.Itoa(decCodecContext.Height()),
		// "re":           "",
	}); err != nil {
		err = fmt.Errorf("main: creating buffersrc2 context failed: %w", err)
		return
	}
	if bufferSinkContext, err = filterGraph.NewFilterContext(buffersink, "in", nil); err != nil {
		err = fmt.Errorf("main: creating buffersink context failed: %w", err)
		return
	}
	if bufferSinkContext, err = filterGraph.NewFilterContext(buffersink, "in1", nil); err != nil {
		err = fmt.Errorf("main: creating buffersink context failed: %w", err)
		return
	}

	// Update outputs
	outputs.SetName("in")
	outputs.SetFilterContext(bufferSrcContext)
	outputs.SetPadIdx(0)
	// outputs.SetNext(outputs)

	// outputs.SetName("in1")
	// outputs.SetFilterContext(bufferOverlayContext)
	// outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	// Update inputs
	inputs.SetName("out")
	inputs.SetFilterContext(bufferSinkContext)
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// Parse
	if err = filterGraph.Parse(filterDesc, inputs, outputs); err != nil {
		err = fmt.Errorf("main: parsing filter failed: %w", err)
		return
	}

	// Configure
	if err = filterGraph.Configure(); err != nil {
		err = fmt.Errorf("main: configuring filter failed: %w", err)
		return
	}

	// Alloc frame
	filterFrame = astiav.AllocFrame()
	c.Add(filterFrame.Free)
	return
}

func filter(f *astiav.Frame) (err error) {
	// Add frame
	if err = bufferSrcContext.BuffersrcAddFrame(f, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
		err = fmt.Errorf("main: adding frame failed: %w", err)
		return
	}

	// Loop
	for {
		// Unref frame
		filterFrame.Unref()

		// Get frame
		if err = bufferSinkContext.BuffersinkGetFrame(filterFrame, astiav.NewBuffersinkFlags()); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				err = nil
				break
			}
			err = fmt.Errorf("main: getting frame failed: %w", err)
			return
		}

		// Do something with filtered frame
		log.Printf("new filtered frame: %dx%d\n", filterFrame.Width(), filterFrame.Height())
	}
	return
}
