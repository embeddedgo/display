// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/embeddedgo/display/images"
)

type RGB struct {
	frame images.RGB
	flush func(frame *images.RGB) error
	color     color.RGBA64
	err   error
}

func NewRGB(width, height int, flush func(frame *images.RGB) error) *RGB {
	d := new(RGB)
	d.frame.Rect.Max.X = width
	d.frame.Rect.Max.Y = height
	d.frame.Stride = 3 * width
	d.frame.Pix = make([]uint8, d.frame.Stride*height)
	d.flush = flush
	return d
}

func (d *RGB) SetDir(dir int) image.Rectangle {
	return d.frame.Rect
}

func (d *RGB) Flush() {
	if d.flush != nil && d.err != nil {
		d.err = d.flush(&d.frame)
	}
}

func (d *RGB) Err(clear bool) error {
	err := d.err
	if clear {
		d.err = nil
	}
	return err
}

func (d *RGB) Frame() draw.Image {
	return &d.frame
}

func (d *RGB) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	// TODO
}

func (d *RGB) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	d.color.R = uint16(r)
	d.color.G = uint16(g)
	d.color.B = uint16(b)
	d.color.A = uint16(a)
}

func (d *RGB) Fill(r image.Rectangle) {
	// TODO
}
