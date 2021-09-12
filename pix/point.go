// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import "image"

// DrawPoint draws a point with a given radius.
func (a *Area) DrawPoint(p image.Point, r int) {
	setColor(a)
	if r <= 0 {
		if r == 0 {
			drawPixel(a, p)
		}
		return
	}
	// based on Alois Zingl algorithm
	// fill the four sides of the circle
	x, y, e := -r, 0, 2*(1-r)
	for x+y < 0 {
		ny := y + 1
		ne := e + 2*ny + 1
		if e > x || ne > ny {
			x0, x1 := p.X-y, p.X+y
			hline(a, x0, p.Y-x, x1)
			hline(a, x0, p.Y+x, x1)
			y0, y1 := p.Y-y, p.Y+y
			vline(a, p.X+x, y0, y1)
			vline(a, p.X-x, y0, y1)
			x += 1
			ne += 2*x + 1
		}
		y = ny
		e = ne
	}
	// fill the center rectangle
	rect := image.Rectangle{
		p.Sub(image.Point{x - 1, x - 1}),
		p.Add(image.Point{x, x}),
	}
	rect = rect.Canon().Intersect(a.Bounds())
	if !rect.Empty() {
		rawFill(a, rect)
	}
}
