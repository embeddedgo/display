// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

// Ellipse draws an empty or filled ellipse.
func (a *Area) Ellipse(p image.Point, ra, rb int, fill bool) {
	var r image.Rectangle
	if ra <= 0 || rb <= 0 {
		if ra >= 0 || rb >= 0 {
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
				if fill {
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
				} else {
					lminX := p.X + cx
					lmaxX := p.X + x
					rminX := p.X - x + 1
					rmaxX := p.X - cx + 1
					r.Min.Y = p.Y - y + 1
					if cy != 0 {
						r.Max.Y = p.Y - cy + 1
						r.Min.X = lminX
						r.Max.X = lmaxX
						a.Fill(r)
						r.Min.X = rminX
						r.Max.X = rmaxX
						a.Fill(r)
						r.Min.Y = p.Y + cy
					}
					r.Max.Y = p.Y + y
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
	r.Min.X = p.X + cx
	r.Max.X = p.X - cx + 1
	r.Min.Y = p.Y - rb
	r.Max.Y = p.Y - y + 1
	a.Fill(r)
	r.Min.Y = p.Y + y
	r.Max.Y = p.Y + rb + 1
	a.Fill(r)
}
