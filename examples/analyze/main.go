package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	input = flag.String("i", "", "the input path")
)

type stream struct {
	decCodec        *astiav.Codec
	decCodecContext *astiav.CodecContext
	inputStream     *astiav.Stream
}

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelInfo)
	astiav.SetLogCallback(func(l astiav.LogLevel, fmt, msg, parent string) {
		log.Printf("ffmpeg log: %s (level: %s)\n", strings.TrimSpace(msg), l)
	})

	// Parse flags
	flag.Parse()

	// Usage
	if *input == "" {
		log.Println("Usage: <binary path> -i <input path>")
		return
	}

	// Alloc packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	// Alloc frame
	f := astiav.AllocFrame()
	defer f.Free()

	// Alloc input format context
	inputFormatContext := astiav.AllocFormatContext()
	if inputFormatContext == nil {
		log.Fatal(errors.New("main: input format context is nil"))
	}
	defer inputFormatContext.Free()

	// Open input
	if *input != "-" {
		if err := inputFormatContext.OpenInput(*input, nil, nil); err != nil {
			log.Fatal(fmt.Errorf("main: opening input failed: %w", err))
		}
		defer inputFormatContext.CloseInput()
	} else {
		file, err := os.OpenFile("testdata/video.mp4", os.O_RDWR, 777)
		if err != nil {
			log.Fatal(fmt.Errorf("main: opening input failed: %w", err))
		}
		ioCtx := astiav.AllocIOContextReadSeek(
			file, file,
		)
		inputFormatContext.SetPb(ioCtx)
		inputFormatContext.OpenInput("", nil, nil)
	}

	// Find stream info
	dict := &astiav.Dictionary{}
	dict.Set("testing", "111", astiav.NewDictionaryFlags(
		astiav.DictionaryFlagAppend)) // test value
	if err := inputFormatContext.FindStreamInfo(dict); err != nil {
		log.Fatal(fmt.Errorf("main: finding stream info failed: %w", err))
	} else if dict.Len() > 0 {
		buf := make([]byte, 1024)
		dict.Unpack(buf)
		log.Printf("Dictionay Data: %s\n", string(buf))
	}

	// Loop through streams
	streams := make(map[int]*stream) // Indexed by input stream index
	for _, is := range inputFormatContext.Streams() {
		// Only process audio or video
		log.Printf("found media type (%s) for stream%d", is.CodecParameters().MediaType().String(), is.Index())
		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
			is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Create stream
		s := &stream{inputStream: is}

		// Find decoder
		if s.decCodec = astiav.FindDecoder(is.CodecParameters().CodecID()); s.decCodec == nil {
			log.Fatal(errors.New("main: codec is nil"))
		} else {
			log.Printf("main: found decoder for %s: %s\n", is.CodecParameters().CodecID().Name(), s.decCodec.Name())
		}

		// Find decoder
		if x := astiav.FindEncoder(is.CodecParameters().CodecID()); x == nil {
			log.Fatal(errors.New("main: codec is nil"))
		} else {
			log.Printf("main: found encoder for %s: %s\n", is.CodecParameters().CodecID().Name(), x.Name())
		}

		// Find pixel format
		if is.CodecParameters().MediaType() == astiav.MediaTypeVideo {
			log.Printf("pixel format: %s\n", is.CodecParameters().PixelFormat().Name())
			log.Printf("color range id: %d\n", is.CodecParameters().ColorRange())
			log.Printf("color space id: %d\n", is.CodecParameters().ColorSpace())
			log.Printf("color primaries id: %d\n", is.CodecParameters().ColorPrimaries())
		}

		// astiav.FindFilterByName("scale")

		// Alloc codec context
		if s.decCodecContext = astiav.AllocCodecContext(s.decCodec); s.decCodecContext == nil {
			log.Fatal(errors.New("codec context is nil"))
		}
		defer s.decCodecContext.Free()

		// Update codec context
		if err := is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
			log.Fatal(fmt.Errorf("updating codec context failed: %w", err))
		}

		// Open codec context
		if err := s.decCodecContext.Open(s.decCodec, nil); err != nil {
			log.Fatal(fmt.Errorf("opening codec context failed: %w", err))
		}

		// Add stream
		streams[is.Index()] = s
	}

	return

	// // Loop through packets
	// for {
	// 	// Read frame
	// 	if err := inputFormatContext.ReadFrame(pkt); err != nil {
	// 		if errors.Is(err, astiav.ErrEof) {
	// 			break
	// 		}
	// 		log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
	// 	}

	// 	// Get stream
	// 	s, ok := streams[pkt.StreamIndex()]
	// 	if !ok {
	// 		continue
	// 	}

	// 	// Send packet
	// 	if err := s.decCodecContext.SendPacket(pkt); err != nil {
	// 		log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
	// 	}

	// 	// Loop
	// 	for {
	// 		// Receive frame
	// 		if err := s.decCodecContext.ReceiveFrame(f); err != nil {
	// 			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
	// 				break
	// 			}
	// 			log.Fatal(fmt.Errorf("main: receiving fram failed: %w", err))
	// 		}

	// 		// Do something with decoded frame
	// 		// log.Printf("new frame: stream %d - pts: %d", pkt.StreamIndex(), f.Pts())
	// 	}
	// }

	// Success
	log.Println("success")
}
