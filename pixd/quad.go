// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
)

// Quad draws a quadrilateral or a triangle described by the given vertices. If
// fill is true the vertices must describe a convex quadrilateral or triangle.
// To draw triangle make two adjacent vertices identical. The fill parameter is
// provided for symmetry with the other drawing functions but for empty
// quadrilaterals or triangles the drawing performance can be worese than using
// Line method directly.
func (a *Area) Quad(p0, p1, p2, p3 image.Point, fill bool) {
	if fill {
		fillQuad(a, p0, p1, p2, p3, false)
	} else {
		pixel := true
		if p0.X != p1.X || p0.Y != p1.Y {
			a.Line(p0, p1)
			pixel = false
		}
		if p1.X != p2.X || p1.Y != p2.Y {
			a.Line(p1, p2)
			pixel = false
		}
		if p2.X != p3.X || p2.Y != p3.Y {
			a.Line(p2, p3)
			pixel = false
		}
		if p3.X != p0.X || p3.Y != p0.Y {
			a.Line(p3, p0)
			pixel = false
		}
		if pixel {
			a.Pixel(p0.X, p0.Y)
		}
	}
}

// FillQuad draws a filled convex quadrilateral or a triangle described by given
// vertices. It obeys the filling rules so can be used to draw filled polygons
// composed of adjacent quadrilaterals/triangles ensuring that the common edges
// are drawn once. To draw triangle make two adjacent vertices identical.
func (a *Area) FillQuad(p0, p1, p2, p3 image.Point) {
	fillQuad(a, p0, p1, p2, p3, true)
}

func min4(a, b, c, d int) int {
	if a > b {
		a = b
	}
	if a > c {
		a = c
	}
	if a > d {
		a = d
	}
	return a
}

func max4(a, b, c, d int) int {
	if a < b {
		a = b
	}
	if a < c {
		a = c
	}
	if a < d {
		a = d
	}
	return a
}

// notTopLeft reports whether the directed edge p0->p1 of a clockwise triangle
// is not the top and not the left one.
func notTopLeft(p0, p1 image.Point) bool {
	return p0.Y < p1.Y || p0.Y == p1.Y && p0.X > p1.X
}

// orient returns >0 if p lies on the right side of the oriented line defined by
// vector p0->p1, <0 if p lies on the left side of the line, =0 if p lies
// on the line (p0, p1, p are colinear).
func orient(p0, p1, p image.Point) int {
	return (p1.X-p0.X)*(p.Y-p0.Y) - (p1.Y-p0.Y)*(p.X-p0.X)
}

func fillQuad(a *Area, p0, p1, p2, p3 image.Point, fillRules bool) {
	// make the quad clockwise
	var q1, q2 image.Point
	if p1.Eq(p0) {
		q1 = p2
		q2 = p3
	} else {
		q1 = p1
		if p2.Eq(p1) {
			q2 = p3
		} else {
			q2 = p2
		}
	}
	if orient(p0, q1, q2) < 0 {
		p0, p2 = p2, p0
	}
	// bounding box
	var box image.Rectangle
	box.Min.X = min4(p0.X, p1.X, p2.X, p3.X)
	box.Min.Y = min4(p0.Y, p1.Y, p2.Y, p3.Y)
	box.Max.X = max4(p0.X, p1.X, p2.X, p3.X) + 1
	box.Max.Y = max4(p0.Y, p1.Y, p2.Y, p3.Y) + 1
	box = box.Intersect(a.Bounds())
	// setup
	a01 := p0.Y - p1.Y
	a12 := p1.Y - p2.Y
	a23 := p2.Y - p3.Y
	a30 := p3.Y - p0.Y
	b01 := p1.X - p0.X
	b12 := p2.X - p1.X
	b23 := p3.X - p2.X
	b30 := p0.X - p3.X
	w0 := orient(p1, p2, box.Min)
	w1 := orient(p2, p3, box.Min)
	w2 := orient(p3, p0, box.Min)
	w3 := orient(p0, p1, box.Min)
	if fillRules {
		if notTopLeft(p1, p2) {
			w0--
		}
		if notTopLeft(p2, p3) {
			w1--
		}
		if notTopLeft(p3, p0) {
			w2--
		}
		if notTopLeft(p0, p1) {
			w3--
		}
	}
	// fill
	var r image.Rectangle
	if box.Dx() >= box.Dy() {
		for r.Min.Y = box.Min.Y; r.Min.Y < box.Max.Y; r.Min.Y++ {
			m0, m1, m2, m3 := w0, w1, w2, w3
			r.Min.X = box.Max.X
			r.Max.X = box.Min.X
			for r.Max.X < box.Max.X {
				m := m0 | m1 | m2 | m3
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
				m1 += a23
				m2 += a30
				m3 += a01
				r.Max.X++
			}
			r.Max.Y = r.Min.Y + 1
			a.Fill(r)
			w0 += b12
			w1 += b23
			w2 += b30
			w3 += b01
		}
	} else {
		for r.Min.X = box.Min.X; r.Min.X < box.Max.X; r.Min.X++ {
			m0, m1, m2, m3 := w0, w1, w2, w3
			r.Min.Y = box.Max.Y
			r.Max.Y = box.Min.Y
			for r.Max.Y < box.Max.Y {
				m := m0 | m1 | m2 | m3
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
				m1 += b23
				m2 += b30
				m3 += b01
				r.Max.Y++
			}
			r.Max.X = r.Min.X + 1
			a.Fill(r)
			w0 += a12
			w1 += a23
			w2 += a30
			w3 += a01
		}
	}
}
