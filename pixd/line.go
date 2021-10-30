// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

func absSign(x int) (abs, sign int) {
	abs = x
	sign = 1
	if abs < 0 {
		abs = -abs
		sign = -sign
	}
	return
}

func line(a *Area, x0, y0, x1, y1 int) {
	// based on Alois Zingl algorithm
	dx, sx := absSign(x1 - x0)
	dy, sy := absSign(y1 - y0)
	dy = -dy
	e := dx + dy
	x := x0
	y := y0
	for {
		e2 := 2 * e
		h := e2 >= dy
		v := e2 <= dx
		if h {
			if x == x1 {
				break
			}
			if v {
				if y != y0 {
					vline(a, x, y0, y)
				} else {
					hline(a, x0, y, x)
				}
			}
			e += dy
			x += sx
		}
		if v {
			if y == y1 {
				break
			}
			e += dx
			y += sy
			if h {
				y0 = y
				x0 = x
			}
		}
	}
	if y != y0 {
		vline(a, x, y0, y)
	} else {
		hline(a, x0, y, x)
	}
}

// Line connects the given points by drawing segments of straight line.
func (a *Area) Line(points ...image.Point) {
	if len(points) <= 1 {
		if len(points) == 1 {
			a.Pixel(points[0].X, points[0].Y)
		}
		return
	}
	p0 := points[0]
	for _, p1 := range points[1:] {
		line(a, p0.X, p0.Y, p1.X, p1.Y)
		p0 = p1
	}
}
