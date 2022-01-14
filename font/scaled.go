// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package font

import (
	"image"
	"image/color"
)

// Scaling mode.
const (
	// Nearest means nearest-neighbor scaling. Use if the pixelated look of
	// returned glyphs is acceptable or desirable.
	Nearest byte = iota

	// Bilinear means bilinear interpolation. Use if the pixelated look of
	// Nearest is undesirable.
	Bilinear
)

type scaled struct {
	img  image.Image
	mul  int
	mode byte
}

func (p *scaled) ColorModel() color.Model {
	return p.img.ColorModel()
}

func (p *scaled) Bounds() image.Rectangle {
	r := p.img.Bounds()
	r.Min = r.Min.Mul(p.mul)
	r.Max = r.Max.Mul(p.mul)
	return r
}

func (p *scaled) At(x, y int) color.Color {
	if p.mode == Nearest {
		return p.img.At(x/p.mul, y/p.mul)
	}
	return p.img.At(x/p.mul, y/p.mul)
}

// Scaled can wrap Face to scale it up at runtime by integer factor. It can be
// useful if the available fonts are too small for a display. Scaled is UNSAFE
// for concurrent use.
type Scaled struct {
	fa Face
	si scaled
}

// NewScaled wraps face to scale it up by mul factor usind given scaling mode.
func NewScaled(face Face, mul int, mode byte) *Scaled {
	p := new(Scaled)
	p.fa = face
	p.si.mul = mul
	p.si.mode = mode
	return p
}

// Size implements Face interface.
func (p *Scaled) Size() (height, ascent int) {
	height, ascent = p.fa.Size()
	return height * p.si.mul, ascent * p.si.mul
}

// Advance implements Face interface.
func (p *Scaled) Advance(r rune) int {
	return p.fa.Advance(r) * p.si.mul
}

// Glyph implements Face interface.
func (p *Scaled) Glyph(r rune) (img image.Image, origin image.Point, advance int) {
	img, origin, advance = p.fa.Glyph(r)
	p.si.img = img
	return &p.si, origin.Mul(p.si.mul), advance * p.si.mul

}
