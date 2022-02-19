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

// Mirror can be used to wrap an image to reflect it through three different
// axes: (0, t) for MX, (t, 0) for MY and (t, t) for MV (parentheses contain
// parametric descriptions of axes).
type Mirror struct {
	Image image.Image
	Mode  int
}

func NewMirror(img image.Image, mvxy int) *Mirror {
	return &Mirror{img, mvxy}
}

// ColorModel implements image.Image interface.
func (p *Mirror) ColorModel() color.Model {
	return p.Image.ColorModel()
}

// Bounds implements image.Image interface.
func (p *Mirror) Bounds() image.Rectangle {
	r := p.Image.Bounds()
	if p.Mode&MV != 0 {
		r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
		r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
	}
	if p.Mode&MX != 0 {
		r.Min.X, r.Max.X = -r.Max.X, -r.Min.X
	}
	if p.Mode&MY != 0 {
		r.Min.Y, r.Max.Y = -r.Max.Y, -r.Min.Y
	}
	return r
}

// At implements image.Image interface.
func (p *Mirror) At(x, y int) color.Color {
	if p.Mode&MX != 0 {
		x = -1 - x
	}
	if p.Mode&MY != 0 {
		y = -1 - y
	}
	if p.Mode&MV != 0 {
		x, y = y, x
	}
	return p.Image.At(x, y)
}

// RGBA64At implements image.RGBA64Image interface.
func (p *Mirror) RGBA64At(x, y int) color.RGBA64 {
	if p.Mode&MX != 0 {
		x = -1 - x
	}
	if p.Mode&MY != 0 {
		y = -1 - y
	}
	if p.Mode&MV != 0 {
		x, y = y, x
	}
	if img, ok := p.Image.(RGBA64Image); ok {
		return img.RGBA64At(x, y)
	}
	r, g, b, a := p.Image.At(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

/*
Reflecing relative to the center of the coordinate system (as seen above) has
many nice properties and overall looks like the only correct way to reflect
Go images.

Reflecting relative to the center of the image does not change the image
position/bounds in case of Mode&MV == 0 but it must change it anyway for MV
reflection. It also  has slower At method, mainly because of calling Bounds
method (see code below).

func (p *Mirror) Bounds() image.Rectangle {
	r := p.Image.Bounds()
	if p.Mode&MV != 0 {
		r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
		r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
	}
	return r
}

func (p *Mirror) At(x, y int) color.Color {
	r := p.Image.Bounds()
	if p.Mode&MX != 0 {
		x = r.Min.X + r.Max.X - 1 - x
	}
	if p.Mode&MY != 0 {
		y = r.Min.Y + r.Max.Y - 1 - y
	}
	if p.Mode&MV != 0 {
		x, y = y, x
	}
	return p.Image.At(x, y)
}
*/
