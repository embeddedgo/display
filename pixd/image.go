// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/color"
)

const (
	Alpha1Model AlphaNModel = 1 // 2 levels of transparency
	Alpha2Model AlphaNModel = 2 // 4 levels of transparency
	Alpha4Model AlphaNModel = 4 // 16 levels of transparency
	Alpha8Model AlphaNModel = 8 // 256 levels of transparency
)

// AlphaNModel is a color model for n-bit transparency
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

var (
	RGBModel   = color.ModelFunc(rgbModel)
	RGB16Model = color.ModelFunc(rgb16Model)
)

func rgbModel(c color.Color) color.Color {
	if c, ok := c.(color.RGBA); ok {
		c.A = 255
		return c
	}
	r, g, b, _ := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
}

func rgb16Model(c color.Color) color.Color {
	var r, g, b uint32
	if c, ok := c.(color.RGBA); ok {
		r = uint32(c.R)
		g = uint32(c.G)
		b = uint32(c.B)
	} else {
		r, g, b, _ = c.RGBA()
		r >>= 8
		g >>= 8
		b >>= 8
	}
	return color.RGBA{
		uint8(r&^7 | r>>5),
		uint8(g&^3 | g>>6),
		uint8(b&^7 | b>>5),
		255,
	}
}

func rgb16torgba(h, l uint8) color.RGBA {
	r := h &^ 7
	g := (h<<5 | l>>3) &^ 3
	b := l << 3
	return color.RGBA{
		r | r>>5,
		g | g>>6,
		b | b>>5,
		255,
	}
}

func logNStride(r image.Rectangle, nbpp int) (logN uint8, stride int) {
	switch nbpp {
	case 1:
		logN = 0
	case 2:
		logN = 1
	case 4:
		logN = 2
	case 8:
		logN = 3
	default:
		panic("unsupported bpp")
	}
	stride = (r.Dx() + 7>>logN) >> (3 - logN)
	return
}

func alphaLogN(a uint, logN uint8) color.Alpha {
	if logN == 0 {
		a = -(a & 1)
	} else {
		if logN == 1 {
			a &= 3
			a |= a << 2
		} else if logN < 3 {
			a &= 15
		}
		a |= a << 4
	}
	return color.Alpha{uint8(a)}
}

// AlphaN is an in-memory image whose At method returns color.Alpha with
// 1, 2, 4 or 8 bit precision.
type AlphaN struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	LogN   uint8           // 1<<LogN is the number of bits per pixel
	Shift  uint8           // the bit offest in Pix[0] to the first pixel
	Pix    []uint8         // the image pixels
}

// NewAlphaN returns a new AlphaN image with the given bounds and
// number of bits per pixel.
func NewAlphaN(r image.Rectangle, nbpp int) *AlphaN {
	p := new(AlphaN)
	p.Rect = r
	p.LogN, p.Stride = logNStride(r, nbpp)
	p.Pix = make([]uint8, p.Stride*r.Dy())
	return p
}

func (p *AlphaN) ColorModel() color.Model {
	return AlphaNModel(1 << p.LogN)
}

func (p *AlphaN) Bounds() image.Rectangle {
	return p.Rect
}

func (p *AlphaN) AlphaAt(x, y int) color.Alpha {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	return alphaLogN(uint(p.Pix[i])>>s, p.LogN)
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
	if !(image.Pt(x, y).In(p.Rect)) {
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
	if !(image.Pt(x, y).In(p.Rect)) {
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

// ImmAlphaN is an immutable counterpart of AlphaN.
type ImmAlphaN struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	LogN   uint8           // 1<<LogN is the number of bits per pixel
	Shift  uint8           // the bit offest in Pix[0] to the first pixel
	Pix    string          // the image pixels
}

// NewImmAlphaN returns a new ImmAlpha image with the given bounds and content.
func NewImmAlphaN(r image.Rectangle, nbpp int, bits string) *ImmAlphaN {
	p := new(ImmAlphaN)
	p.Rect = r
	p.LogN, p.Stride = logNStride(r, nbpp)
	p.Pix = bits
	return p
}

func (p *ImmAlphaN) ColorModel() color.Model { return AlphaNModel(1 << p.LogN) }
func (p *ImmAlphaN) Bounds() image.Rectangle { return p.Rect }

func (p *ImmAlphaN) AlphaAt(x, y int) color.Alpha {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	return alphaLogN(uint(p.Pix[i])>>s, p.LogN)
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
		return &ImmAlphaN{}
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

// RGB is an in-memory image whose At method returns RGB values.
type RGB struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	Pix    []uint8         // the image pixels
}

// NewRGB returns a new RGB image with the given bounds.
func NewRGB(r image.Rectangle) *RGB {
	return &RGB{
		Rect:   r,
		Stride: 3 * r.Dx(),
		Pix:    make([]uint8, 3*r.Dx()*r.Dy()),
	}
}

func (p *RGB) ColorModel() color.Model { return RGBModel }
func (p *RGB) Bounds() image.Rectangle { return p.Rect }
func (p *RGB) Opaque() bool            { return true }

func (p *RGB) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *RGB) RGBAAt(x, y int) color.RGBA {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	return color.RGBA{s[0], s[1], s[2], 255}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGB) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *RGB) Set(x, y int, c color.Color) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := RGBModel.Convert(c).(color.RGBA)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
}

func (p *RGB) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGB) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGB{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGB{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// ImmRGB is an immutable counterpart of RGB.
type ImmRGB struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	Pix    string          // the image pixels
}

// NewImmRGB returns a new ImmRGB image with the given bounds and content
func NewImmRGB(r image.Rectangle, bits string) *ImmRGB {
	return &ImmRGB{
		Rect:   r,
		Stride: 3 * r.Dx(),
		Pix:    bits,
	}
}

func (p *ImmRGB) ColorModel() color.Model { return RGBModel }
func (p *ImmRGB) Bounds() image.Rectangle { return p.Rect }
func (p *ImmRGB) Opaque() bool            { return true }

func (p *ImmRGB) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *ImmRGB) RGBAAt(x, y int) color.RGBA {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3]
	return color.RGBA{s[0], s[1], s[2], 255}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ImmRGB) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImmRGB) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImmRGB{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImmRGB{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// RGB16 is an in-memory image whose At method returns RGB16 values.
type RGB16 struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	Pix    []uint8         // the image pixels
}

// NewRGB16 returns a new RGB16 image with the given bounds.
func NewRGB16(r image.Rectangle) *RGB16 {
	return &RGB16{
		Rect:   r,
		Stride: 2 * r.Dx(),
		Pix:    make([]uint8, 2*r.Dx()*r.Dy()),
	}
}

func (p *RGB16) ColorModel() color.Model { return RGB16Model }
func (p *RGB16) Bounds() image.Rectangle { return p.Rect }
func (p *RGB16) Opaque() bool            { return true }

func (p *RGB16) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *RGB16) RGBAAt(x, y int) color.RGBA {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	return rgb16torgba(s[0], s[1])
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGB16) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
}

func (p *RGB16) Set(x, y int, c color.Color) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	var r, g, b uint32
	if c, ok := c.(color.RGBA); ok {
		r = uint32(c.R)
		g = uint32(c.G)
		b = uint32(c.B)
	} else {
		r, g, b, _ = c.RGBA()
		r >>= 8
		g >>= 8
		b >>= 8
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(r&^7 | g>>5)
	s[1] = uint8((g&^3)<<3 | b>>3)
}

func (p *RGB16) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.R&^7 | c.G>>5
	s[1] = (c.G&^3)<<3 | c.B>>3
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGB16) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGB16{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGB{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// ImmRGB16 is an immutable counterpart of RGB16.
type ImmRGB16 struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	Pix    string          // the image pixels
}

// NewImmRGB16 returns a new ImmRGB16 image with the given bounds and content.
func NewImmRGB16(r image.Rectangle, bits string) *ImmRGB16 {
	return &ImmRGB16{
		Rect:   r,
		Stride: 2 * r.Dx(),
		Pix:    bits,
	}
}

func (p *ImmRGB16) ColorModel() color.Model { return RGB16Model }
func (p *ImmRGB16) Bounds() image.Rectangle { return p.Rect }
func (p *ImmRGB16) Opaque() bool            { return true }

func (p *ImmRGB16) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *ImmRGB16) RGBAAt(x, y int) color.RGBA {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2]
	return rgb16torgba(s[0], s[1])
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ImmRGB16) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImmRGB16) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImmRGB16{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImmRGB16{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}
