// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/color"
	"image/draw"
)

// Driver defines the interface to be implemented by display controller drivers
// or anything else that wants to use pix package drawing capability.
type Driver interface {
	// SetDir sets the display direction rotating its default coordinate system
	// by dir*90 degrees. It returns the bounds of the frame memory for the new
	// direction. NewDisplay calls SetDir(0) before any other Driver method is
	// called so this is a good place for initialization code if required.
	SetDir(dir int) image.Rectangle

	// Draw works like draw.DrawMask with dst set to the image representing the
	// whole frame memory.
	//
	// The draw.Over operator may be implemented in a limited way but it must
	// at least support 1-bit transparency.
	//
	// Draw may assume the r is a non-empty rectangle that fits entirely on the
	// display and is entirely covered by src and mask.
	//
	// Draw is actually the only drawing operation required from a display
	// controller. The Fill operation below can be easily implemented by drawing
	// an uniform image.
	Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op)

	// SetColor sets the color used by Fill method.
	SetColor(c color.Color)

	// Fill helps to increase prformance when drawing filled rectangles which
	// are heavily used when drawing various geometric shapes.
	//
	// Fill(r) is intended to be a faster counterpart of
	//
	//	Draw(r, image.Uniform{c}, image.Point{}, nil, image.Point{}, draw.Over)
	//
	// where c is the color set by SetColor.
	Fill(r image.Rectangle)

	// Flush allows to flush the drivers internal buffers. Drivers are allowed
	// to implement any kind of buffering if the direct drawing to the frame
	// memory is problematic or inefficient.
	Flush()

	// Err returns the saved error and clears it if the clear is true. If an
	// error has occured it is recommended that the Driver avoids any further
	// operations until the error is cleared.
	Err(clear bool) error
}

/*

// ReadDriver is an optional interface that driver may implement to allow read
// the content of frame memory.
type ReadDriver interface {
	// NewImage returns a new image with the given bounds. The returned image
	// type is best fitted to be used with Read method.
	NewImage(r image.Rectangle) draw.Image

	// Read reads the part of frame memory specified by r into dst at dp.
	Read(r image.Rectangle, dst draw.Image, dp image.Point)
}

*/
