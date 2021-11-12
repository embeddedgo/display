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

func (a *Area) Line(p0, p1 image.Point) {
	// based on Alois Zingl algorithm
	dx, sx := absSign(p1.X - p0.X)
	dy, sy := absSign(p1.Y - p0.Y)
	dy = -dy
	e := dx + dy
	x := p0.X
	y := p0.Y
	var h, v bool
	for {
		e2 := 2 * e
		h = e2 >= dy
		v = e2 <= dx
		if h {
			if v {
				if y != p0.Y {
					vline(a, x, p0.Y, y)
				} else {
					hline(a, p0.X, y, x)
				}
			}
			if x == p1.X {
				break
			}
			e += dy
			x += sx
		}
		if v {
			if y == p1.Y {
				break
			}
			e += dx
			y += sy
			if h {
				p0.X = x
				p0.Y = y
			}
		}
	}
	if h != v {
		if y != p0.Y {
			vline(a, x, p0.Y, y)
		} else {
			hline(a, p0.X, y, x)
		}
	}
}
