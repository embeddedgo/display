// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
)

// DrawCircle draws empty circle with the center in p. Notice that p is the
// center of symmetry of the image.Rect(p.X-r, p.Y-r, p.X+r, p.Y+r). If r <= 0
// then DrawCircle returns wthout drawing anything.
func (a *Area) DrawCircle(p image.Point, r int) {
	if r <= 0 {
		return
	}
	x, y := 0, r-1
	rr := r*r*4 - r
	yy := (2*y + 1)
	yy *= yy
	for {
		a.DrawPixel(p.Add(image.Pt(x, y)))
		a.DrawPixel(p.Add(image.Pt(-x-1, y)))
		a.DrawPixel(p.Add(image.Pt(x, -y-1)))
		a.DrawPixel(p.Add(image.Pt(-x-1, -y-1)))
		a.DrawPixel(p.Add(image.Pt(y, x)))
		a.DrawPixel(p.Add(image.Pt(-y-1, x)))
		a.DrawPixel(p.Add(image.Pt(y, -x-1)))
		a.DrawPixel(p.Add(image.Pt(-y-1, -x-1)))
		x += 1
		if x >= y {
			break
		}
		xx := (2*x + 1)
		xx *= xx
		if xx+yy >= rr {
			y -= 1
			yy = (2*y + 1)
			yy *= yy
		}
	}
}

func (a *Area) FillCircle(p image.Point, r int) {
	if r <= 0 {
		return
	}
	x, y := 0, r-1
	rr := r*r*4 - r
	yy := (2*y + 1)
	yy *= yy
	for {
		//y0, y1 := p.Y-y-1, p.Y+y
		x0, x1 := p.X-x-1, p.X+x
		a.hline(x0, p.Y-y-1, x1)
		a.hline(x0, p.Y+y, x1)
		/*
			a.DrawPixel(p.Add(image.Pt(x, y)))
			a.DrawPixel(p.Add(image.Pt(-x-1, y)))
			a.DrawPixel(p.Add(image.Pt(x, -y-1)))
			a.DrawPixel(p.Add(image.Pt(-x-1, -y-1)))
			a.DrawPixel(p.Add(image.Pt(y, x)))
			a.DrawPixel(p.Add(image.Pt(-y-1, x)))
			a.DrawPixel(p.Add(image.Pt(y, -x-1)))
			a.DrawPixel(p.Add(image.Pt(-y-1, -x-1)))
		*/
		x += 1
		if x >= y {
			break
		}
		xx := (2*x + 1)
		xx *= xx
		if xx+yy >= rr {
			y -= 1
			yy = (2*y + 1)
			yy *= yy
		}
	}
	// Fill central rectangle.
	//a.FillRect(image.Rectangle{
	//	p.Sub(image.Pt(x, y)), p.Add(image.Pt(x+1, y+1)),
	//})
}

// DrawCircle1 draws an empty circle around the pixel defined by the
// image.Rect(p.X, p.Y, p.X+1, p.Y+1). Notice that the circle center is the
// center of this pixel. If r < 0 then DrawCircle1 returns wthout drawing
// anything. If r == 0 it is equivalent to DrawPixel(p).
func (a *Area) DrawCircle1(p image.Point, r int) {
	if r <= 0 {
		if r == 0 {
			a.DrawPixel(p)
		}
		return
	}
	x, y, e := r, 0, 1-r
	for x >= y {
		a.DrawPixel(p.Add(image.Pt(-x, y)))
		a.DrawPixel(p.Add(image.Pt(x, y)))
		a.DrawPixel(p.Add(image.Pt(-x, -y)))
		a.DrawPixel(p.Add(image.Pt(x, -y)))
		a.DrawPixel(p.Add(image.Pt(-y, x)))
		a.DrawPixel(p.Add(image.Pt(y, x)))
		a.DrawPixel(p.Add(image.Pt(-y, -x)))
		a.DrawPixel(p.Add(image.Pt(y, -x)))
		y++
		e += 2*y + 1
		if e > 0 {
			x--
			e -= 2 * x
		}
	}
}

// FillCircle1 is like DrawCircle1 but draws a filled circle.
func (a *Area) FillCircle1(p image.Point, r int) {
	if r <= 0 {
		if r == 0 {
			a.DrawPixel(p)
		}
		return
	}
	// Fill four sides.
	x, y, e := r, 0, 1-r
	for x > y {
		e += 2*y + 3
		if e > 0 {
			y0, y1 := p.Y-y, p.Y+y
			a.vline(p.X+x, y0, y1)
			a.vline(p.X-x, y0, y1)
			x0, x1 := p.X-y, p.X+y
			a.hline(x0, p.Y-x, x1)
			a.hline(x0, p.Y+x, x1)
			x--
			e -= 2 * x
		}
		y++
	}
	// Fill central rectangle.
	a.FillRect(image.Rectangle{
		p.Sub(image.Pt(x, y)), p.Add(image.Pt(x+1, y+1)),
	})
}
