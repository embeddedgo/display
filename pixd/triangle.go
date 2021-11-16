// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

/*
func (a *Area) Triangle(p0, p1, p2 image.Point, fill bool) {
	// Order the points as top, left, right.
	if p0.Y > p1.Y {
		p0, p1 = p1, p0
	}
	if p0.Y > p2.Y {
		p0, p2 = p2, p0
	}
	if p1.X > p2.X {
		p1, p2 = p2, p1
	}
	if fill {
		fillTriangle(a, p0, p1, p2)
	} else {
		// The same drawing direction is required for filed and empty triangle
		// to ensure equal edges.
		a.Line(p0, p1)
		a.Line(p0, p2)
		if p1.Y > p2.Y {
			p1, p2 = p2, p1
		}
		a.Line(p1, p2)
		return
	}
}
*/

func (a *Area) Triangle(p0, p1, p2 image.Point, fill bool) {
	// order the vertices by Y to determine the triangle height
	if p0.Y > p2.Y {
		p0, p2 = p2, p0
	}
	if p0.Y > p1.Y {
		p0, p1 = p1, p0
	} else if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}
	height := p2.Y - p0.Y
	// order the vertices by X to determine the triangle width
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
		// order the vertices as left-top-bottom
		if p1.Y > q2.Y {
			p1, p2 = p2, p1
		}
	} else {
		// order the vertices as top-left-right
		if p1.X > p2.X {
			p1, p2 = p2, p1
		}
	}
	if fill {
		fillTriangle(a, p0, p1, p2, rot)
	} else {
		// ensure the same drawing direction for filed and empty triangle
		// for same edges (Bresenham / Zingl lines are not symetrical)
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

// nextX returns x for the point of a line where y changes to y+1
func nextX(x, dx, sx, end, dy, e int) (nx, ne int) {
	// based on Alois Zingl line drawing algorithm
	for {
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
	return x, e
}

// If rot is false fillTriangle fills the triangle described by top, left,
// right vertices drawing horizontal lines from top to bottom. If rot is true
// it fills the triangle drawing vertical lines from left to right. In this
// case top vertex is left one, left is top, right is bottom.
func fillTriangle(a *Area, top, left, right image.Point, rot bool) {
	if rot {
		top.X, top.Y = top.Y, top.X
		left.X, left.Y = left.Y, left.X
		right.X, right.Y = right.Y, right.X
	}
	ldx, lsx, ldy, le := dxsxdye(left, top)
	lx := top.X
	rdx, rsx, rdy, re := dxsxdye(right, top)
	rx := top.X
	y := top.Y
	bottom := left.Y
	if bottom < right.Y {
		bottom = right.Y
	}
	lend := left.X
	rend := right.X
	for {
		flx := lx
		lx, le = nextX(lx, ldx, lsx, lend, ldy, le)
		if lsx < 0 {
			flx = lx
		}
		frx := rx
		rx, re = nextX(rx, rdx, rsx, rend, rdy, re)
		if rsx > 0 {
			frx = rx
		}
		var r image.Rectangle
		if rot {
			r.Min.X = y
			r.Max.X = y + 1
			r.Min.Y = flx
			r.Max.Y = frx + 1
		} else {
			r.Min.X = flx
			r.Max.X = frx + 1
			r.Min.Y = y
			r.Max.Y = y + 1
		}
		a.Fill(r)
		if y == left.Y {
			if y == bottom {
				break
			}
			ldx, lsx, ldy, le = dxsxdye(right, left)
			lend = right.X
			lx, le = nextX(lx, ldx, lsx, lend, ldy, le)
		} else if y == right.Y {
			if y == bottom {
				break
			}
			rdx, rsx, rdy, re = dxsxdye(left, right)
			rend = left.X
			rx, re = nextX(rx, rdx, rsx, rend, rdy, re)
		}
		le += ldx
		re += rdx
		y++
	}
}

/*
func fillTriangle(a *Area, top, left, right image.Point, rot bool) {
	if rot {
		top.X, top.Y = top.Y, top.X
		left.X, left.Y = left.Y, left.X
		right.X, right.Y = right.Y, right.X
	}
	ldx, lsx, ldy, le := dxsxdye(left, top)
	lx := top.X
	rdx, rsx, rdy, re := dxsxdye(right, top)
	rx := top.X
	y := top.Y
	bottom := left.Y
	if bottom < right.Y {
		bottom = right.Y
	}
	lend := left.X
	rend := right.X
	for {
		flx := lx
		for {
			le2 := 2 * le
			if le2 >= ldy {
				if lx == lend {
					break
				}
				le += ldy
				lx += lsx
			}
			if le2 <= ldx {
				break
			}
		}
		if lsx < 0 {
			flx = lx
		}
		frx := rx
		for {
			re2 := 2 * re
			if re2 >= rdy {
				if rx == rend {
					break
				}
				re += rdy
				rx += rsx
			}
			if re2 <= rdx {
				break
			}
		}
		if rsx > 0 {
			frx = rx
		}
		var r image.Rectangle
		if rot {
			r.Min.X = y
			r.Max.X = y + 1
			r.Min.Y = flx
			r.Max.Y = frx + 1
		} else {
			r.Min.X = flx
			r.Max.X = frx + 1
			r.Min.Y = y
			r.Max.Y = y + 1
		}
		a.Fill(r)
		if y == left.Y {
			if y == bottom {
				break
			}
			ldx, lsx, ldy, le = dxsxdye(right, left)
			lend = right.X
			for {
				le2 := 2 * le
				if le2 >= ldy {
					if lx == lend {
						break
					}
					le += ldy
					lx += lsx
				}
				if le2 <= ldx {
					break
				}
			}
		} else if y == right.Y {
			if y == bottom {
				break
			}
			rdx, rsx, rdy, re = dxsxdye(left, right)
			rend = left.X
			for {
				re2 := 2 * re
				if re2 >= rdy {
					if rx == rend {
						break
					}
					re += rdy
					rx += rsx
				}
				if re2 <= rdx {
					break
				}
			}
		}
		le += ldx
		re += rdx
		y++
	}
}
*/
