// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/color"
	"image/draw"
)

// Driver lists the operations expected from a display driver.
type Driver interface {
	// Size returns the display dimensions.
	Size() image.Point

	// Draw works like draw.DrawMask with dst set to the image representing the
	// whole display.
	//
	// The draw.Over operator can be implemented in a limited way but it must
	// at least support one-bit transparency.
	//
	// Draw can assume that r is a non-empty rectangle that fits entirely on the
	// display and is entirely covered by src and mask.
	//
	// Draw is actually the only operation required from a display controller.
	// The Fill operation below can be easily implemented by drawing an uniform
	// image.
	Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op)

	// SetColor sets the color used by Fill method.
	SetColor(c color.Color)

	// Fill helps to increase prformance when drawing filled rectangles which
	// are heavily used in drawing various geometric shapes.
	Fill(r image.Rectangle)

	// Flush allows to flush the drivers internal buffers. Drivers are allowed
	// to implement any kind of buffering if the direct drawing to the display
	// is problematic or inefficient.
	Flush()

	// Err returns the saved error and clears it if the clear is true. If an
	// error has occured it is recommended that the Driver avoids any further
	// operations until the error is cleared.
	Err(clear bool) error
}
