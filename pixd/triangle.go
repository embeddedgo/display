// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

func (a *Area) Triangle(p0, p1, p2 image.Point, fill bool) {
	// order the vertices by Y to determine the height of triangle
	if p0.Y > p2.Y {
		p0, p2 = p2, p0
	}
	if p0.Y > p1.Y {
		p0, p1 = p1, p0
	} else if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}
	height := p2.Y - p0.Y
	// order the vertices by X to determine the width of triangle
	q0, q1, q2 := p0, p1, p2
	if q0.X > q2.X {
		q0, q2 = q2, q0
	}
	if q0.X > q1.X {
		q0, q1 = q1, q0
	} else if q1.X > q2.X {
		q1, q2 = q2, q1
	}
	width := q2.X - q0.X
	// draw left to right if width < height, otherwise top to bottom
	rot := width < height
	if rot {
		p0, p1, p2 = q0, q1, q2
	}
	// order the vertices as left-top-bottom / top-left-right
	q1, q2 = p1.Sub(p0), p2.Sub(p0)
	z := q1.X*q2.Y - q1.Y*q2.X // BUG: a huge triangle can overflow 32-bit int
	if rot {
		z = -z
	}
	if z > 0 {
		p1, p2 = p2, p1
	}
	if fill {
		// Read the code below as if it was always drawing horizontal lines
		// (flx,y)-(frx,y), line by line, from top (p0) to bottom. If rot is
		// true the actual drawing is from left to right, using vertical lines.
		if rot {
			p0.X, p0.Y = p0.Y, p0.X
			p1.X, p1.Y = p1.Y, p1.X
			p2.X, p2.Y = p2.Y, p2.X
		}
		ldx, lsx, ldy, le := dxsxdye(p1, p0)
		lx := p0.X
		rdx, rsx, rdy, re := dxsxdye(p2, p0)
		rx := p0.X
		y := p0.Y
		bottom := p1.Y
		if bottom < p2.Y {
			bottom = p2.Y
		}
		lend := p1.X
		rend := p2.X
		for {
			var px int
			flx := lx
			lx, px, le = nextX(lx, ldx, lsx, lend, ldy, le)
			if lsx < 0 {
				flx = px
			}
			frx := rx
			rx, px, re = nextX(rx, rdx, rsx, rend, rdy, re)
			if rsx > 0 {
				frx = px
			}
			r := image.Rectangle{image.Pt(flx, y), image.Pt(frx+1, y+1)}
			if rot {
				r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
				r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
			}
			a.Fill(r)
			if y == p1.Y {
				if y == bottom {
					break
				}
				ldx, lsx, ldy, le = dxsxdye(p2, p1)
				lend = p2.X
				lx, _, le = nextX(lx, ldx, lsx, lend, ldy, le)
			} else if y == p2.Y {
				if y == bottom {
					break
				}
				rdx, rsx, rdy, re = dxsxdye(p1, p2)
				rend = p1.X
				rx, _, re = nextX(rx, rdx, rsx, rend, rdy, re)
			}
			le += ldx
			re += rdx
			y++
		}
	} else {
		// ensure the same drawing direction for filed and empty triangle to
		// obtain the  same edges (Bresenham / Zingl lines are not symetrical)
		a.Line(p0, p1)
		a.Line(p0, p2)
		if rot {
			if p1.X > p2.X {
				p1, p2 = p2, p1
			}
		} else {
			if p1.Y > p2.Y {
				p1, p2 = p2, p1
			}
		}
		a.Line(p1, p2)
	}
}

func dxsxdye(p, top image.Point) (dx, sx, dy, e int) {
	dx, sx = absSign(p.X - top.X)
	dy = top.Y - p.Y
	e = dx + dy
	return
}

// nextX returns x for the point on the line where y changes to y+1
func nextX(x, dx, sx, end, dy, e int) (nx, px, ne int) {
	// based on Alois Zingl line drawing algorithm
	for {
		px = x
		e2 := 2 * e
		if e2 >= dy {
			if x == end {
				break
			}
			e += dy
			x += sx
		}
		if e2 <= dx {
			break
		}
	}
	return x, px, e
}
