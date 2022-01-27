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

type Mono struct {
	frame images.Mono
	flush func(frame *images.Mono) error
	color uint8
	err   error
}

func NewMono(width, height int, flush func(frame *images.Mono) error) *Mono {
	d := new(Mono)
	d.frame.Rect.Max.X = width
	d.frame.Rect.Max.Y = height
	d.frame.Stride = (width + 7) >> 3
	d.frame.Pix = make([]uint8, d.frame.Stride*height)
	d.flush = flush
	return d
}

func (d *Mono) SetDir(dir int) image.Rectangle {
	return d.frame.Rect
}

func (d *Mono) Flush() {
	if d.flush != nil && d.err == nil {
		d.err = d.flush(&d.frame)
	}
}

func (d *Mono) Err(clear bool) error {
	err := d.err
	if clear {
		d.err = nil
	}
	return err
}

func (d *Mono) Frame() draw.Image {
	return &d.frame
}

func (d *Mono) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	// TODO
}

func (d *Mono) SetColor(c color.Color) {
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

func (d *Mono) Fill(r image.Rectangle) {
	if d.color == ctrans {
		return
	}
	color := uint8(-int(d.color))
	offset, shift := d.frame.PixOffset(r.Min.X, r.Min.Y)
	width := r.Dx()
	length := d.frame.Stride * r.Dy()
	n := 8 - int(shift)
	if width < n {
		n = width
	}
	if n < 8 {
		rs := uint(8 - n)
		color0 := color >> rs << shift
		mask := uint8(0xff) >> rs << shift
		maxi := offset + length
		for i := offset; i < maxi; i += d.frame.Stride {
			d.frame.Pix[i] = d.frame.Pix[i]&^mask | color0
		}
		width -= n
		offset++
	}
	if n = width / 8; n != 0 {
		maxi := offset + length
		for i := offset; i < maxi; i += d.frame.Stride {
			for k, maxk := i, i+n; k < maxk; k++ {
				d.frame.Pix[k] = color
			}
		}
		offset += n
		width -= n * 8
	}
	if width != 0 {
		rs := uint(8 - width)
		color >>= rs
		mask := uint8(0xff) >> rs
		maxi := offset + length
		for i := offset; i < maxi; i += d.frame.Stride {
			d.frame.Pix[i] = d.frame.Pix[i]&^mask | color
		}
	}
}
