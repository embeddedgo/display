// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/draw"
)

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

// The simple implemetations of DrawPixel, hline and vline are almost as fast as
// the optimized ones in case of real hardware that is slow in terms of
// transferring commands and data. For example the simple implementation of
// DrawPixel below is only 1.0003 slower than the optimized one that uses
// Point.In and drv.Fill directly (ili9341.Driver, STM32 21 MHz DMA SPI) but
// takes 500 bytes less of Flash.

// DrawPixel provides a convenient way to fill one-pixel rectangle.
func (a *Area) DrawPixel(x, y int) {
	a.Fill(image.Rectangle{image.Pt(x, y), image.Pt(x+1, y+1)})
}

func hline(a *Area, x0, y0, x1 int) {
	if x0 > x1 {
		x1, x0 = x0, x1
	}
	a.Fill(image.Rectangle{image.Pt(x0, y0), image.Pt(x1+1, y0+1)})
}

func vline(a *Area, x0, y0, y1 int) {
	if y0 > y1 {
		y1, y0 = y0, y1
	}
	a.Fill(image.Rectangle{image.Pt(x0, y0), image.Pt(x0+1, y1+1)})
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
