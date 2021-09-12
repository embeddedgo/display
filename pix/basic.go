// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/draw"
)

func drawPixel(a *Area, p image.Point) {
	if !p.In(a.Bounds()) {
		return
	}
	p = p.Add(a.P0())
	a.disp.drv.Fill(image.Rect(p.X, p.Y, p.X+1, p.Y+1))
}

func rawFill(a *Area, r image.Rectangle) {
	a.disp.drv.Fill(r.Add(a.P0()))
}

func hline(a *Area, x0, y0, x1 int) {
	r := a.Bounds()
	if y0 < r.Min.Y || y0 >= r.Max.Y {
		return
	}
	if x1 < x0 {
		x1, x0 = x0, x1
	}
	x1 += 1
	if x0 < r.Min.X {
		x0 = r.Min.X
	}
	if x1 >= r.Max.X {
		x1 = r.Max.X
	}
	if x0 <= x1 {
		rawFill(a, image.Rect(x0, y0, x1, y0+1))
	}
}

func vline(a *Area, x0, y0, y1 int) {
	r := a.Bounds()
	if x0 < r.Min.X || x0 >= r.Max.X {
		return
	}
	if y1 < y0 {
		y1, y0 = y0, y1
	}
	y1 += 1
	if y0 < r.Min.Y {
		y0 = r.Min.Y
	}
	if y1 >= r.Max.Y {
		y1 = r.Max.Y
	}
	if y0 <= y1 {
		rawFill(a, image.Rect(x0, y0, x0+1, y1))
	}
}

// Fill fills the given rectangle.
func (a *Area) Fill(r image.Rectangle) {
	setColor(a)
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		rawFill(a, r)
	}
}

// Draw works like draw.DrawMask with dst set to the image representing the
// whole area.
func (a *Area) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		a.disp.drv.Draw(r.Add(a.P0()), src, sp, mask, mp, op)
	}
}
