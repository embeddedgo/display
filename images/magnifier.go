// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package images

import (
	"image"
	"image/color"
)

// Scaling mode.
const (
	// Nearest means nearest-neighbor scaling. Use if the pixelated look of the
	// scaled image is acceptable or desirable.
	Nearest byte = iota

	// Bilinear means bilinear interpolation. Use if the pixelated look of the
	// Nearest is undesirable.
	Bilinear
)

// Magnifier can be used to wrap an image to scale it up at runtime by integer
// factor.
type Magnifier struct {
	Image  image.Image
	Sx, Sy int  // scaling factors along X and Y axes
	Mode   byte // scaling mode: Nearest or Bilinear
}

// Magnify wraps img into Magnifier to scale it up by scale factor usind given
// scaling mode.
func Magnify(img image.Image, sx, sy int, mode byte) *Magnifier {
	return &Magnifier{img, sx, sy, mode}
}

// ColorModel implements image.Image interface.
func (p *Magnifier) ColorModel() color.Model {
	return p.Image.ColorModel()
}

// Bounds implements image.Image interface.
func (p *Magnifier) Bounds() image.Rectangle {
	r := p.Image.Bounds()
	r.Min.X *= p.Sx
	r.Min.Y *= p.Sy
	r.Max.X *= p.Sx
	r.Max.Y *= p.Sy
	return r
}

// At implements image.Image interface.
func (p *Magnifier) At(x, y int) color.Color {
	if p.Mode != Nearest {
		x -= p.Sx / 2
		y -= p.Sy / 2
	}
	x0 := x
	if x0 < 0 {
		x0 -= p.Sx - 1 // make division round down even for negative x0
	}
	y0 := y
	if y0 < 0 {
		y0 -= p.Sy - 1 // make division round down even for negative y0
	}
	x0 /= p.Sx
	y0 /= p.Sy
	if p.Mode == Nearest {
		return p.Image.At(x0, y0)
	}
	return magnify(p, x, y, x0, y0)
}

// RGBA64At implements image.RGBA64Image interface.
func (p *Magnifier) RGBA64At(x, y int) color.RGBA64 {
	if p.Mode != Nearest {
		x -= p.Sx / 2
		y -= p.Sy / 2
	}
	x0 := x
	if x0 < 0 {
		x0 -= p.Sx - 1 // make division round down even for negative x0
	}
	y0 := y
	if y0 < 0 {
		y0 -= p.Sy - 1 // make division round down even for negative y0
	}
	x0 /= p.Sx
	y0 /= p.Sy
	if p.Mode == Nearest {
		if img, ok := p.Image.(image.RGBA64Image); ok {
			return img.RGBA64At(x0, y0)
		}
		r, g, b, a := p.Image.At(x0, y0).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	}
	return magnify(p, x, y, x0, y0)
}

func magnify(p *Magnifier, x, y, x0, y0 int) color.RGBA64 {
	x1 := x0 + 1
	y1 := y0 + 1
	var (
		r00, g00, b00, a00 uint32
		r10, g10, b10, a10 uint32
		r01, g01, b01, a01 uint32
		r11, g11, b11, a11 uint32
	)
	if img, ok := p.Image.(image.RGBA64Image); ok {
		c := img.RGBA64At(x0, y0)
		r00, g00, b00, a00 = uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
		c = img.RGBA64At(x1, y0)
		r10, g10, b10, a10 = uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
		c = img.RGBA64At(x0, y1)
		r01, g01, b01, a01 = uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
		c = img.RGBA64At(x1, y1)
		r11, g11, b11, a11 = uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	} else {
		r00, g00, b00, a00 = p.Image.At(x0, y0).RGBA()
		r10, g10, b10, a10 = p.Image.At(x1, y0).RGBA()
		r01, g01, b01, a01 = p.Image.At(x0, y1).RGBA()
		r11, g11, b11, a11 = p.Image.At(x1, y1).RGBA()
	}
	x0 *= p.Sx
	x1 *= p.Sx
	y0 *= p.Sy
	y1 *= p.Sy
	dx0 := uint(x - x0)
	dx1 := uint(x1 - x)
	dy0 := uint(y - y0)
	dy1 := uint(y1 - y)
	r0 := uint(r00)*dx1 + uint(r10)*dx0
	g0 := uint(g00)*dx1 + uint(g10)*dx0
	b0 := uint(b00)*dx1 + uint(b10)*dx0
	a0 := uint(a00)*dx1 + uint(a10)*dx0
	r1 := uint(r01)*dx1 + uint(r11)*dx0
	g1 := uint(g01)*dx1 + uint(g11)*dx0
	b1 := uint(b01)*dx1 + uint(b11)*dx0
	a1 := uint(a01)*dx1 + uint(a11)*dx0
	div := uint(p.Sx) * uint(p.Sy)
	r := (r0*dy1 + r1*dy0) / div
	g := (g0*dy1 + g1*dy0) / div
	b := (b0*dy1 + b1*dy0) / div
	a := (a0*dy1 + a1*dy0) / div
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}

}
