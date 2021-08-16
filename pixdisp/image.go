// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
)

var Alpha1Model color.Model = color.ModelFunc(alpha1Model)

func alpha1Model(c color.Color) color.Color {
	var alpha uint32
	if a, ok := c.(color.Alpha); ok {
		alpha = uint32(a.A)
	} else {
		_, _, _, alpha = c.RGBA()
	}
	if alpha >= 0x80 {
		alpha = 0xff
	}
	return color.Alpha{uint8(alpha)}
}

// Alpha1 is an in-memory image whose At method returns color.Alpha with only
// two possible values: 0, 0xff.
type Alpha1 struct {
	Rect   image.Rectangle
	Stride int
	Shift  int
	Pix    []uint8
}

// NewAlpha1 returns a new Alpha1 image with the given bounds.
func NewAlpha1(r image.Rectangle) *Alpha1 {
	stride := (r.Dx() + 7) / 8
	return &Alpha1{
		Rect:   r,
		Stride: stride,
		Pix:    make([]uint8, stride*r.Dy()),
	}
}

func (p *Alpha1) ColorModel() color.Model { return Alpha1Model }
func (p *Alpha1) Bounds() image.Rectangle { return p.Rect }

func (p *Alpha1) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	return color.Alpha{-(p.Pix[i] >> s & 1)}
}

func (p *Alpha1) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index of bit in that element that determines the
// pixel value.
func (p *Alpha1) PixOffset(x, y int) (offset, shift int) {
	x += p.Shift - p.Rect.Min.X
	y -= p.Rect.Min.Y
	col := x / 8
	offset = y*p.Stride + col
	shift = x - col*8
	return
}

func (p *Alpha1) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	a := alpha1Model(c).(color.Alpha)
	i, shift := p.PixOffset(x, y)
	mask := uint8(1 << shift)
	p.Pix[i] = p.Pix[i]&^mask | a.A&mask
}

func (p *Alpha1) SetAlpha(x, y int, c color.Alpha) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i, shift := p.PixOffset(x, y)
	mask := uint8(1 << shift)
	if c.A != 0 {
		c.A = 0xff
	}
	p.Pix[i] = p.Pix[i]&^mask | c.A&mask
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Alpha1) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Alpha1{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &Alpha1{
		Rect:   r,
		Stride: p.Stride,
		Shift:  shift,
		Pix:    p.Pix[i:],
	}
}

// ImmAlpha1 is like Alpha1 but its pixels are stored in a string so it is
// immutable.
type ImmAlpha1 struct {
	Rect   image.Rectangle
	Stride int
	Shift  int
	Pix    string
}

// NewImmAlpha1 returns a new ImmMono image with the given bounds and the
// pixel values stored in string.
func NewImmAlpha1(r image.Rectangle, bits string) *ImmAlpha1 {
	stride := (r.Dx() + 7) / 8
	return &ImmAlpha1{
		Rect:   r,
		Stride: stride,
		Pix:    bits,
	}
}

func (p *ImmAlpha1) ColorModel() color.Model { return Alpha1Model }
func (p *ImmAlpha1) Bounds() image.Rectangle { return p.Rect }

func (p *ImmAlpha1) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	return color.Alpha{p.Pix[i] >> s & 1 * 0xff}
}

func (p *ImmAlpha1) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index of bit in that element that determines the
// pixel value.
func (p *ImmAlpha1) PixOffset(x, y int) (offset, shift int) {
	x += p.Shift - p.Rect.Min.X
	y -= p.Rect.Min.Y
	col := x / 8
	offset = y*p.Stride + col
	shift = x - col*8
	return
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImmAlpha1) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImmAlpha1{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImmAlpha1{
		Rect:   r,
		Stride: p.Stride,
		Shift:  shift,
		Pix:    p.Pix[i:],
	}
}
