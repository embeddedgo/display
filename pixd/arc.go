// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"

	"github.com/embeddedgo/display/math2d"
)

const frac = 18 // maxa, maxb up to 8191 without 64-bit multiplication

func mulfi(p image.Point, mx, my int) image.Point {
	const round = 1 << (frac - 1)
	p.X = (p.X*mx + round) >> frac
	p.Y = (p.Y*my + round) >> frac
	return p
}

func (a *Area) Arc(p image.Point, mina, minb, maxa, maxb int, th0, th1 int32, fill bool) {
	// bounding box
	var box image.Rectangle
	box.Min.X = -maxa
	box.Max.X = maxa
	box.Min.Y = -maxb
	box.Max.Y = maxb
	if a0, a1 := uint32(th0), uint32(th1); a0 < a1 && a1 <= math2d.FullAngle/2 {
		box.Min.Y = 0
	} else if th0 < th1 && th1 <= 0 {
		box.Max.Y = 0
	}
	tr0 := th0 + math2d.RightAngle
	tr1 := th1 + math2d.RightAngle
	if a0, a1 := uint32(tr0), uint32(tr1); a0 < a1 && a1 <= math2d.FullAngle/2 {
		box.Min.X = 0
	} else if tr0 < tr1 && tr1 <= 0 {
		box.Max.X = 0
	}
	box.Max.X++
	box.Max.Y++
	box = box.Intersect(a.Bounds().Sub(p))
	// setup ellipses
	minaa := mina * mina
	minbb := minb * minb
	maxaa := maxa * maxa
	maxbb := maxb * maxb
	xx := box.Min.X * box.Min.X
	yy := box.Min.Y * box.Min.Y
	wmin := (xx*minbb + yy*minaa) - minaa*minbb
	wmax := maxaa*maxbb - (xx*maxbb + yy*maxaa)
	minaa2 := minaa * 2
	minbb2 := minbb * 2
	maxaa2 := maxaa * 2
	maxbb2 := maxbb * 2
	dwymin := minaa2*box.Min.Y + minaa
	dwymax := maxaa2*box.Min.Y + maxaa
	dwxmin := minbb2*box.Min.X + minbb
	dwxmax := maxbb2*box.Min.X + maxbb
	// setup two sides
	one := image.Pt(1<<frac, 0)
	cosSin := math2d.Rotate(one, th0)
	pmin0 := mulfi(cosSin, mina, minb)
	pmax0 := mulfi(cosSin, maxa, maxb)
	cosSin = math2d.Rotate(one, th1)
	pmin1 := mulfi(cosSin, mina, minb)
	pmax1 := mulfi(cosSin, maxa, maxb)
	dx0 := pmax0.X - pmin0.X
	dx1 := pmin1.X - pmax1.X
	dy0 := pmin0.Y - pmax0.Y
	dy1 := pmax1.Y - pmin1.Y
	w0 := dx0*(box.Min.Y-pmin0.Y) + dy0*(box.Min.X-pmin0.X)
	w1 := dx1*(box.Min.Y-pmin1.Y) + dy1*(box.Min.X-pmin1.X)
	of := int(th1 - th0)
	if of < 0 {
		dy0, dy1 = -dy0, -dy1
		dx0, dx1 = -dx0, -dx1
		w0, w1 = -w0, -w1
	}
	// fill
	for y := box.Min.Y; y < box.Max.Y; y++ {
		m0, m1 := w0, w1
		mmin, mmax := wmin, wmax
		dmmin, dmmax := dwxmin, dwxmax
		for x := box.Min.X; x < box.Max.X; x++ {
			if (m0|m1)^of|(mmin|mmax) >= 0 {
				a.Pixel(p.X+x, p.Y+y)
			}
			m0 += dy0
			m1 += dy1
			mmin += dmmin
			mmax -= dmmax
			dmmin += minbb2
			dmmax += maxbb2
		}
		w0 += dx0
		w1 += dx1
		wmin += dwymin
		wmax -= dwymax
		dwymin += minaa2
		dwymax += maxaa2
	}
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
