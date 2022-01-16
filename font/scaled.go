// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package font

import (
	"image"

	"github.com/embeddedgo/display/images"
)

// Scaling mode.
const (
	// Nearest means nearest-neighbor scaling. Use if the pixelated look of the
	// scaled font is acceptable or desirable.
	Nearest  = images.Nearest

	// Bilinear means bilinear interpolation. Use if the pixelated look of the
	// Nearest is undesirable.
	Bilinear = images.Bilinear
)

// Scaled can wrap Face to scale it up at runtime by integer factor. It can be
// useful if the available fonts are too small for a display. Scaled is UNSAFE
// for concurrent use.
type Scaled struct {
	fa Face
	si images.ScaledUp
}

// NewScaled wraps face to scale it up by mul factor usind given scaling mode.
func NewScaled(face Face, mul int, mode byte) *Scaled {
	p := new(Scaled)
	p.fa = face
	p.si.Mul = mul
	p.si.Mode = mode
	return p
}

// Size implements Face interface.
func (p *Scaled) Size() (height, ascent int) {
	height, ascent = p.fa.Size()
	return height * p.si.Mul, ascent * p.si.Mul
}

// Advance implements Face interface.
func (p *Scaled) Advance(r rune) int {
	return p.fa.Advance(r) * p.si.Mul
}

// Glyph implements Face interface.
func (p *Scaled) Glyph(r rune) (img image.Image, origin image.Point, advance int) {
	img, origin, advance = p.fa.Glyph(r)
	p.si.Image = img
	return &p.si, origin.Mul(p.si.Mul), advance * p.si.Mul

}