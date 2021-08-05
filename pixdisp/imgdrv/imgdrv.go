// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package imgdrv provides a driver based on Go image and image/draw packages.
// It allows to use pixdisp package without any real display, mainly for tests.
package imgdrv

import (
	"image"
	"image/color"
	"image/draw"
)

type Driver struct{ Image draw.Image }

func (d Driver) Dim() (width, height int) {
	r := d.Image.Bounds()
	return r.Dx(), r.Dy()
}

func (d Driver) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	draw.Draw(d.Image, r.Add(d.Image.Bounds().Min), src, sp, draw.Src)
}

func (d Driver) Color(c color.Color) uint64 {
	r, g, b, a := c.RGBA()
	return uint64(r<<16|g)<<32 | uint64(b<<16|a)
}

func (d Driver) Fill(r image.Rectangle, c uint64) {
	rgba64 := color.RGBA64{
		uint16(c >> 48),
		uint16(c >> 32),
		uint16(c >> 16),
		uint16(c),
	}
	d.Draw(r, &image.Uniform{rgba64}, image.Point{})
}

func (d Driver) Flush()               {}
func (d Driver) Err(clear bool) error { return nil }
