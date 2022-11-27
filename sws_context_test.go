package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
)

// TestNewSwsContext
func TestNewSwsContext(t *testing.T) {
	// Init
	var (
		s = astiav.NewSwsContext(
			100, 100, astiav.PixelFormatYuv420P,
			100, 100, astiav.PixelFormatYuv420P,
			0, nil, nil,
		)
	)
	// Free
	s.Free()
}

// TestInitContext

// TestCachedContext

// TestScale

// TestScaleFrames

// TestScaleDstFrame

// TestScaleSrcFrame
