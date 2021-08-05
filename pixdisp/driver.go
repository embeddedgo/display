// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
)

// Driver lists the operations expected from a display driver.
type Driver interface {
	// Dim returns the display dimensions.
	Dim() (width, height int)

	// Draw draws an image. It works similar to the draw.Draw function with op
	// set to draw.Src.
	//
	// Draw is actually the only operation required from a display cntroller.
	// The Fill operation below can be easily implemented by drawing an uniform
	// image.
	Draw(r image.Rectangle, src image.Image, sp image.Point)

	// Color allows to convert any color to the drivers/controller internal
	// representation. It helps to reduce the number of required color
	// conversions.
	Color(c color.Color) uint64

	// Fill helps to increase prformance when drawing filled rectangles. Any
	// geometric shape or font is drawed using finite number of filed
	// rectangles so its heavily used operation which is worth optimizing.
	Fill(r image.Rectangle, c uint64)

	// Flush allows to flush the drivers internal buffers. Drivers is allowed to
	// implement any kind of buffering if the direct drawing to the display is
	// problematic or inefficient.
	Flush()

	// Err returns the saved error and clears it if the clear is true. If an
	// error has occured it is recommended that the Driver avoid any further
	// operation until the error is cleared.
	Err(clear bool) error
}
