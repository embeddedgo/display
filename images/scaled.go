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

// ScaledUp can wrap image to scale it up at runtime by integer factor.
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
	if p.Mode == Nearest {
		return p.Image.At(x/p.Mul, y/p.Mul)
	}
	return p.Image.At(x/p.Mul, y/p.Mul)
}
