// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

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
	if _, ok := c.(RGB); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

func rgb16Model(c color.Color) color.Color {
	if _, ok := c.(RGB16); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	r >>= 11
	g >>= 10
	b >>= 11
	return RGB16(r<<11 | g<<5 | b)
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

// ImageAlphaN is an in-memory image whose At method returns color.Alpha with
// 1, 2, 4 or 8 bit precision.
type ImageAlphaN struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	LogN   uint8           // 1<<LogN is the number of bits per pixel
	Shift  uint8           // the bit offest in Pix[0] to the first pixel
	Pix    []uint8         // the image pixels
}

// NewImageAlphaN returns a new ImageAlphaN image with the given bounds and
// number of bits per pixel.
func NewImageAlphaN(r image.Rectangle, nbpp int) *ImageAlphaN {
	p := new(ImageAlphaN)
	p.Rect = r
	p.LogN, p.Stride = logNStride(r, nbpp)
	p.Pix = make([]uint8, p.Stride*r.Dy())
	return p
}

func (p *ImageAlphaN) ColorModel() color.Model {
	return AlphaNModel(1 << p.LogN)
}

func (p *ImageAlphaN) Bounds() image.Rectangle {
	return p.Rect
}

func (p *ImageAlphaN) AlphaAt(x, y int) color.Alpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Alpha{}
	}
	i, s := p.PixOffset(x, y)
	return alphaLogN(uint(p.Pix[i])>>s, p.LogN)
}

func (p *ImageAlphaN) At(x, y int) color.Color {
	return p.AlphaAt(x, y)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y) and the index to the bits in that element that
// determines the pixel value.
func (p *ImageAlphaN) PixOffset(x, y int) (offset int, shift uint) {
	x += int(p.Shift)>>p.LogN - p.Rect.Min.X
	y -= p.Rect.Min.Y
	cs := 3 - p.LogN
	col := x >> cs
	offset = y*p.Stride + col
	shift = uint(x-col<<cs) << p.LogN
	return
}

func (p *ImageAlphaN) Set(x, y int, c color.Color) {
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

func (p *ImageAlphaN) SetAlpha(x, y int, c color.Alpha) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	rshift := uint(8) - 1<<p.LogN
	i, lshift := p.PixOffset(x, y)
	p.Pix[i] = p.Pix[i]&^(0xff>>rshift<<lshift) | c.A>>rshift<<lshift
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImageAlphaN) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be
	// inside either r1 or r2 if the intersection is empty. Without explicitly
	// checking for this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImageAlphaN{}
	}
	i, shift := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImageAlphaN{
		Rect:   r,
		LogN:   p.LogN,
		Shift:  uint8(shift),
		Stride: p.Stride,
		Pix:    p.Pix[i:],
	}
}

// ImmAlphaN is an immutable counterpart of ImageAlphaN.
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
	if !(image.Point{x, y}.In(p.Rect)) {
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

// ImageRGB is an in-memory image whose At method returns RGB values.
type ImageRGB struct {
	Rect   image.Rectangle // image bounds
	Stride int             // stride (in bytes) between vertically adjacent pixels
	Pix    []uint8         // the image pixels
}

// NewImageRGB returns a new ImageRGB image with the given bounds.
func NewImageRGB(r image.Rectangle) *ImageRGB {
	return &ImageRGB{
		Rect:   r,
		Stride: 3 * r.Dx(),
		Pix:    make([]uint8, 3*r.Dx()*r.Dy()),
	}
}

func (p *ImageRGB) ColorModel() color.Model { return RGBModel }
func (p *ImageRGB) Bounds() image.Rectangle { return p.Rect }
func (p *ImageRGB) Opaque() bool            { return true }

func (p *ImageRGB) At(x, y int) color.Color {
	return p.RGBAt(x, y)
}

func (p *ImageRGB) RGBAt(x, y int) RGB {
	if !(image.Point{x, y}.In(p.Rect)) {
		return RGB{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	return RGB{s[0], s[1], s[2]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ImageRGB) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *ImageRGB) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := RGBModel.Convert(c).(RGB)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
}

func (p *ImageRGB) SetRGB(x, y int, c RGB) {
	if !(image.Point{x, y}.In(p.Rect)) {
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
func (p *ImageRGB) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImageRGB{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImageRGB{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// ImmRGB is an immutable counterpart of ImageRGB.
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
	return p.RGBAt(x, y)
}

func (p *ImmRGB) RGBAt(x, y int) RGB {
	if !(image.Point{x, y}.In(p.Rect)) {
		return RGB{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	return RGB{s[0], s[1], s[2]}
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

// ImageRGB16 is an in-memory image whose At method returns RGB16 values.
type ImageRGB16 struct {
	Rect   image.Rectangle // image bounds
	Stride int             // Pix stride (in bytes) between vertically adjacent pixels
	Pix    []uint8         // the image pixels
}

// NewImageRGB16 returns a new ImageRGB16 image with the given bounds.
func NewImageRGB16(r image.Rectangle) *ImageRGB16 {
	return &ImageRGB16{
		Rect:   r,
		Stride: 2 * r.Dx(),
		Pix:    make([]uint8, 2*r.Dx()*r.Dy()),
	}
}

func (p *ImageRGB16) ColorModel() color.Model { return RGB16Model }
func (p *ImageRGB16) Bounds() image.Rectangle { return p.Rect }
func (p *ImageRGB16) Opaque() bool            { return true }

func (p *ImageRGB16) At(x, y int) color.Color {
	return p.RGB16At(x, y)
}

func (p *ImageRGB16) RGB16At(x, y int) RGB16 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return 0
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	return RGB16(uint(s[0])<<8 | uint(s[1]))
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *ImageRGB16) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
}

func (p *ImageRGB16) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := RGB16Model.Convert(c).(RGB16)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(c1 >> 8)
	s[1] = uint8(c1)
}

func (p *ImageRGB16) SetRGB16(x, y int, c RGB16) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(c >> 8)
	s[1] = uint8(c)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *ImageRGB16) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &ImageRGB16{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ImageRGB{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// ImmRGB16 is an immutable counterpart of ImageRGB16.
type ImmRGB16 struct {
	Rect   image.Rectangle // image bounds
	Stride int             // Pix stride (in bytes) between vertically adjacent pixels
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
	return p.RGB16At(x, y)
}

func (p *ImmRGB16) RGB16At(x, y int) RGB16 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return 0
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	return RGB16(uint(s[0])<<8 | uint(s[1]))
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
