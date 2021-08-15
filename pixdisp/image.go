// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
)

var MonoModel color.Model = color.ModelFunc(monoModel)

func monoModel(c color.Color) color.Color {
	var alpha uint32
	if a, ok := c.(color.Alpha); ok {
		alpha = uint32(a.A)
	} else {
		_, _, _, alpha = c.RGBA()
	}
	if alpha != 0 {
		alpha = 0xff
	}
	return color.Alpha{uint8(alpha)}
}

// Mono is an in-memory image whose At method returns color.Alpha with only two
// possible values: 0, 0xff.
type Mono struct {
	Rect   image.Rectangle
	Stride int
	Shift  int
	Pix    []uint8
}

// NewMono returns a new Mono image with the given bounds.
func NewMono(r image.Rectangle) *Mono {
	stride := (r.Dx() + 7) / 8
	return &Mono{
		Rect:   r,
		Stride: stride,
		Pix:    make([]uint8, stride*r.Dy()),
	}
}

func (p *Mono) ColorModel() color.Model { return MonoModel }
func (p *Mono) Bounds() image.Rectangle { return p.Rect }

func (p *Mono) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	return color.Alpha{p.Pix[i] >> s & 1 * 0xff}
}

func (p *Mono) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index of bit in that element that determines the
// pixel value.
func (p *Mono) PixOffset(x, y int) (offset, shift int) {
	x += p.Shift - p.Rect.Min.X
	y -= p.Rect.Min.Y
	col := x / 8
	offset = y*p.Stride + col
	shift = x - col*8
	return
}

func (p *Mono) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	a := monoModel(c).(color.Alpha)
	i, shift := p.PixOffset(x, y)
	mask := uint8(1 << shift)
	p.Pix[i] = p.Pix[i]&^mask | a.A&mask
}

func (p *Mono) SetAlpha(x, y int, c color.Alpha) {
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
func (p *Mono) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Mono{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &Mono{
		Rect:   r,
		Stride: p.Stride,
		Shift:  shift,
		Pix:    p.Pix[i:],
	}
}

// ImmMono is like Mono but its pixels are stored in string so it is immutable.
type ImmMono struct {
	Rect   image.Rectangle
	Stride int
	Shift  int
	Pix    string
}

// NewImmMono returns a new ImmMono image with the given bounds and the pixels
// values stored in string.
func NewImmMono(r image.Rectangle, bits string) *ImmMono {
	stride := (r.Dx() + 7) / 8
	return &ImmMono{
		Rect:   r,
		Stride: stride,
		Pix:    bits,
	}
}

func (p *ImmMono) ColorModel() color.Model { return MonoModel }
func (p *ImmMono) Bounds() image.Rectangle { return p.Rect }

func (p *ImmMono) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	return color.Alpha{p.Pix[i] >> s & 1 * 0xff}
}

func (p *ImmMono) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index of bit in that element that determines the
// pixel value.
func (p *ImmMono) PixOffset(x, y int) (offset, shift int) {
	x += p.Shift - p.Rect.Min.X
	y -= p.Rect.Min.Y
	col := x / 8
	offset = y*p.Stride + col
	shift = x - col*8
	return
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImmMono) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImmMono{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImmMono{
		Rect:   r,
		Stride: p.Stride,
		Shift:  shift,
		Pix:    p.Pix[i:],
	}
}
