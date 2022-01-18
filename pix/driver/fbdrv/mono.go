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

const (
	cblack = 0
	cwhite = 1
	ctrans = 2
)

type DriverMono struct {
	frame images.Mono
	flush func(frame *images.Mono) error
	color uint8
	err   error
}

func NewDriverMono(width, height int, flush func(frame *images.Mono) error) *DriverMono {
	d := new(DriverMono)
	d.frame.Rect.Max.X = width
	d.frame.Rect.Max.Y = height
	d.frame.Stride = width / 8
	d.frame.Pix = make([]uint8, d.frame.Stride*height)
	d.flush = flush
	return d
}

func (d *DriverMono) SetDir(dir int) image.Rectangle {
	return d.frame.Rect
}

func (d *DriverMono) Flush() {
	if d.flush != nil && d.err != nil {
		d.err = d.flush(&d.frame)
	}
}

func (d *DriverMono) Err(clear bool) error {
	err := d.err
	if clear {
		d.err = nil
	}
	return err
}

func (d *DriverMono) Frame() draw.Image {
	return &d.frame
}

func (d *DriverMono) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	// TODO
}

func (d *DriverMono) SetColor(c color.Color) {
	if g, ok := c.(color.Gray); ok {
		d.color = g.Y >> 7
		return
	}
	r, g, b, a := c.RGBA()
	if a>>15 == 0 {
		d.color = ctrans
	} else {
		y := 19595*r + 38470*g + 7471*b + 1<<15
		d.color = uint8(int32(y) >> 31)
	}
}

func (d *DriverMono) Fill(r image.Rectangle) {
	if d.color == ctrans {
		return
	}
	i, s := d.frame.PixOffset(r.Min.X, r.Min.Y)
	_ = i
	_ = s
}
