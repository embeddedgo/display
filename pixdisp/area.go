// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
)

type Area struct {
	color  uint64
	disp   *Display
	rect   image.Rectangle
	p0     image.Point
	width  int
	height int
	swapWH bool
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

func (a *Area) Bounds() image.Rectangle {
	if a.swapWH != a.disp.swapWH {
		a.updateBounds()
	}
	return image.Rectangle{Max: image.Pt(int(a.width), int(a.height))}
}

// SetColor sets the color used by drawing methods.
func (a *Area) SetColor(c color.Color) {
	a.color = a.disp.drv.Color(c)
}

// SetColorRGB is a convenient wrapper over SetColor(RGB{r, g, b}).
func (a *Area) SetColorRGB(r, g, b uint8) {
	a.color = a.disp.drv.Color(color.RGBA{r, g, b, 255})
}
