// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package images

import (
	"image"
	"image/color"
)

const (
	MI = 0         // identity (no operation)
	MV = 1 << iota // swap X with Y
	MX             // mirror X axis
	MY             // mirror Y axis
)

type Mirrored struct {
	Image image.Image
	Mode  int
}

func NewMirrored(img image.Image, mvxy int) *Mirrored {
	return &Mirrored{img, mvxy}
}

// ColorModel implements image.Image interface.
func (p *Mirrored) ColorModel() color.Model {
	return p.Image.ColorModel()
}

// Bounds implements image.Image interface.
func (p *Mirrored) Bounds() image.Rectangle {
	r := p.Image.Bounds()
	if p.Mode&MV != 0 {
		r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
		r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
	}
	return r
}

// At implements image.Image interface.
func (p *Mirrored) At(x, y int) color.Color {
	r := p.Image.Bounds()
	if p.Mode&MX != 0 {
		x = r.Max.X - 1 - x
	}
	if p.Mode&MY != 0 {
		y = r.Max.Y - 1 - y
	}
	if p.Mode&MV != 0 {
		x, y = y, x
	}
	return p.Image.At(x, y)
}
