// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv

import (
	"image"
	"image/color"
	"image/draw"
)

const (
	cblack = 0
	cwhite = 1
	ctrans = 2
)

type Mono struct {
	fb     FrameBuffer
	pix    []byte
	width  int
	height int
	stride int
	shift  uint8
	mvxy   uint8
	color  uint8
}

func NewMono(fb FrameBuffer) *Mono {
	return &Mono{fb: fb}
}

func (d *Mono) SetDir(dir int) image.Rectangle {
	d.pix, d.width, d.height, d.stride, d.shift, d.mvxy = d.fb.SetDir(dir)
	var r image.Rectangle
	if d.mvxy&MV == 0 {
		r.Max.X = d.width
		r.Max.Y = d.height
	} else {
		r.Max.X = d.height
		r.Max.Y = d.width
	}
	return r
}

func (d *Mono) Flush()               { d.pix = d.fb.Flush() }
func (d *Mono) Err(clear bool) error { return d.fb.Err(clear) }

func monoPixOffset(d *Mono, x, y int) (offset int, shift uint) {
	x += int(d.shift)
	offset = y*d.stride + x>>3
	shift = uint(x & 7)
	return
}

func (d *Mono) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	var sl, sa uint = 0, 0xffff
	srcIsUniform := false
	srcGray, _ := src.(interface{ GrayAt(x, y int) color.Gray })
	if srcGray == nil {
		if u, ok := src.(*image.Uniform); ok {
			r, g, b, a := u.At(0, 0).RGBA()
			sl = lum32(r, g, b) >> 16
			sa = uint(a)
			srcIsUniform = true
		}
	}
	minx, miny := r.Min.X, r.Min.Y
	ox, oy := 1, d.stride*8
	if d.mvxy&MV != 0 {
		minx, miny = miny, minx
	}
	if d.mvxy&MX != 0 {
		minx = d.width - 1 - minx
		ox = -ox
	}
	if d.mvxy&MY != 0 {
		miny = d.height - 1 - miny
		oy = -oy
	}
	if d.mvxy&MV != 0 {
		ox, oy = oy, ox
	}
	offset, shift := monoPixOffset(d, minx, miny)
	offset = offset*8 + int(shift)
	width, height := r.Dx(), r.Dy()
	for y := 0; y < height; y++ {
		o := offset
		for x := 0; x < width; x++ {
			switch {
			case srcIsUniform:
				// sl, sa are constant
			case srcGray != nil:
				sl = uint(srcGray.GrayAt(sp.X+x, sp.Y+y).Y)
				sl |= sl << 8
				// sa is constant
			default:
				r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
				sl = lum32(r, g, b) >> 16
				sa = uint(a)
			}
			ma := uint(0xffff)
			if mask != nil {
				_, _, _, a := mask.At(mp.X+x, mp.Y+y).RGBA()
				ma = uint(a)
			}
			s := uint(o & 7)
			pix8 := uint(d.pix[o>>3])
			dl := sl
			switch {
			case sa&ma == 0xffff:
				// dl is ok
			case op == draw.Over:
				dl = -(pix8 >> s & 1) & 0xffff
				a := 0xffff - (sa * ma / 0xffff)
				dl = (dl*a + sl*ma) / 0xffff
			case ma != 0xffff:
				dl = sl * ma / 0xffff
			}
			d.pix[o>>3] = uint8(pix8&^(1<<s) | (dl>>15)<<s)
			o += ox
		}
		offset += oy
	}
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
		d.color = uint8(lum32(r, g, b) >> 31)
	}
}

func (d *Mono) Fill(r image.Rectangle) {
	if d.color == ctrans {
		return
	}
	if d.mvxy&MV != 0 {
		r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
		r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
	}
	if d.mvxy&MX != 0 {
		r.Max.X, r.Min.X = d.width-r.Min.X, d.width-r.Max.X
	}
	if d.mvxy&MY != 0 {
		r.Max.Y, r.Min.Y = d.height-r.Min.Y, d.height-r.Max.Y
	}
	offset, shift := monoPixOffset(d, r.Min.X, r.Min.Y)
	width := r.Dx()
	length := d.stride * r.Dy()
	n := 8 - int(shift)
	if width < n {
		n = width
	}
	color := uint8(-int(d.color))
	if n < 8 {
		rs := uint(8 - n)
		color0 := color >> rs << shift
		mask := uint8(0xff) >> rs << shift
		maxi := offset + length
		for i := offset; i < maxi; i += d.stride {
			d.pix[i] = d.pix[i]&^mask | color0
		}
		width -= n
		offset++
	}
	if n = width / 8; n != 0 {
		maxi := offset + length
		for i := offset; i < maxi; i += d.stride {
			for k, maxk := i, i+n; k < maxk; k++ {
				d.pix[k] = color
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
		for i := offset; i < maxi; i += d.stride {
			d.pix[i] = d.pix[i]&^mask | color
		}
	}
}

func lum32(r, g, b uint32) uint {
	return uint(19595*r + 38470*g + 7471*b + 1<<15)
}
