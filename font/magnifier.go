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
	Nearest = images.Nearest

	// Bilinear means bilinear interpolation. Use if the pixelated look of the
	// Nearest is undesirable.
	Bilinear = images.Bilinear
)

// Magnifier can wrap Face to scale it up at runtime by integer factor. It can
// be useful if the available fonts are too small for a display. Magnifier is
// UNSAFE for concurrent use. A goroutine should magnify the selected font face
// just before using it and should not share magnified faces with other
// goroutines.
type Magnifier struct {
	fa Face
	mi images.Magnifier
}

// Magnify wraps face into Magnifier to scale it up by scale factor usind given
// scaling mode.
func Magnify(face Face, scale int, mode byte) *Magnifier {
	p := new(Magnifier)
	p.fa = face
	p.mi.Scale = scale
	p.mi.Mode = mode
	return p
}

// Face returns the base font face.
func (p *Magnifier) Face() Face {
	return p.fa
}

// SetFace sets the base font face.
func (p *Magnifier) SetFace(face Face) {
	p.fa = face
}

// Scale returns current scaling factor.
func (p *Magnifier) Scale() int {
	return p.mi.Scale
}

// SetScale sets the scaling factor.
func (p *Magnifier) SetScale(scale int) {
	p.mi.Scale = scale
}

// Size implements Face interface.
func (p *Magnifier) Size() (height, ascent int) {
	height, ascent = p.fa.Size()
	return height * p.mi.Scale, ascent * p.mi.Scale
}

// Advance implements Face interface.
func (p *Magnifier) Advance(r rune) int {
	return p.fa.Advance(r) * p.mi.Scale
}

// Glyph implements Face interface.
func (p *Magnifier) Glyph(r rune) (img image.Image, origin image.Point, advance int) {
	img, origin, advance = p.fa.Glyph(r)
	p.mi.Image = img
	return &p.mi, origin.Mul(p.mi.Scale), advance * p.mi.Scale

}
