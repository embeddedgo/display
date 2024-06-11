// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dummydrv provides a dummy implementation of pix.Driver. It is
// intended to help bencharking the pix library and its drivers, providing
// the fastest possible implementation that real drivers can be compared to.
package dummydrv

import (
	"image"
	"image/color"
	"image/draw"
)

// A Driver provides a no-op implementation of pix.Driver. All its methods
// except SetDir does nothing. SetDir returns the rectangle set by the New
// function.
type Driver struct {
	r image.Rectangle
}

// New returns a new Driver that simulates a display with given bounds.
func New(r image.Rectangle) *Driver { return &Driver{r} }

func (d *Driver) SetDir(dir int) image.Rectangle { return d.r }
func (d *Driver) SetColor(c color.Color)         { return }
func (d *Driver) Fill(r image.Rectangle)         { return }
func (d *Driver) Flush()                         { return }
func (d *Driver) Err(clear bool) error           { return nil }

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	return
}
