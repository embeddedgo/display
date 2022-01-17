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

// ScaledUp can be used to wrap image to scale it up at runtime by integer
// factor.
type ScaledUp struct {
	Image image.Image
	Mul   int
	Mode  byte
}

// NewScaledUp wraps img to scale it up by mul factor usind given scaling mode.
func NewScaledUp(img image.Image, mul int, mode byte) *ScaledUp {
	return &ScaledUp{img, mul, mode}
}

// ColorModel implements image.Image interface.
func (p *ScaledUp) ColorModel() color.Model {
	return p.Image.ColorModel()
}

// Bounds implements image.Image interface.
func (p *ScaledUp) Bounds() image.Rectangle {
	r := p.Image.Bounds()
	r.Min = r.Min.Mul(p.Mul)
	r.Max = r.Max.Mul(p.Mul)
	return r
}

// At implements image.Image interface.
func (p *ScaledUp) At(x, y int) color.Color {
	if p.Mode != Nearest {
		center := p.Mul / 2
		x -= center
		y -= center
	}
	round := p.Mul - 1
	x0 := x
	if x0 < 0 {
		x0 -= round // make division round down even for negative x0
	}
	y0 := y
	if y0 < 0 {
		y0 -= round // make division round down even for negative y0
	}
	x0 /= p.Mul
	y0 /= p.Mul
	if p.Mode == Nearest {
		return p.Image.At(x0, y0)
	}
	x1 := x0 + 1
	y1 := y0 + 1
	r00, g00, b00, a00 := p.Image.At(x0, y0).RGBA()
	r10, g10, b10, a10 := p.Image.At(x1, y0).RGBA()
	r01, g01, b01, a01 := p.Image.At(x0, y1).RGBA()
	r11, g11, b11, a11 := p.Image.At(x1, y1).RGBA()
	x0 *= p.Mul
	x1 *= p.Mul
	y0 *= p.Mul
	y1 *= p.Mul
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
	div := uint(p.Mul) * uint(p.Mul)
	r := (r0*dy1 + r1*dy0) / div
	g := (g0*dy1 + g1*dy0) / div
	b := (b0*dy1 + b1*dy0) / div
	a := (a0*dy1 + a1*dy0) / div
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}
