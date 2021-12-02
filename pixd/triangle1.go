// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

func (a *Area) Triangle(p0, p1, p2 image.Point, fill bool) {
	if fill {
		fillTriangle(a, p0, p1, p2, false)
	} else {
		a.Line(p0, p1)
		a.Line(p1, p2)
		a.Line(p2, p0)
	}
}

func (a *Area) FillTriangle(p0, p1, p2 image.Point) {
	fillTriangle(a, p0, p1, p2, true)
}

func min3(a, b, c int) int {
	if a > b {
		a = b
	}
	if a > c {
		a = c
	}
	return a
}

func max3(a, b, c int) int {
	if a < b {
		a = b
	}
	if a < c {
		a = c
	}
	return a
}

// orient returns >0 if the triangle described by p0, p1, p2 is clockwise,
// <0 if it is counter clockwise, =0 if p0, p1, p2 are collinear.
func orient(p0, p1, p2 image.Point) int {
	return (p1.X-p0.X)*(p2.Y-p0.Y) - (p1.Y-p0.Y)*(p2.X-p0.X)
}

// notTopLeft reports whether the directed edge p0-p1 of a clockwise triangle
// is not the top and not the left one.
func notTopLeft(p0, p1 image.Point) bool {
	return p0.Y < p1.Y || p0.Y == p1.Y && p0.X > p1.X
}

// For description of the algorithm see:
// https://fgiesen.wordpress.com/2013/02/08/triangle-rasterization-in-practice/

func fillTriangle(a *Area, p0, p1, p2 image.Point, fillRules bool) {
	if orient(p0, p1, p2) < 0 {
		p1, p2 = p2, p1 // make the triangle clockwise oriented
	}
	// bounding box
	var box image.Rectangle
	box.Min.X = min3(p0.X, p1.X, p2.X)
	box.Min.Y = min3(p0.Y, p1.Y, p2.Y)
	box.Max.X = max3(p0.X, p1.X, p2.X) + 1
	box.Max.Y = max3(p0.Y, p1.Y, p2.Y) + 1
	box = box.Intersect(a.Bounds())
	// setup
	a01 := p0.Y - p1.Y
	a12 := p1.Y - p2.Y
	a20 := p2.Y - p0.Y
	b01 := p1.X - p0.X
	b12 := p2.X - p1.X
	b20 := p0.X - p2.X
	w0 := orient(p1, p2, box.Min)
	w1 := orient(p2, p0, box.Min)
	w2 := orient(p0, p1, box.Min)
	if fillRules {
		if notTopLeft(p1, p2) {
			w0--
		}
		if notTopLeft(p2, p0) {
			w1--
		}
		if notTopLeft(p0, p1) {
			w2--
		}
	}
	// fill
	var r image.Rectangle
	if box.Dx() >= box.Dy() {
		for r.Min.Y = box.Min.Y; r.Min.Y < box.Max.Y; r.Min.Y++ {
			m0, m1, m2 := w0, w1, w2
			r.Min.X = box.Max.X
			r.Max.X = box.Min.X
			for r.Max.X < box.Max.X {
				m := m0 | m1 | m2
				if r.Min.X == box.Max.X {
					if m >= 0 {
						r.Min.X = r.Max.X
					}
				} else {
					if m < 0 {
						break
					}
				}
				m0 += a12
				m1 += a20
				m2 += a01
				r.Max.X++
			}
			r.Max.Y = r.Min.Y + 1
			a.Fill(r)
			w0 += b12
			w1 += b20
			w2 += b01
		}
	} else {
		for r.Min.X = box.Min.X; r.Min.X < box.Max.X; r.Min.X++ {
			m0, m1, m2 := w0, w1, w2
			r.Min.Y = box.Max.Y
			r.Max.Y = box.Min.Y
			for r.Max.Y < box.Max.Y {
				m := m0 | m1 | m2
				if r.Min.Y == box.Max.Y {
					if m >= 0 {
						r.Min.Y = r.Max.Y
					}
				} else {
					if m < 0 {
						break
					}
				}
				m0 += b12
				m1 += b20
				m2 += b01
				r.Max.Y++
			}
			r.Max.X = r.Min.X + 1
			a.Fill(r)
			w0 += a12
			w1 += a20
			w2 += a01
		}
	}
}

/*
	// select drawing orientation
	rot := box.Dx() < box.Dy()
	if rot {
		box.Min.X, box.Min.Y = box.Min.Y, box.Min.X
		box.Max.X, box.Max.Y = box.Max.Y, box.Max.X
		a01, b01 = b01, a01
		a12, b12 = b12, a12
		a20, b20 = b20, a20
	}
	// fill
	for y := box.Min.Y; y < box.Max.Y; y++ {
		m0, m1, m2 := w0, w1, w2
		minX := box.Max.X
		maxX := box.Min.X
		for maxX < box.Max.X {
			m := m0 | m1 | m2
			if minX == box.Max.X {
				if m >= 0 {
					minX = maxX
				}
			} else {
				if m < 0 {
					break
				}
			}
			m0 += a12
			m1 += a20
			m2 += a01
			maxX++
		}
		var r image.Rectangle
		if rot {
			r.Min.Y = minX
			r.Max.Y = maxX
			r.Min.X = y
			r.Max.X = y + 1
		} else {
			r.Min.X = minX
			r.Max.X = maxX
			r.Min.Y = y
			r.Max.Y = y + 1
		}
		a.Fill(r)
		w0 += b12
		w1 += b20
		w2 += b01
	}
*/
