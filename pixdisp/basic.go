// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

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
		rawFill(a, image.Rect(x0, y0, x1, y0+1))
	}
}

func vline(a *Area, x0, y0, y1 int) {
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
		rawFill(a, image.Rect(x0, y0, x0+1, y1))
	}
}

// DrawPixel is a faster counterpart of
// a.Fill(image.Rect(p.X, p.Y, p.X+1, p.Y+1)).
func (a *Area) DrawPixel(p image.Point) {
	setColor(a)
	drawPixel(a, p)
}

// Fill fills the given rectangle.
func (a *Area) Fill(r image.Rectangle) {
	setColor(a)
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		rawFill(a, r)
	}
}

// FillRect draws a filled rectangle for the given diagonal.
func (a *Area) FillRect(p0, p1 image.Point) {
	a.Fill(image.Rectangle{p0, p1.Add(image.Point{1, 1})})
}

// DrawRect draws a rectangle for the given diagonal.
//func (a *Area) DrawRect(p0, p1 image.Point) {
//	a.Fill(image.Rectangle{p0, p1.Add(image.Point{1, 1})})
//}

// Draw works like draw.Draw with dst set to the image representing the whole
// area.
func (a *Area) Draw(r image.Rectangle, src image.Image, sp image.Point, op draw.Op) {
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		a.disp.drv.Draw(r.Add(a.P0()), src, sp, nil, image.Point{}, op)
	}
}

// DrawMask works like draw.DrawMask with dst set to the image representing the
// whole area.
func (a *Area) DrawMask(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		a.disp.drv.Draw(r.Add(a.P0()), src, sp, mask, mp, op)
	}
}
