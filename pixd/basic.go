// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/draw"
)

func drawPixel(a *Area, p image.Point) {
	for ad := &a.ad; ad != nil; ad = ad.link {
		tp := p.Add(ad.tr)
		if tp.In(ad.visible) {
			setColor(a, ad.disp)
			ad.disp.drv.Fill(image.Rect(tp.X, tp.Y, tp.X+1, tp.Y+1))
		}
	}
}

/*
func hline(a *Area, x0, y0, x1 int) {
	var r image.Rectangle
	r.Min.X = x0
	r.Min.Y = y0
	r.Max.X = x1 + 1
	r.Max.Y = y0 + 1
	a.Fill(r)
}

func vline(a *Area, x0, y0, y1 int) {
	var r image.Rectangle
	r.Min.X = x0
	r.Min.Y = y0
	r.Max.X = x0 + 1
	r.Max.Y = y1 + 1
	a.Fill(r)
}
*/

func hline(a *Area, x0, y0, x1 int) {
	for ad := &a.ad; ad != nil; ad = ad.link {
		v := ad.visible
		ty0 := y0 + ad.tr.Y
		if ty0 < v.Min.Y || v.Max.Y <= ty0 {
			continue
		}
		tx0 := x0 + ad.tr.X
		tx1 := x1 + ad.tr.X
		if tx1 < tx0 {
			tx1, tx0 = tx0, tx1
		}
		tx1 += 1
		if tx0 < v.Min.X {
			tx0 = v.Min.X
		}
		if tx1 >= v.Max.X {
			tx1 = v.Max.X
		}
		if tx0 <= tx1 {
			r := image.Rectangle{
				image.Point{tx0, ty0},
				image.Point{tx1, ty0 + 1},
			}
			setColor(a, ad.disp)
			ad.disp.drv.Fill(r)
		}
	}
}

func vline(a *Area, x0, y0, y1 int) {
	for ad := &a.ad; ad != nil; ad = ad.link {
		v := ad.visible
		tx0 := x0 + ad.tr.X
		if tx0 < v.Min.X || v.Max.X <= tx0 {
			continue
		}
		ty0 := y0 + ad.tr.Y
		ty1 := y1 + ad.tr.Y
		if ty1 < ty0 {
			ty1, ty0 = ty0, ty1
		}
		ty1 += 1
		if ty0 < v.Min.Y {
			ty0 = v.Min.Y
		}
		if ty1 >= v.Max.Y {
			ty1 = v.Max.Y
		}
		if ty0 <= ty1 {
			r := image.Rectangle{
				image.Point{tx0, ty0},
				image.Point{tx0 + 1, ty1},
			}
			setColor(a, ad.disp)
			ad.disp.drv.Fill(r)
		}
	}
}

// Fill fills the given rectangle.
func (a *Area) Fill(r image.Rectangle) {
	for ad := &a.ad; ad != nil; ad = ad.link {
		tr := r.Add(ad.tr).Intersect(ad.visible)
		if !tr.Empty() {
			setColor(a, ad.disp)
			ad.disp.drv.Fill(tr)
		}
	}
}

// Draw works like draw.DrawMask with dst set to the image representing the
// whole area.
func (a *Area) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	orig := r.Min
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
	for ad := &a.ad; ad != nil; ad = ad.link {
		trt := r.Add(ad.tr)
		orig := trt.Min
		trt = trt.Intersect(ad.visible)
		if trt.Empty() {
			continue
		}
		delta := trt.Min.Sub(orig)
		tsp := sp.Add(delta)
		var tmp image.Point
		if mask != nil {
			tmp = mp.Add(delta)
		}
		ad.disp.drv.Draw(trt, src, tsp, mask, tmp, op)
	}
}
