// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

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

func dxsxdye(p, top image.Point) (dx, sx, dy, e int) {
	dx, sx = absSign(p.X - top.X)
	dy = top.Y - p.Y
	e = dx + dy
	return
}

func fillTriangle(a *Area, top, left, right image.Point) {
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
		plx := lx
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
			if lsx < 0 {
				plx = lx
			}
		}
		prx := rx
		for {
			prx = rx
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
			if rsx > 0 {
				prx = rx
			}
		}
		a.Fill(image.Rectangle{image.Pt(plx, y), image.Pt(prx+1, y+1)})
		switch y {
		case bottom:
			return
		case left.Y:
			ldx, lsx, ldy, le = dxsxdye(right, left)
			lend = right.X
		case right.Y:
			rdx, rsx, rdy, re = dxsxdye(left, right)
			rend = left.X
		}
		le += ldx
		re += rdx
		y++
	}
}
