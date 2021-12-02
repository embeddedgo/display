// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"

	"github.com/embeddedgo/display/math2d"
)

func (a *Area) Arc(p image.Point, ra0, rb0, ra1, rb1 int, theta0, theta1 int32, fill bool) {
	arc(a, p, ra0, rb0, theta0, theta1)
}

func arc(a *Area, p image.Point, ra, rb int, theta0, theta1 int32) {
	const (
		frac  = 29
		round = 1 << (frac - 1)
	)
	one := image.Pt(1<<frac, 0)
	p0 := math2d.Rotate(one, theta0)
	p0.X = int((int64(ra)*int64(p0.X) + round) >> frac)
	p0.Y = int((int64(rb)*int64(p0.Y) + round) >> frac)
	p1 := math2d.Rotate(one, theta1)
	p1.X = int((int64(ra)*int64(p1.X) + round) >> frac)
	p1.Y = int((int64(rb)*int64(p1.Y) + round) >> frac)

}

/*
func (a *Area) DrawEllipse(p image.Point, ra, rb int) {
	setColor(a)
	// Alois Zingl algorithm
	x := -ra
	y := 0
	e2 := rb
	dx := (1 + 2*x) * e2 * e2
	dy := x * x
	err := dx + dy
	bb2 := 2 * rb * rb
	aa2 := 2 * ra * ra
	for x <= 0 {
		drawPixel(a, p.Add(image.Point{-x, y}))
		drawPixel(a, p.Add(image.Point{x, y}))
		drawPixel(a, p.Add(image.Point{x, -y}))
		drawPixel(a, p.Add(image.Point{-x, -y}))
		e2 = 2 * err
		if e2 >= dx {
			x++
			dx += bb2
			err += dx
		}
		if e2 <= dy {
			y++
			dy += aa2
			err += dy
		}
	}
	for y < rb {
		y++
		drawPixel(a, p.Add(image.Point{0, y}))
		drawPixel(a, p.Add(image.Point{0, -y}))
	}
}
*/