// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
)

const (
	Alpha1Model AlphaNModel = 1
	Alpha2Model AlphaNModel = 2
	Alpha4Model AlphaNModel = 4
)

type AlphaNModel uint8

func (m AlphaNModel) Convert(c color.Color) color.Color {
	var alpha uint32
	if a, ok := c.(color.Alpha); ok {
		alpha = uint32(a.A) >> (8 - m)
	} else {
		_, _, _, alpha = c.RGBA()
		alpha >>= 16 - m
	}
	if m == 1 {
		alpha = -alpha
	} else {
		if m == 2 {
			alpha |= alpha << 2
		}
		alpha |= alpha << 4
	}
	return color.Alpha{uint8(alpha)}
}

// AlphaN is an in-memory image whose At method returns color.Alpha with 1, 2 or
// 4 bit precision.
type AlphaN struct {
	Rect   image.Rectangle
	Stride int
	LogN   uint8
	Shift  uint8
	Pix    []uint8
}

// NewAlphaN returns a new AlphaN image with the given bounds and number of bits
// per pixel.
func NewAlphaN(r image.Rectangle, n int) *AlphaN {
	p := new(AlphaN)
	p.Rect = r
	switch n {
	case 1:
		p.LogN = 0
	case 2:
		p.LogN = 1
	case 4:
		p.LogN = 2
	default:
		panic("unsupported bpp")
	}
	p.Stride = (r.Dx() + 7>>p.LogN) >> (3 - p.LogN)
	p.Pix = make([]uint8, p.Stride*r.Dy())
	return p
}

func (p *AlphaN) ColorModel() color.Model { return AlphaNModel(1 << p.LogN) }
func (p *AlphaN) Bounds() image.Rectangle { return p.Rect }

func (p *AlphaN) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	a := uint(p.Pix[i]) >> s
	if p.LogN == 0 {
		a = -(a & 1)
	} else {
		if p.LogN == 1 {
			a &= 3
			a |= a << 2
		} else {
			a &= 15
		}
		a |= a << 4
	}
	return color.Alpha{uint8(a)}
}

func (p *AlphaN) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index to the bits in that element that
// determines the pixel value.
func (p *AlphaN) PixOffset(x, y int) (offset int, shift uint) {
	x += int(p.Shift)>>p.LogN - p.Rect.Min.X
	y -= p.Rect.Min.Y
	cs := 3 - p.LogN
	col := x >> cs
	offset = y*p.Stride + col
	shift = uint(x-col<<cs) << p.LogN
	return
}

func (p *AlphaN) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	var alpha uint32
	if a, ok := c.(color.Alpha); ok {
		alpha = uint32(a.A)
	} else {
		_, _, _, alpha = c.RGBA()
		alpha >>= 8
	}
	rshift := uint(8) - 1<<p.LogN
	i, lshift := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(0xff>>rshift<<lshift) | uint8(alpha>>rshift<<lshift)
}

func (p *AlphaN) SetAlpha(x, y int, c color.Alpha) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	rshift := uint(8) - 1<<p.LogN
	i, lshift := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(0xff>>rshift<<lshift) | c.A>>rshift<<lshift
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *AlphaN) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &AlphaN{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &AlphaN{
		Rect:   r,
		LogN:   p.LogN,
		Shift:  uint8(shift),
		Stride: p.Stride,
		Pix:    p.Pix[i:],
	}
}

// ImmAlphaN is immutable counterpart of AlphaN.
type ImmAlphaN struct {
	Rect   image.Rectangle
	Stride int
	LogN   uint8
	Shift  uint8
	Pix    string
}

// NewImmAlphaN returns a new ImmAlpha image with the given bounds and content.
func NewImmAlphaN(r image.Rectangle, n int, bits string) *ImmAlphaN {
	p := new(ImmAlphaN)
	p.Rect = r
	switch n {
	case 1:
		p.LogN = 0
	case 2:
		p.LogN = 1
	case 4:
		p.LogN = 2
	default:
		panic("unsupported bpp")
	}
	p.Stride = (r.Dx() + 7>>p.LogN) >> (3 - p.LogN)
	p.Pix = bits
	return p
}

func (p *ImmAlphaN) ColorModel() color.Model { return AlphaNModel(1 << p.LogN) }
func (p *ImmAlphaN) Bounds() image.Rectangle { return p.Rect }

func (p *ImmAlphaN) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	a := uint(p.Pix[i]) >> s
	if p.LogN == 0 {
		a = -(a & 1)
	} else {
		if p.LogN == 1 {
			a &= 3
			a |= a << 2
		} else {
			a &= 15
		}
		a |= a << 4
	}
	return color.Alpha{uint8(a)}
}

func (p *ImmAlphaN) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index to the bits in that element that
// determines the pixel value.
func (p *ImmAlphaN) PixOffset(x, y int) (offset int, shift uint) {
	x += int(p.Shift)>>p.LogN - p.Rect.Min.X
	y -= p.Rect.Min.Y
	cs := 3 - p.LogN
	col := x >> cs
	offset = y*p.Stride + col
	shift = uint(x-col<<cs) << p.LogN
	return
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImmAlphaN) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &AlphaN{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImmAlphaN{
		Rect:   r,
		LogN:   p.LogN,
		Shift:  uint8(shift),
		Stride: p.Stride,
		Pix:    p.Pix[i:],
	}
}

// ImmAlpha is immutable counterpart of image.Alpha.
type ImmAlpha struct {
	Rect   image.Rectangle
	Stride int
	Pix    string
}

// NewImmAlpha returns a new ImmAlpha image with the given bounds and content.
func NewImmAlpha(r image.Rectangle, bits string) *ImmAlpha {
	return &ImmAlpha{
		Rect:   r,
		Stride: r.Dx(),
		Pix:    bits,
	}
}

func (p *ImmAlpha) ColorModel() color.Model { return color.AlphaModel }
func (p *ImmAlpha) Bounds() image.Rectangle { return p.Rect }

func (p *ImmAlpha) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i := p.PixOffset(x, y)
	return color.Alpha{p.Pix[i]}
}

func (p *ImmAlpha) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index to bits in that element that determines the
// pixel value.
func (p *ImmAlpha) PixOffset(x, y int) int {
	x -= p.Rect.Min.X
	y -= p.Rect.Min.Y
	return y*p.Stride + x
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImmAlpha) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImmAlpha{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImmAlpha{
		Rect:   r,
		Stride: p.Stride,
		Pix:    p.Pix[i:],
	}
}
