// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

// Quad draws a quadrilateral or triangle described by the given vertices. If
// fill is true the quadrilateral must be convex. Quad draws a triangle if two
// adjacent vertices are identical.
func (a *Area) Quad(p0, p1, p2, p3 image.Point, fill bool) {
	if fill {
		fillQuad(a, p0, p1, p2, p3, false)
	} else {
		a.Line(p0, p1)
		if p1.X != p2.X || p1.Y != p2.Y {
			a.Line(p1, p2)
		}
		if p2.X != p3.X || p2.Y != p3.Y {
			a.Line(p2, p3)
		}
		if p3.X != p0.X || p3.Y != p0.Y {
			a.Line(p3, p0)
		}
	}
}

// FillQuad draws a filled convex quadrilateral or triangle described by the
// given vertices. It obeys the filling rules so can be used to draw filled
// polygons composed of adjacent quadrilaterals/triangles ensuring that the
// common edges are drawn once. FillQuad draws a triangle if two adjacent
// vertices are identical.
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
	if (q1.X-p0.X)*(q2.Y-p0.Y) < (q1.Y-p0.Y)*(q2.X-p0.X) {
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
	dx01 := p1.X - p0.X
	dx12 := p2.X - p1.X
	dx23 := p3.X - p2.X
	dx30 := p0.X - p3.X
	dy01 := p0.Y - p1.Y
	dy12 := p1.Y - p2.Y
	dy23 := p2.Y - p3.Y
	dy30 := p3.Y - p0.Y
	w0 := dx12*(box.Min.Y-p1.Y) + dy12*(box.Min.X-p1.X)
	w1 := dx23*(box.Min.Y-p2.Y) + dy23*(box.Min.X-p2.X)
	w2 := dx30*(box.Min.Y-p3.Y) + dy30*(box.Min.X-p3.X)
	w3 := dx01*(box.Min.Y-p0.Y) + dy01*(box.Min.X-p0.X)
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
				m0 += dy12
				m1 += dy23
				m2 += dy30
				m3 += dy01
				r.Max.X++
			}
			r.Max.Y = r.Min.Y + 1
			a.Fill(r)
			w0 += dx12
			w1 += dx23
			w2 += dx30
			w3 += dx01
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
				m0 += dx12
				m1 += dx23
				m2 += dx30
				m3 += dx01
				r.Max.Y++
			}
			r.Max.X = r.Min.X + 1
			a.Fill(r)
			w0 += dy12
			w1 += dy23
			w2 += dy30
			w3 += dy01
		}
	}
}

// IsConvex reports whether the given vertices describe a convex polygon.
func IsConvex(vertices ...image.Point) bool {
	if len(vertices) == 0 {
		return false
	}
	lastxp := 0
	a := vertices[len(vertices)-1]
	b := vertices[0]
	for _, c := range vertices[1:] {
		xp := (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
		a, b = b, c
		if xp == 0 {
			continue
		}
		if lastxp == 0 {
			lastxp = xp
			continue
		}
		if xp^lastxp < 0 {
			return false
		}
	}
	return true
}
