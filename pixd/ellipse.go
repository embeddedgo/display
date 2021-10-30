// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
)

// DrawEllipse draws an empty ellipse.
func (a *Area) DrawEllipse(p image.Point, ra, rb int) {
	if ra < 0 || rb < 0 {
		return
	}
	// Alois Zingl algorithm
	x := -ra
	y := 0
	e2 := rb
	dx := (1 + 2*x) * e2 * e2
	dy := x * x
	err := dx + dy
	bb2 := 2 * rb * rb
	aa2 := 2 * ra * ra
	a.DrawPixel(p.X-x, p.Y)
	a.DrawPixel(p.X+x, p.Y)
	for {
		e2 = 2 * err
		if e2 >= dx {
			if x++; x > 0 {
				break
			}
			dx += bb2
			err += dx
		}
		if e2 <= dy {
			y++
			dy += aa2
			err += dy
		}
		a.DrawPixel(p.X-x, p.Y+y)
		a.DrawPixel(p.X+x, p.Y+y)
		a.DrawPixel(p.X+x, p.Y-y)
		a.DrawPixel(p.X-x, p.Y-y)
	}
	for y < rb {
		y++
		a.DrawPixel(p.X, p.Y+y)
		a.DrawPixel(p.X, p.Y-y)
	}
}

// FillEllipse draws a filled ellipse.
func (a *Area) FillEllipse(p image.Point, ra, rb int) {
	if ra <= 0 || rb <= 0 {
		if ra >= 0 || rb >= 0 {
			var r image.Rectangle
			r.Min.X = p.X - ra
			r.Max.X = p.X + ra + 1
			r.Min.Y = p.Y - rb
			r.Max.Y = p.Y + rb + 1
			a.Fill(r)
		}
		return
	}
	// based on Alois Zingl algorithm
	x := -ra
	y := 0
	cx := x
	cy := y
	e2 := rb
	dx := (1 + 2*x) * e2 * e2
	dy := x * x
	err := dx + dy
	aa2 := 2 * ra * ra
	bb2 := 2 * rb * rb
	for {
		e2 = 2 * err
		if e2 >= dx {
			if x++; x > 0 {
				break
			}
			dx += bb2
			err += dx
		}
		if e2 <= dy {
			y++
			if x != cx {
				// x, y both changed
				var r image.Rectangle
				r.Min.X = p.X + cx
				r.Max.X = p.X - cx + 1
				r.Min.Y = p.Y - y + 1
				if cy != 0 {
					r.Max.Y = p.Y - cy + 1
					a.Fill(r)
					r.Min.Y = p.Y + cy
				}
				r.Max.Y = p.Y + y
				a.Fill(r)
				cx = x
				cy = y
			}
			dy += aa2
			err += dy
		}
	}
	var r image.Rectangle
	r.Min.X = p.X + cx
	r.Max.X = p.X - cx + 1
	r.Min.Y = p.Y - rb
	r.Max.Y = p.Y - y + 1
	a.Fill(r)
	r.Min.Y = p.Y + y
	r.Max.Y = p.Y + rb + 1
	a.Fill(r)
}
