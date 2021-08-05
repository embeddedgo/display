// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
)

// DrawPixel is a faster counterpart of
// a.FillRect(image.Rect(p.X, p.Y, p.X+1, p.Y+1)).
func (a *Area) DrawPixel(p image.Point) {
	if !p.In(a.Bounds()) {
		return
	}
	p = p.Add(a.P0())
	a.disp.drv.Fill(image.Rect(p.X, p.Y, p.X+1, p.Y+1), a.color)
}

func (a *Area) rawFillRect(r image.Rectangle) {
	a.disp.drv.Fill(r.Add(a.P0()), a.color)
}

// FillRect draws a filled rectangle.
func (a *Area) FillRect(r image.Rectangle) {
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		a.rawFillRect(r)
	}
}

func (a *Area) hline(x0, y0, x1 int) {
	x1 += 1
	r := a.Bounds()
	if y0 < r.Min.Y || y0 >= r.Max.Y {
		return
	}
	if x0 < r.Min.X {
		x0 = r.Min.X
	}
	if x1 >= r.Max.X {
		x1 = r.Max.X
	}
	if x0 <= x1 {
		a.rawFillRect(image.Rect(x0, y0, x1, y0+1))
	}
}

func (a *Area) vline(x0, y0, y1 int) {
	y1 += 1
	r := a.Bounds()
	if x0 < r.Min.X || x0 >= r.Max.X {
		return
	}
	if y0 < r.Min.Y {
		y0 = r.Min.Y
	}
	if y1 >= r.Max.Y {
		y1 = r.Max.Y
	}
	if y0 <= y1 {
		a.rawFillRect(image.Rect(x0, y0, x0+1, y1))
	}
}
