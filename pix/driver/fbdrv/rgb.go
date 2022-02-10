// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv

import (
	"image"
	"image/color"
	"image/draw"
)

type RGB struct {
	fb         FrameBuffer
	pix        []byte
	width      int
	height     int
	stride     int
	r, g, b, a uint16
	mvxy       uint8
}

func NewRGB(fb FrameBuffer) *RGB {
	return &RGB{fb: fb}
}

func (d *RGB) SetDir(dir int) image.Rectangle {
	d.pix, d.width, d.height, d.stride, _, d.mvxy = d.fb.SetDir(dir)
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

func (d *RGB) Flush()               { d.pix = d.fb.Flush() }
func (d *RGB) Err(clear bool) error { return d.fb.Err(clear) }

func (d *RGB) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	var sr, sg, sb, sa uint32
	srcIsUniform := false
	if u, ok := src.(*image.Uniform); ok {
		sr, sg, sb, sa = u.At(0, 0).RGBA()
		srcIsUniform = true
	}
	minx, miny := r.Min.X, r.Min.Y
	ox, oy := 3, d.stride
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
	offset := miny*d.stride + minx*3
	width, height := r.Dx(), r.Dy()
	for y := 0; y < height; y++ {
		o := offset
		for x := 0; x < width; x++ {
			if !srcIsUniform {
				sr, sg, sb, sa = src.At(sp.X+x, sp.Y+y).RGBA()
			}
			ma := uint32(0xffff)
			if mask != nil {
				_, _, _, ma = mask.At(mp.X+x, mp.Y+y).RGBA()
			}
			pix := d.pix[o : o+3 : o+3] // see https://golang.org/issue/27857
			dr, dg, db := sr, sg, sb
			if sa&ma != 0xffff {
				if op == draw.Over {
					dr = uint32(pix[0])
					dg = uint32(pix[1])
					db = uint32(pix[2])
					a := 0xffff - (sa * ma / 0xffff)
					dr = ((dr|dr<<8)*a + sr*ma) / 0xffff
					dg = ((dg|dg<<8)*a + sg*ma) / 0xffff
					db = ((db|db<<8)*a + sb*ma) / 0xffff
				} else if ma != 0xffff {
					dr = sr * ma / 0xffff
					dg = sg * ma / 0xffff
					db = sb * ma / 0xffff
				} else {
				}
			}
			pix[0] = uint8(dr >> 8)
			pix[1] = uint8(dg >> 8)
			pix[2] = uint8(db >> 8)
			o += ox
		}
		offset += oy
	}
}

func (d *RGB) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	d.r = uint16(r)
	d.g = uint16(g)
	d.b = uint16(b)
	d.a = uint16(a)
}

func (d *RGB) Fill(r image.Rectangle) {
	if d.a == 0 {
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
	offset := r.Min.Y*d.stride + r.Min.X*3
	maxi := offset + d.stride*r.Dy()
	n := r.Dx() * 3
	cr, cg, cb := uint(d.r), uint(d.g), uint(d.b)
	a := 0xffff - uint(d.a)
	for i := offset; i < maxi; i += d.stride {
		for k, maxk := i, i+n; k < maxk; k += 3 {
			pix := d.pix[k : k+3 : k+3] // see https://golang.org/issue/27857
			r := uint(pix[0])
			g := uint(pix[1])
			b := uint(pix[2])
			r = (r|r<<8)*a/0xffff + cr
			g = (g|g<<8)*a/0xffff + cg
			b = (b|b<<8)*a/0xffff + cb
			pix[0] = uint8(r >> 8)
			pix[1] = uint8(g >> 8)
			pix[2] = uint8(b >> 8)
		}
	}
}
