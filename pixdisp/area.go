// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
)

type Area struct {
	disp   *Display
	color  color.Color
	rect   image.Rectangle
	p0     image.Point
	width  int
	height int
	swapWH bool
}

func setColor(a *Area) {
	if a.disp.lastColor != a.color {
		a.disp.lastColor = a.color
		a.disp.drv.SetColor(a.color)
	}
}

func (a *Area) P0() image.Point {
	return a.p0
}

func (a *Area) updateBounds() {
	wh := a.rect.Intersect(a.disp.Bounds())
	a.p0 = wh.Min
	a.width = wh.Dx()
	a.height = wh.Dy()
	a.swapWH = a.disp.swapWH
}

func (a *Area) Rect() image.Rectangle {
	return a.rect
}

func (a *Area) SetRect(r image.Rectangle) {
	a.rect = r.Canon()
	a.updateBounds()
}

func (a *Area) Bounds() image.Rectangle {
	if a.swapWH != a.disp.swapWH {
		a.updateBounds()
	}
	return image.Rectangle{Max: image.Point{int(a.width), int(a.height)}}
}

// SetColor sets the color used by drawing methods.
func (a *Area) SetColor(c color.Color) {
	a.color = c
}

// SetColorRGB is a convenient wrapper over SetColor(RGB{r, g, b}).
func (a *Area) SetColorRGB(r, g, b uint8) {
	a.color = RGB{r, g, b}
}

func (a *Area) TextWriter(f FontFace) *TextWriter {
	_, ascent := f.Size()
	return &TextWriter{
		Area:  a,
		Face:  f,
		Color: a.color,
		Pos:   image.Point{0, ascent},
	}
}
