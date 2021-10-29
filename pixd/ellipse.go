// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"time"
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
	drawPixel(a, p.Add(image.Pt(-x, 0)))
	drawPixel(a, p.Add(image.Pt(x, 0)))
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
		time.Sleep(500 * time.Millisecond)
		drawPixel(a, p.Add(image.Pt(-x, y)))
		drawPixel(a, p.Add(image.Pt(x, y)))
		drawPixel(a, p.Add(image.Pt(x, -y)))
		drawPixel(a, p.Add(image.Pt(-x, -y)))
	}
	for y < rb {
		y++
		drawPixel(a, p.Add(image.Pt(0, y)))
		drawPixel(a, p.Add(image.Pt(0, -y)))
	}
}

// FillEllipse draws a filled ellipse.
func (a *Area) FillEllipse(p image.Point, ra, rb int) {
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
	hline(a, p.X-x, p.Y, p.X+x)
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
			hline(a, p.X-x, p.Y-y, p.X+x)
			hline(a, p.X-x, p.Y+y, p.X+x)
		}
	}
	for y < rb {
		y++
		drawPixel(a, p.Add(image.Pt(0, y)))
		drawPixel(a, p.Add(image.Pt(0, -y)))
	}

}
