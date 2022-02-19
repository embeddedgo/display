// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package images

import (
	"image"
	"image/color"
)

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

func rgb16torgba64(h, l uint8) color.RGBA64 {
	r := uint(h&^7) << 8
	g := (uint(h)<<13 | uint(l)<<5) & 0xfc00
	b := (uint(l) << 11) & 0xf800
	return color.RGBA64{
		uint16(r | r>>5 | r>>10 | r>>15),
		uint16(g | g>>6 | g>>12),
		uint16(b | b>>5 | b>>10 | b>>15),
		0xffff,
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

func (p *RGB) RGBA64At(x, y int) color.RGBA64 {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	r := uint16(s[0])
	g := uint16(s[1])
	b := uint16(s[2])
	return color.RGBA64{r | r<<8, g | g<<8, b | b<<8, 255}
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

func (p *RGB) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(c.R >> 8)
	s[1] = uint8(c.G >> 8)
	s[2] = uint8(c.B >> 8)
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
	s := p.Pix[i : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	return color.RGBA{s[0], s[1], s[2], 255}
}

func (p *ImmRGB) RGBA64At(x, y int) color.RGBA64 {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	r := uint16(s[0])
	g := uint16(s[1])
	b := uint16(s[2])
	return color.RGBA64{r | r<<8, g | g<<8, b | b<<8, 255}
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

func (p *RGB16) RGBA64At(x, y int) color.RGBA64 {
	if !(image.Pt(x, y).In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	return rgb16torgba64(s[0], s[1])
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

func (p *RGB16) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Pt(x, y).In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+2 : i+2] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8((c.R>>8)&^7 | c.G>>13)
	s[1] = uint8((c.G&0xfc00)>>5 | c.B>>11)
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
