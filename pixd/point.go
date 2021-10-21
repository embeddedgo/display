// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
)

// DrawPoint draws a point with a given radius.
func (a *Area) DrawPoint(p image.Point, r int) {
	if r <= 0 {
		if r == 0 {
			drawPixel(a, p)
		}
		return
	}
	// based on Alois Zingl algorithm
	// fill the four sides of the circle
	x, y, e := -r, 0, 2*(1-r)
	for x+y < 0 {
		ny := y + 1
		ne := e + 2*ny + 1
		if e > x || ne > ny {
			x0, x1 := p.X-y, p.X+y
			hline(a, x0, p.Y-x, x1)
			hline(a, x0, p.Y+x, x1)
			y0, y1 := p.Y-y, p.Y+y
			vline(a, p.X+x, y0, y1)
			vline(a, p.X-x, y0, y1)
			x += 1
			ne += 2*x + 1
		}
		y = ny
		e = ne
	}
	// fill the center rectangle
	rect := image.Rectangle{
		p.Add(image.Point{x, x}),
		p.Sub(image.Point{x - 1, x - 1}),
	}
	a.Fill(rect)
}

/*
The implementation below is 1.5x slower than the above one in case of ILI9341,
21 MHz SPI.

func (a *Area) DrawPoint(p image.Point, r int) {
	if r <= 0 {
		if r == 0 {
			drawPixel(a, p)
		}
		return
	}
	// based on Alois Zingl algorithm
	x, y, e := -r, 1, 5-2*r
	hline(a, p.X-x, p.Y, p.X+x)
	for {
		e0 := e
		if e0 <= y {
			hline(a, p.X-x, p.Y-y, p.X+x)
			hline(a, p.X-x, p.Y+y, p.X+x)
			y++
			e += 2*y + 1
		}
		if e0 > x || e > y {
			if x++; x >= 0 {
				break
			}
			e += 2*x + 1
		}
	}
}
*/

/*
// DrawPoint draws a point with a given radius.
func (a *Area) DrawPoint(p image.Point, r int) {
	if r <= 0 {
		if r == 0 {
			drawPixel(a, p)
		}
		return
	}
	// based on Alois Zingl algorithm
	// fill the four sides of the circle
	x, y, e := -r, 0, 2*(1-r)
	for x+y < 0 {
		ny := y + 1
		ne := e + 2*ny + 1
		if e > x || ne > ny {
			var rect image.Rectangle
			rect.Min.X = p.X - y
			rect.Min.Y = p.Y - x
			rect.Max.X = p.X + y + 1
			rect.Max.Y = rect.Min.Y + 1
			a.Fill(rect)
			rect.Min.Y = p.Y + x
			rect.Max.Y = rect.Min.Y + 1
			a.Fill(rect)
			rect.Min.X = p.X + x
			rect.Min.Y = p.Y - y
			rect.Max.X = rect.Min.X + 1
			rect.Max.Y = p.Y + y
			a.Fill(rect)
			rect.Min.X = p.X - x
			rect.Max.X = rect.Min.X + 1
			a.Fill(rect)
			x += 1
			ne += 2*x + 1
		}
		y = ny
		e = ne
	}
	// fill the center rectangle
	rect := image.Rectangle{
		p.Add(image.Point{x, x}),
		p.Sub(image.Point{x - 1, x - 1}),
	}
	a.Fill(rect)
}
*/
