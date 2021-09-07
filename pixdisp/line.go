// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
)

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

// DrawLine draws a line from p0 to p1 (including both points).
func (a *Area) DrawLine(p0, p1 image.Point) {
	setColor(a)
	dp := p1.Sub(p0)
	if dp.Y == 0 {
		if dp.X < 0 {
			p1.X, p0.X = p0.X, p1.X
		}
		hline(a, p0.X, p0.Y, p1.X)
		return
	}
	if dp.X == 0 {
		if dp.Y < 0 {
			p1.Y, p0.Y = p0.Y, p1.Y
		}
		vline(a, p0.X, p0.Y, p1.Y)
		return
	}
	vl := abs(dp.Y) > abs(dp.X)
	if vl {
		p0.X, p0.Y = p0.Y, p0.X
		p1.X, p1.Y = p1.Y, p1.X
	}
	if p0.X > p1.X {
		p0, p1 = p1, p0
	}
	dp = p1.Sub(p0).Mul(2)
	sy := 1
	if dp.Y < 0 {
		dp.Y = -dp.Y
		sy = -sy
	}
	e := p0.X - p1.X
	for x := p0.X; x <= p1.X; x++ {
		e += dp.Y
		if e > 0 {
			if vl {
				vline(a, p0.Y, p0.X, x)
			} else {
				hline(a, p0.X, p0.Y, x)
			}
			p0.X = x + 1
			p0.Y += sy
			e -= dp.X
		}
	}
	if p0.X <= p1.X {
		if vl {
			vline(a, p0.Y, p0.X, p1.X)
		} else {
			hline(a, p0.X, p0.Y, p1.X)
		}
	}
}

/*
// DrawLine_ draws a line from p0 to p1 (including both pointsc). DrawLine_
///uses less memory for code than DrawLine but it is generally slower (can be
// faster for very short lines: 1-3 pixels). Use DrawLine_ if you are very
// short of Flash space and do not care about speed or to draw very short lines.
func (a *Area) DrawLine_(p0, p1 image.Point) {
	setColor(a)
	dp := p1.Sub(p0)
	sx, sy := 1, 1
	if dp.X < 0 {
		sx = -sx
	}
	if dp.Y < 0 {
		sy = -sy
	}
	dp.X = abs(dp.X)
	dp.Y = abs(dp.Y)
	e := dp.X - dp.Y
	for {
		drawPixel(a, p0)
		if p0 == p1 {
			return
		}
		e2 := 2 * e
		if e2 > -dp.Y {
			e -= dp.Y
			p0.X += sx
		}
		if e2 < dp.X {
			e += dp.X
			p0.Y += sy
		}
	}
}
*/

/*
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func setPixelColor(a *Area, x, y, c int) {
	a.SetColor(color.Alpha{uint8(255 - c)})
	a.DrawPoint(image.Point{x, y}, 0)
}


func (a *Area) PlotLine(x0, y0, x1, y1 int, wd float32) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy
	var e2, x2, y2 int
	ed := float32(1)
	if dx+dy != 0 {
		ed = float32(math.Sqrt(float64(dx*dx + dy*dy)))
	}
	wd = (wd + 1) / 2
	for {
		setPixelColor(a, x0, y0, max(0, int(255*(float32(abs(err-dx+dy))/ed-wd+1))))
		e2 = err
		x2 = x0
		if 2*e2 >= -dx {
			e2 += dy
			y2 = y0
			for float32(e2) < ed*wd && (y1 != y2 || dx > dy) {
				y2 += sy
				setPixelColor(a, x0, y2, max(0, int(255*(float32(abs(e2))/ed-wd+1))))
				e2 += dx
			}
			if x0 == x1 {
				break
			}
			e2 = err
			err -= dy
			x0 += sx
		}
		if 2*e2 <= dy {
			e2 = dx - e2
			for float32(e2) < ed*wd && (x1 != x2 || dx < dy) {
				x2 += sx
				setPixelColor(a, x2, y0, max(0, int(255*(float32(abs(e2))/ed-wd+1))))
				e2 += dy
			}
			if y0 == y1 {
				break
			}
			err += dx
			y0 += sy
		}
	}
}
*/
