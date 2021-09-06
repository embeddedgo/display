// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import "image"

// The following code uses algorithm described by Alois Zingl in
// "A Rasterizing Algorithm for Drawing Curves".

// DrawCircle draws an empty circle.
func (a *Area) DrawCircle(p image.Point, r int) {
	setColor(a)
	if r <= 0 {
		if r == 0 {
			drawPixel(a, p)
		}
		return
	}
	x, y, e := -r, 0, 2*(1-r)
	for x+y <= 0 {
		drawPixel(a, p.Add(image.Point{-x, y}))
		drawPixel(a, p.Add(image.Point{x, y}))
		drawPixel(a, p.Add(image.Point{-x, -y}))
		drawPixel(a, p.Add(image.Point{x, -y}))
		if x+y != 0 {
			drawPixel(a, p.Add(image.Point{-y, x}))
			drawPixel(a, p.Add(image.Point{y, x}))
			drawPixel(a, p.Add(image.Point{-y, -x}))
			drawPixel(a, p.Add(image.Point{y, -x}))
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
