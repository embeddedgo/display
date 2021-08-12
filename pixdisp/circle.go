// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import "image"

// The following code uses algorithm described by Alois Zingl in
// "A Rasterizing Algorithm for Drawing Curves".

// DrawCircle draws an empty circle.
func (a *Area) DrawCircle(p image.Point, r int) {
	if r <= 0 {
		if r == 0 {
			a.DrawPixel(p)
		}
		return
	}
	x, y, e := -r, 0, 2*(1-r)
	for x+y <= 0 {
		a.DrawPixel(p.Add(image.Pt(-x, y)))
		a.DrawPixel(p.Add(image.Pt(x, y)))
		a.DrawPixel(p.Add(image.Pt(-x, -y)))
		a.DrawPixel(p.Add(image.Pt(x, -y)))
		if x+y != 0 {
			a.DrawPixel(p.Add(image.Pt(-y, x)))
			a.DrawPixel(p.Add(image.Pt(y, x)))
			a.DrawPixel(p.Add(image.Pt(-y, -x)))
			a.DrawPixel(p.Add(image.Pt(y, -x)))
		}
		y += 1
		ne := e + 2*y + 1
		if e > x || ne > y {
			x += 1
			ne += 2*x + 1
		}
		e = ne
	}
}

// FillCircle draws a filled circle.
func (a *Area) FillCircle(p image.Point, r int) {
	if r <= 0 {
		if r == 0 {
			a.DrawPixel(p)
		}
		return
	}
	// Fill the four sides.
	x, y, e := -r, 0, 2*(1-r)
	for x+y < 0 {
		ny := y + 1
		ne := e + 2*ny + 1
		if e > x || ne > ny {
			x0, x1 := p.X-y, p.X+y
			a.hline(x0, p.Y-x, x1)
			a.hline(x0, p.Y+x, x1)
			y0, y1 := p.Y-y, p.Y+y
			a.vline(p.X+x, y0, y1)
			a.vline(p.X-x, y0, y1)
			x += 1
			ne += 2*x + 1
		}
		y = ny
		e = ne
	}
	// Fill the center.
	a.FillRect(image.Rectangle{
		p.Add(image.Pt(x, x)), p.Sub(image.Pt(x-1, x-1)),
	})
}
