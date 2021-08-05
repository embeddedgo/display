// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tft

import (
	"image"
)

// DrawPoint draws a point (one pixel).
func (a *Area) DrawPoint(p image.Point) {
	if !p.In(a.Bounds()) {
		return
	}
	p = p.Add(a.P0())
	a.disp.fillRect(a.disp.dci, p.X, p.X, p.Y, p.Y, 1, a.color)
}

func (a *Area) rawFillRect(x0, y0, x1, y1, wxh int) {
	x0 += int(a.x0)
	y0 += int(a.y0)
	x1 += int(a.x0)
	y1 += int(a.y0)
	a.disp.fillRect(a.disp.dci, x0, y0, x1, y1, wxh, a.color)
}

// FillRect draws a filled rectangle.
func (a *Area) FillRect(r image.Rectangle) {
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		a.rawFillRect(r.Min.X, r.Min.Y, r.Max.X-1, r.Max.Y-1, r.Dx()*r.Dy())
	}
}

func (a *Area) hline(x0, y0, x1 int) {
	r := a.Bounds()
	if y0 < r.Min.Y || y0 >= r.Max.Y {
		return
	}
	if x0 < r.Min.X {
		x0 = r.Min.X
	}
	if x1 >= r.Max.X {
		x1 = r.Max.X - 1
	}
	if x0 <= x1 {
		a.rawFillRect(x0, y0, x1, y0, x1-x0+1)
	}
}

func (a *Area) vline(x0, y0, y1 int) {
	r := a.Bounds()
	if x0 < r.Min.X || x0 >= r.Max.X {
		return
	}
	if y0 < r.Min.Y {
		y0 = r.Min.Y
	}
	if y1 >= r.Max.Y {
		y1 = r.Max.Y - 1
	}
	if y0 <= y1 {
		a.rawFillRect(x0, y0, x0, y1, y1-y0+1)
	}
}
