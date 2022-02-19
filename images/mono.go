// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package images

import (
	"image"
	"image/color"
)

var MonoModel = color.ModelFunc(monoModel)

func monoModel(c color.Color) color.Color {
	if g, ok := c.(color.Gray); ok {
		return color.Gray{uint8(int8(g.Y) >> 7)}
	}
	r, g, b, _ := c.RGBA()
	y := int32(19595*r+38470*g+7471*b+1<<15) >> 31
	return color.Gray{uint8(y)}
}

// Mono is an in-memory image whose At method returns color.Gray with two
// possible values: 0 or 255.
type Mono struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	Shift  uint            // the bit offest in Pix[0] to the first pixel
	Pix    []byte          // the image pixels
}

// NewMono returns a new Mono image.
func NewMono(r image.Rectangle) *Mono {
	p := new(Mono)
	p.Rect = r
	p.Stride = (r.Dx() + 7) >> 3
	p.Pix = make([]uint8, p.Stride*r.Dy())
	return p
}

func (p *Mono) ColorModel() color.Model {
	return MonoModel
}

func (p *Mono) Bounds() image.Rectangle {
	return p.Rect
}

func (p *Mono) GrayAt(x, y int) color.Gray {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.Gray{}
	}
	i, s := p.PixOffset(x, y)
	return color.Gray{uint8(-int(p.Pix[i] >> s & 1))}
}

func (p *Mono) At(x, y int) color.Color {
	return p.GrayAt(x, y)
}

func (p *Mono) RGBA64At(x, y int) color.RGBA64 {
	l := uint16(p.GrayAt(x, y).Y)
	l |= l << 8
	return color.RGBA64{l, l, l, 0xffff}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index to the bit in that element that determines
// the pixel value.
func (p *Mono) PixOffset(x, y int) (offset int, shift uint) {
	x += int(p.Shift) - p.Rect.Min.X
	y -= p.Rect.Min.Y
	offset = y*p.Stride + x>>3
	shift = uint(x & 7)
	return
}

func (p *Mono) Set(x, y int, c color.Color) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	var bit uint8
	if g, ok := c.(color.Gray); ok {
		bit = uint8(int8(g.Y) >> 7)
	} else {
		r, g, b, _ := c.RGBA()
		bit = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 31)
	}
	i, s := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(1<<s) | bit<<s
}

func (p *Mono) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	r, g, b := uint32(c.R), uint32(c.G), uint32(c.B)
	bit := uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 31)
	i, s := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(1<<s) | bit<<s
}

func (p *Mono) SetGray(x, y int, c color.Gray) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	bit := uint8(int8(c.Y) >> 7)
	i, s := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(1<<s) | bit<<s
}
