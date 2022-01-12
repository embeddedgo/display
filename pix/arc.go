// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"

	"github.com/embeddedgo/display/math2d"
)

const frac = 18 // maxa, maxb up to 8191 without 64-bit multiplication

func (a *Area) Arc(p image.Point, mina, minb, maxa, maxb int, th0, th1 int32, fill bool) {
	// bounding box
	var box image.Rectangle
	box.Min.X = -maxa
	box.Max.X = maxa
	box.Min.Y = -maxb
	box.Max.Y = maxb
	// setup two sides
	var w0, w1, dx0, dx1, dy0, dy1, of int
	if th0 != th1 {
		one := image.Pt(1<<frac, 0)
		cosSin := math2d.Rotate(one, th0)
		pmin0 := mulfi(cosSin, mina, minb)
		pmax0 := mulfi(cosSin, maxa, maxb)
		cosSin = math2d.Rotate(one, th1)
		pmin1 := mulfi(cosSin, mina, minb)
		pmax1 := mulfi(cosSin, maxa, maxb)
		dx0 = pmax0.X - pmin0.X
		dx1 = pmin1.X - pmax1.X
		dy0 = pmin0.Y - pmax0.Y
		dy1 = pmax1.Y - pmin1.Y
		of = int(th1 - th0)
		if of < 0 {
			dy0, dy1 = -dy0, -dy1
			dx0, dx1 = -dx0, -dx1
		}
		top := th0 < th1 && th1 <= 0
		bottom := uint32(th0) < uint32(th1) && uint32(th1) <= math2d.FullAngle/2
		tr0, tr1 := th0+math2d.RightAngle, th1+math2d.RightAngle
		left := tr0 < tr1 && tr1 <= 0
		right := uint32(tr0) < uint32(tr1) && uint32(tr1) <= math2d.FullAngle/2
		if !fill {
			a.Line(pmin0.Add(p), pmax0.Add(p))
			w0 := -dx0*pmin0.Y - dy0*pmin0.X
			w1 := -dx1*pmin1.Y - dy1*pmin1.X
			if !top {
				if !left {
					arc(a, p, mina, minb, w0, dx0, dy0, w1, dx1, dy1, of, -1, 1)
				}
				if !right {
					arc(a, p, mina, minb, w0, dx0, dy0, w1, dx1, dy1, of, 1, 1)
				}
			}
			if !bottom {
				if !right {
					arc(a, p, mina, minb, w0, dx0, dy0, w1, dx1, dy1, of, 1, -1)
				}
				if !left {
					arc(a, p, mina, minb, w0, dx0, dy0, w1, dx1, dy1, of, -1, -1)
				}
			}
			w0 = -dx0*pmax0.Y - dy0*pmax0.X
			w1 = -dx1*pmax1.Y - dy1*pmax1.X
			if !top {
				if !left {
					arc(a, p, maxa, maxb, w0, dx0, dy0, w1, dx1, dy1, of, -1, 1)
				}
				if !right {
					arc(a, p, maxa, maxb, w0, dx0, dy0, w1, dx1, dy1, of, 1, 1)
				}
			}
			if !bottom {
				if !right {
					arc(a, p, maxa, maxb, w0, dx0, dy0, w1, dx1, dy1, of, 1, -1)
				}
				if !left {
					arc(a, p, maxa, maxb, w0, dx0, dy0, w1, dx1, dy1, of, -1, -1)
				}
			}
			a.Line(pmin1.Add(p), pmax1.Add(p))
			return
		}
		// more fitted box
		if bottom {
			box.Min.X = pmax1.X
			box.Max.X = pmax0.X
		} else if top {
			box.Min.X = pmax0.X
			box.Max.X = pmax1.X
		}
		if right {
			box.Min.Y = pmax0.Y
			box.Max.Y = pmax1.Y
		} else if left {
			box.Min.Y = pmax1.Y
			box.Max.Y = pmax0.Y
		}
		if bottom {
			box.Min.Y = min(pmin0.Y, pmin1.Y)
		} else if top {
			box.Max.Y = max(pmin0.Y, pmin1.Y)
		}
		if right {
			box.Min.X = min(pmin0.X, pmin1.X)
		} else if left {
			box.Max.X = max(pmin0.X, pmin1.X)
		}
		w0 = dx0*(box.Min.Y-pmin0.Y) + dy0*(box.Min.X-pmin0.X)
		w1 = dx1*(box.Min.Y-pmin1.Y) + dy1*(box.Min.X-pmin1.X)
		if of < 0 {
			w0--
		} else {
			w1--
		}
	} else if !fill {
		a.RoundRect(p, p, mina, minb, false)
		a.RoundRect(p, p, maxa, maxb, false)
		return
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
	wmax := maxaa*maxbb - (xx*maxbb + yy*maxaa) - 1
	minaa2 := minaa * 2
	minbb2 := minbb * 2
	maxaa2 := maxaa * 2
	maxbb2 := maxbb * 2
	dwymin := minaa2*box.Min.Y + minaa
	dwymax := maxaa2*box.Min.Y + maxaa
	dwxmin := minbb2*box.Min.X + minbb
	dwxmax := maxbb2*box.Min.X + maxbb
	// fill
	box = box.Add(p)
	var r image.Rectangle
	if box.Dx() >= box.Dy() {
		for r.Min.Y = box.Min.Y; r.Min.Y < box.Max.Y; r.Min.Y++ {
			m0, m1 := w0, w1
			mmin, mmax := wmin, wmax
			dmmin, dmmax := dwxmin, dwxmax
			r.Min.X = box.Max.X
			r.Max.X = box.Min.X
			for {
				m := (m0 | m1) ^ of | (mmin | mmax)
				if r.Min.X == box.Max.X {
					if r.Max.X == box.Max.X {
						break
					}
					if m >= 0 {
						r.Min.X = r.Max.X
					}
				} else {
					if m < 0 || r.Max.X == box.Max.X {
						r.Max.Y = r.Min.Y + 1
						a.Fill(r)
						if r.Max.X == box.Max.X {
							break
						}
						r.Min.X = box.Max.X
					}
				}
				m0 += dy0
				m1 += dy1
				mmin += dmmin
				mmax -= dmmax
				dmmin += minbb2
				dmmax += maxbb2
				r.Max.X++
			}
			w0 += dx0
			w1 += dx1
			wmin += dwymin
			wmax -= dwymax
			dwymin += minaa2
			dwymax += maxaa2
		}
	} else {
		for r.Min.X = box.Min.X; r.Min.X < box.Max.X; r.Min.X++ {
			m0, m1 := w0, w1
			mmin, mmax := wmin, wmax
			dmmin, dmmax := dwymin, dwymax
			r.Min.Y = box.Max.Y
			r.Max.Y = box.Min.Y
			for {
				m := (m0 | m1) ^ of | (mmin | mmax)
				if r.Min.Y == box.Max.Y {
					if r.Max.Y == box.Max.Y {
						break
					}
					if m >= 0 {
						r.Min.Y = r.Max.Y
					}
				} else {
					if m < 0 || r.Max.Y == box.Max.Y {
						r.Max.X = r.Min.X + 1
						a.Fill(r)
						if r.Max.Y == box.Max.Y {
							break
						}
						r.Min.Y = box.Max.Y
					}
				}
				m0 += dx0
				m1 += dx1
				mmin += dmmin
				mmax -= dmmax
				dmmin += minaa2
				dmmax += maxaa2
				r.Max.Y++
			}
			w0 += dy0
			w1 += dy1
			wmin += dwxmin
			wmax -= dwxmax
			dwxmin += minbb2
			dwxmax += maxbb2
		}
	}
}

func mulfi(p image.Point, mx, my int) image.Point {
	const round = 1 << (frac - 1)
	p.X = (p.X*mx + round) >> frac
	p.Y = (p.Y*my + round) >> frac
	return p
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func arc(a *Area, p image.Point, ra, rb, w0, dx0, dy0, w1, dx1, dy1, of int, dirx, diry int8) {
	if ra|rb == 0 {
		a.Pixel(p.X, p.Y)
		return
	}
	// based on Alois Zingl algorithm
	dx := (1 - 2*ra) * rb * rb
	dy := ra * ra
	e := dx + dy
	bb2 := 2 * rb * rb
	aa2 := 2 * ra * ra
	x := -ra * int(dirx)
	y := 0
	w0 += dy0 * x
	w1 += dy1 * x
	if of >= 0 {
		w0--
		w1--
	}
	dx0 *= int(diry)
	dx1 *= int(diry)
	dy0 *= int(dirx)
	dy1 *= int(dirx)
	if dirx != diry && (w0|w1)^of >= 0 {
		a.Pixel(p.X+x, p.Y)
	}
	for {
		e2 := 2 * e
		if e2 >= dx {
			x += int(dirx)
			if x == 0 && dirx == diry {
				return
			}
			dx += bb2
			e += dx
			w0 += dy0
			w1 += dy1
		}
		if e2 <= dy {
			y += int(diry)
			dy += aa2
			e += dy
			w0 += dx0
			w1 += dx1
		}
		if (w0|w1)^of >= 0 {
			a.Pixel(p.X+x, p.Y+y)
		}
		if x == 0 {
			break
		}
	}
	rb *= int(diry)
	for y != rb {
		y += int(diry)
		w0 += dx0
		w1 += dx1
		if (w0|w1)^of >= 0 {
			a.Pixel(p.X, p.Y+y)
		}
	}
}
