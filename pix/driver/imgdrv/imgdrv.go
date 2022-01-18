// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package imgdrv is an example of a minimal driver based on Go image/draw
// package. This package is intended mainly for tests. See ../fbdrv for similar
// but more efficient drivers.
package imgdrv

import (
	"image"
	"image/color"
	"image/draw"
)

// Driver provides simplest possible implementation of pix.Driver. Its speed is
// not important.
type Driver struct {
	Image draw.Image
	fill  image.Uniform
}

// New returns new Driver using img as frame memory.
func New(img draw.Image) *Driver {
	return &Driver{Image: img}
}

// SetDir implements pix.Driver.
func (d *Driver) SetDir(dir int) image.Rectangle {
	return d.Image.Bounds()
}

// Draw implements pix.Driver.
func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	draw.DrawMask(d.Image, r, src, sp, mask, mp, op)
}

// SetColor implements pix.Driver.
func (d *Driver) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	d.fill.C = color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// Fill implements pix.Driver.
func (d *Driver) Fill(r image.Rectangle) {
	d.Draw(r, &d.fill, image.Point{}, nil, image.Point{}, draw.Over)
}

// Flush implements pix.Driver.
func (d *Driver) Flush() {}

// Err implements pix.Driver.
func (d *Driver) Err(clear bool) error { return nil }
