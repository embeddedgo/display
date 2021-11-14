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
	end := false
	for {
		e2 := 2 * e
		h := e2 >= dy
		v := e2 <= dx
		px, py := x, y
		if h {
			if x == p1.X {
				end = true
			}
			e += dy
			x += sx
		}
		if v {
			if y == p1.Y {
				end = true
			}
			e += dx
			y += sy
		}
		if h && v || end {
			if py != p0.Y {
				vline(a, px, p0.Y, py)
			} else {
				hline(a, p0.X, py, px)
			}
			if end {
				break
			}
			p0.X = x
			p0.Y = y
		}
	}
}

/*
Some other considered implementations

// The original Alois Zingl algorithm
func (a *Area) Line(p0, p1 image.Point) {
	dx, sx := absSign(p1.X - p0.X)
	dy, sy := absSign(p1.Y - p0.Y)
	dy = -dy
	e := dx + dy
	for {
		a.Pixel(p0.X, p0.Y)
		e2 := 2 * e
		if e2 >= dy {
			if p0.X == p1.X {
				break
			}
			e += dy
			p0.X += sx
		}
		if e2 <= dx {
			if p0.Y == p1.Y {
				break
			}
			e += dx
			p0.Y += sy
		}
	}
}

// Avoid px, py
func (a *Area) Line(p0, p1 image.Point) {
	// based on Alois Zingl algorithm
	dx, sx := absSign(p1.X - p0.X)
	dy, sy := absSign(p1.Y - p0.Y)
	dy = -dy
	e := dx + dy
	x := p0.X
	y := p0.Y
	for {
		e2 := 2 * e
		h := e2 >= dy
		v := e2 <= dx
		if end := (h && x == p1.X) || (v && y == p1.Y); h && v || end {
			if y != p0.Y {
				vline(a, x, p0.Y, y)
			} else {
				hline(a, p0.X, y, x)
			}
			if end {
				break
			}
		}
		if h {
			e += dy
			x += sx
		}
		if v {
			e += dx
			y += sy
			if h {
				p0.X = x
				p0.Y = y
			}
		}
	}
}

*/
