// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/draw"
)

// Fill fills the given rectangle.
func (a *Area) Fill(r image.Rectangle) {
	if a.misrc.Mode != MI {
		if a.misrc.Mode&MV != 0 {
			r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
			r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
		}
		if a.misrc.Mode&MX != 0 {
			r.Min.X, r.Max.X = -r.Max.X, -r.Min.X
		}
		if a.misrc.Mode&MY != 0 {
			r.Min.Y, r.Max.Y = -r.Max.Y, -r.Min.Y
		}
	}
	for ad := &a.ad; ad != nil; ad = ad.link {
		rt := r.Add(ad.tr).Intersect(ad.visible)
		if !rt.Empty() {
			ad.disp.mt.Lock()
			if ad.disp.lastColor != a.color {
				ad.disp.lastColor = a.color
				ad.disp.drv.SetColor(a.color)
			}
			ad.disp.drv.Fill(rt)
			ad.disp.mt.Unlock()
		}
	}
}

// The simple implemetations of Point, hline and vline are almost as fast as
// the optimized ones in case of real hardware that is slow in terms of
// transferring commands and data. For example the simple implementation of
// Point below is only 1.0003 slower than the optimized one that uses Point.In
// and drv.Fill directly (ili9341.Driver, STM32 21 MHz DMA SPI) but
// takes 500 bytes less of Flash.

// Point provides a convenient way to fill one-pixel rectangle.
func (a *Area) Point(x, y int) {
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
	mp = mp.Add(delta)
	if a.misrc.Mode != MI {
		if _, ok := src.(*image.Uniform); !ok {
			a.misrc.Image = src
			src = &a.misrc
		}
		if mask != nil {
			if _, ok := mask.(*image.Uniform); !ok {
				a.mimask.Image = mask
				mask = &a.mimask
			}
		}
		if a.misrc.Mode&MV != 0 {
			r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
			r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
			sp.X, sp.Y = sp.Y, sp.X
			mp.X, mp.Y = mp.Y, mp.X
		}
		if a.misrc.Mode&MX != 0 {
			r.Min.X, r.Max.X = -r.Max.X, -r.Min.X
			dx := r.Min.X - r.Max.X
			sp.X = dx - sp.X
			mp.X = dx - mp.X
		}
		if a.misrc.Mode&MY != 0 {
			r.Min.Y, r.Max.Y = -r.Max.Y, -r.Min.Y
			dy := r.Min.Y - r.Max.Y
			sp.Y = dy - sp.Y
			mp.Y = dy - mp.Y
		}
	}
	for ad := &a.ad; ad != nil; ad = ad.link {
		rt := r.Add(ad.tr)
		orig := rt.Min
		rt = rt.Intersect(ad.visible)
		if rt.Empty() {
			continue
		}
		delta := rt.Min.Sub(orig)
		spt := sp.Add(delta)
		var mpt image.Point
		if mask != nil {
			mpt = mp.Add(delta)
		}
		ad.disp.mt.Lock()
		ad.disp.drv.Draw(rt, src, spt, mask, mpt, op)
		ad.disp.mt.Unlock()
	}
}
