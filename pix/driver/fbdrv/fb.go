// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fbdrv provides frame-buffer based drivers for pix graphics library.
// It is intended to be used for the displays without any internal frame buffer
// or for small displays with the internal frame buffer accessible via
// write-only communication interface. It can be also used instead of ../imgdrv
// to draw on any in memory image with added benefit of working SetDir and Flush
// methods.
//
// The subpackages provides drivers to the specific display controllers.
package fbdrv

// Coordination system translation constants
const (
	MV = 1 << iota // swap X with Y
	MX             // mirror X axis
	MY             // mirror Y axis
)

// FrameBuffer defines the interface to an in RAM frame buffer.
type FrameBuffer interface {
	// SetDir returns the frame buffer in buf and its geometry in width, height
	// and stride. In case of sub-byte pixels the bit shift to the first pixel
	// in buf[0] is provided in shift. The mvxy translation (described as
	// a combination of MV, MX, MY constants) should be aplied when drawing
	// in the returned buffer to obtain the desired direction (rotation) of the
	// image on the display. Every SetDir call may return a different buffer
	// with a different geometry and random content.
	SetDir(dir int) (buf []byte, width, height, stride int, shitf, mvxy uint8)

	// Flush exports the content of the internal RAM buffer to the actual
	// display or to an other media like file, etc. The data in buffer may be
	// also used to produce the visible image in real time. In such case Flush
	// may swap buffers if multi-buffering is implemented, sleep until the
	// beginning of the next V-blank period or simply do nothing. Flush returns
	// a buffer which should be used for subsequent drawing operations. The
	// returned buffer may contain random data (BUG: really so?).
	Flush() (buf []byte)

	// Err returns the saved error and clears it if the clear is true.
	Err(clear bool) error
}
