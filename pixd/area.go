// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/color"
)

// Area is the drawing area on the display. It has its own coordinates with the
// (0, 0) origin regardless of its position on the display.
type Area struct {
	disp    *Display
	color   color.Color
	p0      image.Point
	visible image.Rectangle
	size    image.Point
}

func setColor(a *Area) {
	if a.disp.lastColor != a.color {
		a.disp.lastColor = a.color
		a.disp.drv.SetColor(a.color)
	}
}

// Rect returns the area coordinates on the display.
func (a *Area) Rect() image.Rectangle {
	return image.Rectangle{a.p0, a.p0.Add(a.size)}
}

// SetRect chandges the rectangle occupied by the area on the display.
func (a *Area) SetRect(r image.Rectangle) {
	a.p0 = r.Min
	a.size = r.Size()
	a.visible = r.Intersect(a.disp.Bounds())
}

// Bounds return the area bounds in its own coordinate system.
func (a *Area) Bounds() image.Rectangle {
	return image.Rectangle{Max: a.size}
}

// SetColor sets the color used by drawing methods.
func (a *Area) SetColor(c color.Color) {
	a.color = c
}

// SetColorRGBA is equivalent of SetColor(color.RGBA{r, g, b, alpha}). Notice
// that r, g, b must be alpha-premultiplied, they must be less than or equal to
// alpha.
func (a *Area) SetColorRGBA(r, g, b, alpha uint8) {
	a.color = color.RGBA{r, g, b, alpha}
}

// Color returns the color used by drawing methods.
func (a *Area) Color() color.Color {
	return a.color
}

// TextWriter returns a ready to use TextWriter initialized as below:
//	w := new(TextWriter)
//	w.Area = a
//	w.Face = f
//	w.Color = a.Color()
//	_, w.Pos.Y = f.Size() // ascent
func (a *Area) TextWriter(f FontFace) *TextWriter {
	_, ascent := f.Size()
	return &TextWriter{
		Area:  a,
		Face:  f,
		Color: a.color,
		Pos:   image.Point{0, ascent},
	}
}
