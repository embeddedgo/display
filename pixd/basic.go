// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/draw"
)

func drawPixel(a *Area, p image.Point) {
	p = p.Add(a.p0)
	if !p.In(a.visible) {
		return
	}
	a.disp.drv.Fill(image.Rect(p.X, p.Y, p.X+1, p.Y+1))
}

func hline(a *Area, x0, y0, x1 int) {
	x0 += a.p0.X
	x1 += a.p0.X
	y0 += a.p0.Y
	v := a.visible
	if y0 < v.Min.Y || y0 >= v.Max.Y {
		return
	}
	if x1 < x0 {
		x1, x0 = x0, x1
	}
	x1 += 1
	if x0 < v.Min.X {
		x0 = v.Min.X
	}
	if x1 >= v.Max.X {
		x1 = v.Max.X
	}
	if x0 <= x1 {
		r := image.Rectangle{image.Point{x0, y0}, image.Point{x1, y0 + 1}}
		a.disp.drv.Fill(r)
	}
}

func vline(a *Area, x0, y0, y1 int) {
	x0 += a.p0.X
	y0 += a.p0.Y
	y1 += a.p0.Y
	v := a.visible
	if x0 < v.Min.X || x0 >= v.Max.X {
		return
	}
	if y1 < y0 {
		y1, y0 = y0, y1
	}
	y1 += 1
	if y0 < v.Min.Y {
		y0 = v.Min.Y
	}
	if y1 >= v.Max.Y {
		y1 = v.Max.Y
	}
	if y0 <= y1 {
		r := image.Rectangle{image.Point{x0, y0}, image.Point{x0 + 1, y1}}
		a.disp.drv.Fill(r)
	}
}

// Fill fills the given rectangle.
func (a *Area) Fill(r image.Rectangle) {
	setColor(a)
	r = r.Add(a.p0).Intersect(a.visible)
	if !r.Empty() {
		a.disp.drv.Fill(r)
	}
}

// Draw works like draw.DrawMask with dst set to the image representing the
// whole area.
func (a *Area) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	r = r.Add(a.p0)
	orig := r.Min
	r = r.Intersect(a.visible)
	r = r.Intersect(src.Bounds().Add(orig.Sub(sp)))
	if mask != nil {
		r = r.Intersect(mask.Bounds().Add(orig.Sub(mp)))
	}
	if r.Empty() {
		return
	}
	delta := r.Min.Sub(orig)
	sp = sp.Add(delta)
	if mask != nil {
		mp = mp.Add(delta)
	}
	a.disp.drv.Draw(r, src, sp, mask, mp, op)
}
