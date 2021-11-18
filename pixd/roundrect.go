// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"


// RoundRect can draw rectangles, ellipses and circles. If ra = rb = 0 it draws
// an empty or filled rectangle with a given diagonal (p0, p1). Non-zero ra
// enlarges the rectangle horizontally, non-zero rb enlarges it vertically. The
// enlarged rectangle has rounded corners with eliptic shape. If p0 = p1
// RoundRect draws an ellipse with a given minor and major radius.
func (a *Area) RoundRect(p0, p1 image.Point, ra, rb int, fill bool) {
	var r image.Rectangle
	if p0.X > p1.X {
		p1.X, p0.X = p0.X, p1.X
	}
	if p0.Y > p1.Y {
		p1.Y, p0.Y = p0.Y, p1.Y
	}
	if ra <= 0 || rb <= 0 {
		p0.X -= ra
		p1.X += ra
		p0.Y -= rb
		p1.Y += rb
		r.Min = p0
		if fill {
			r.Max.X = p1.X + 1
			r.Max.Y = p1.Y + 1
		} else {
			r.Max.X = p1.X + 1
			r.Max.Y = p0.Y + 1
			a.Fill(r)
			r.Max.X = p0.X + 1
			r.Max.Y = p1.Y + 1
			a.Fill(r)
			r.Max.X = p1.X + 1
			r.Min.Y = p1.Y
			a.Fill(r)
			r.Min.X = p1.X
			r.Min.Y = p0.Y
		}
		a.Fill(r)
		return
	}
	// based on Alois Zingl ellipse algorithm
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
				if fill {
					r.Min.X = p0.X + cx
					r.Max.X = p1.X - cx + 1
					r.Min.Y = p0.Y - y + 1
					if cy != 0 {
						r.Max.Y = p0.Y - cy + 1
						a.Fill(r)
						r.Min.Y = p1.Y + cy
					}
					r.Max.Y = p1.Y + y
					a.Fill(r)
				} else {
					lminX := p0.X + cx
					lmaxX := p0.X + x
					rminX := p1.X - x + 1
					rmaxX := p1.X - cx + 1
					r.Min.Y = p0.Y - y + 1
					if cy != 0 {
						r.Max.Y = p0.Y - cy + 1
						r.Min.X = lminX
						r.Max.X = lmaxX
						a.Fill(r)
						r.Min.X = rminX
						r.Max.X = rmaxX
						a.Fill(r)
						r.Min.Y = p1.Y + cy
					}
					r.Max.Y = p1.Y + y
					r.Min.X = lminX
					r.Max.X = lmaxX
					a.Fill(r)
					r.Min.X = rminX
					r.Max.X = rmaxX
					a.Fill(r)
				}
				cx = x
				cy = y
			}
			dy += aa2
			err += dy
		}
	}
	r.Min.X = p0.X + cx
	r.Max.X = p1.X - cx + 1
	r.Min.Y = p0.Y - rb
	r.Max.Y = p0.Y - y + 1
	a.Fill(r)
	r.Min.Y = p1.Y + y
	r.Max.Y = p1.Y + rb + 1
	a.Fill(r)
}
