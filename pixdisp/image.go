// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
)

var (
	Alpha1Model color.Model = color.ModelFunc(alpha1Model)
	Alpha2Model color.Model = color.ModelFunc(alpha2Model)
)

func alpha1Model(c color.Color) color.Color {
	var alpha uint32
	if a, ok := c.(color.Alpha); ok {
		alpha = uint32(a.A) >> 7
	} else {
		_, _, _, alpha = c.RGBA()
		alpha >>= 15
	}
	return color.Alpha{uint8(-alpha)}
}

func alpha2Model(c color.Color) color.Color {
	var alpha uint32
	if a, ok := c.(color.Alpha); ok {
		alpha = uint32(a.A) >> 6
	} else {
		_, _, _, alpha = c.RGBA()
		alpha >>= 14
	}
	alpha |= alpha << 2
	alpha |= alpha << 4
	return color.Alpha{uint8(alpha)}
}

// Alpha1 is an in-memory image whose At method returns color.Alpha with only
// two possible values: 0, 255.
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

// ImmAlpha1 is like Alpha1 but it is immutable because its pixels are stored in
// a string.
type ImmAlpha1 struct {
	Rect   image.Rectangle
	Stride int
	Shift  int
	Pix    string
}

// NewImmAlpha1 returns a new ImmAlpha1 image with the given bounds and content.
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

// Alpha2 is an in-memory image whose At method returns color.Alpha with four
// possible values: 0, 85, 170, 255.
type Alpha2 struct {
	Rect   image.Rectangle
	Stride int
	Shift  int
	Pix    []uint8
}

// NewAlpha2 returns a new Alpha2 image with the given bounds.
func NewAlpha2(r image.Rectangle) *Alpha2 {
	stride := (r.Dx() + 3) / 4
	return &Alpha2{
		Rect:   r,
		Stride: stride,
		Pix:    make([]uint8, stride*r.Dy()),
	}
}

func (p *Alpha2) ColorModel() color.Model { return Alpha2Model }
func (p *Alpha2) Bounds() image.Rectangle { return p.Rect }

func (p *Alpha2) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	a := uint(p.Pix[i]) >> s & 3
	a |= a << 2
	a |= a << 4
	return color.Alpha{uint8(a)}
}

func (p *Alpha2) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index to bits in that element that determines the
// pixel value.
func (p *Alpha2) PixOffset(x, y int) (offset, shift int) {
	x += p.Shift - p.Rect.Min.X
	y -= p.Rect.Min.Y
	col := x / 4
	offset = y*p.Stride + col
	shift = (x - col*4) * 2
	return
}

func (p *Alpha2) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	a := alpha2Model(c).(color.Alpha)
	i, shift := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(3<<shift) | a.A>>6<<shift
}

func (p *Alpha2) SetAlpha(x, y int, c color.Alpha) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i, shift := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(3<<shift) | c.A>>6<<shift
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Alpha2) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Alpha2{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &Alpha2{
		Rect:   r,
		Stride: p.Stride,
		Shift:  shift,
		Pix:    p.Pix[i:],
	}
}

// ImmAlpha2 is like Alpha2 but it is immutable because its pixels are stored in
// a string.
type ImmAlpha2 struct {
	Rect   image.Rectangle
	Stride int
	Shift  int
	Pix    string
}

// NewImmAlpha2 returns a new ImmAlpha2 image with the given bounds and content.
func NewImmAlpha2(r image.Rectangle, bits string) *ImmAlpha2 {
	return &ImmAlpha2{
		Rect:   r,
		Stride: (r.Dx() + 3) / 4,
		Pix:    bits,
	}
}

func (p *ImmAlpha2) ColorModel() color.Model { return Alpha2Model }
func (p *ImmAlpha2) Bounds() image.Rectangle { return p.Rect }

func (p *ImmAlpha2) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	a := uint(p.Pix[i]) >> s & 3
	a |= a << 2
	a |= a << 4
	return color.Alpha{uint8(a)}
}

func (p *ImmAlpha2) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index to bits in that element that determines the
// pixel value.
func (p *ImmAlpha2) PixOffset(x, y int) (offset, shift int) {
	x += p.Shift - p.Rect.Min.X
	y -= p.Rect.Min.Y
	col := x / 4
	offset = y*p.Stride + col
	shift = (x - col*4) * 2
	return
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImmAlpha2) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImmAlpha2{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImmAlpha2{
		Rect:   r,
		Stride: p.Stride,
		Shift:  shift,
		Pix:    p.Pix[i:],
	}
}

// ImmAlpha is like image.Alpha but it is immutable because its pixels are
// stored in a string.
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
