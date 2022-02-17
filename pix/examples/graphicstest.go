// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains translation of the Adafruit ILI9341 GFX example.
// The original decrtiption is below.
//
//	This is our GFX example for the Adafruit ILI9341 Breakout and Shield
//	----> http://www.adafruit.com/products/1651
//	Check out the links above for our tutorials and wiring diagrams
//	These displays use SPI to communicate, 4 or 5 pins are required to
//	interface (RST is optional)
//	Adafruit invests time and resources providing this open source code,
//	please support Adafruit and open-source hardware by purchasing
//	products from Adafruit!
//	Written by Limor Fried/Ladyada for Adafruit Industries.
//	MIT license, all text above must be included in any redistribution

package examples

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/embeddedgo/display/font"
	"github.com/embeddedgo/display/font/subfont/font9/terminus12"
	"github.com/embeddedgo/display/pix"
)

func delay(ms int) { time.Sleep(time.Duration(ms) * time.Millisecond) }

func GraphicsTest(disp *pix.Display) {
	a := disp.NewArea(disp.Bounds())

	println("\n** Adafruit Graphics Test **")
	println("Benchmark                Time (microseconds)")
	delay(10)

	print("Screen fill              ")
	println(testFillScreen(a))
	delay(500)

	print("Text                     ")
	println(testText(a))
	delay(3000)

	print("Lines                    ")
	println(testLines(a, cyan))
	delay(500)

	print("Horiz/Vert Lines         ")
	println(testFastLines(a, red, blue))
	delay(500)

	print("Rectangles (outline)     ")
	println(testRects(a, green))
	delay(500)

	print("Rectangles (filled)      ")
	println(testFilledRects(a, yellow, magenta))
	delay(500)

	print("Circles (filled)         ")
	println(testFilledCircles(a, 10, magenta))

	print("Circles (outline)        ")
	println(testCircles(a, 10, white))
	delay(500)

	print("Triangles (outline)      ")
	println(testTriangles(a))
	delay(500)

	print("Triangles (filled)       ")
	println(testFilledTriangles(a))
	delay(500)

	print("Rounded rects (outline)  ")
	println(testRoundRects(a))
	delay(500)

	print("Rounded rects (filled)   ")
	println(testFilledRoundRects(a))
	delay(500)

	print("Done!\n")

	for rotation := 0; rotation < 5; rotation++ {
		disp.SetDir(rotation)
		a.SetRect(disp.Bounds())
		testText(a)
		delay(1000)
	}
}

func testFillScreen(a *pix.Area) uint {
	start := time.Now()
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(red)
	a.Fill(r)
	a.SetColor(green)
	a.Fill(r)
	a.SetColor(blue)
	a.Fill(r)
	a.SetColor(black)
	a.Fill(r)
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testText(a *pix.Area) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(white)
	ff := terminus12.NewFace(terminus12.X0020_007e)
	w := a.NewTextWriter(ff)
	start := time.Now()
	fmt.Fprint(w, "Hello World!\n\n")
	w.SetColor(yellow)
	w.Face = font.Magnify(w.Face, 2, 2, font.Nearest)
	fmt.Fprint(w, 1234.56)
	w.SetColor(red)
	w.Face.(*font.Magnifier).SetScale(3, 3)
	fmt.Fprintf(w, "\n%X", uint32(0xDEADBEEF))
	w.SetColor(green)
	w.Face.(*font.Magnifier).SetScale(5, 5)
	fmt.Fprint(w, "\nGroop")
	w.Face.(*font.Magnifier).SetScale(2, 2)
	fmt.Fprint(w, "\nI implore thee,")
	w.Face = ff
	fmt.Fprintln(w, "\nmy foonting turlingdromes.")
	fmt.Fprintln(w, "And hooptiously drangle me")
	fmt.Fprintln(w, "with crinkly bindlewurdles,")
	fmt.Fprintln(w, "Or I will rend thee")
	fmt.Fprintln(w, "in the gobberwarts")
	fmt.Fprintln(w, "with my blurglecruncheon,")
	fmt.Fprintln(w, "see if I don't!")
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testLines(a *pix.Area, color color.Color) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(color)
	w, h := r.Max.X, r.Max.Y
	var p1, p2 image.Point
	p2.Y = h - 1
	start := time.Now()
	for p2.X = 0; p2.X < w; p2.X += 6 {
		a.Line(p1, p2)
	}
	p2.X = w - 1
	for p2.Y = 0; p2.Y < h; p2.Y += 6 {
		a.Line(p1, p2)
	}
	t := time.Now().Sub(start) // fillScreen doesn't count against timing
	a.Flush()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(color)
	p1.X = w - 1
	p1.Y = 0
	p2.Y = h - 1
	start = time.Now()
	for p2.X = 0; p2.X < w; p2.X += 6 {
		a.Line(p1, p2)
	}
	p2.X = 0
	for p2.Y = 0; p2.Y < h; p2.Y += 6 {
		a.Line(p1, p2)
	}
	t += time.Now().Sub(start)
	a.Flush()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(color)
	p1.X = 0
	p1.Y = h - 1
	p2.Y = 0
	start = time.Now()
	for p2.X = 0; p2.X < w; p2.X += 6 {
		a.Line(p1, p2)
	}
	p2.X = w - 1
	for p2.Y = 0; p2.Y < h; p2.Y += 6 {
		a.Line(p1, p2)
	}
	t += time.Now().Sub(start)
	a.Flush()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(color)
	p1.X = w - 1
	p1.Y = h - 1
	p2.Y = 0
	start = time.Now()
	for p2.X = 0; p2.X < w; p2.X += 6 {
		a.Line(p1, p2)
	}
	p2.X = 0
	for p2.Y = 0; p2.Y < h; p2.Y += 6 {
		a.Line(p1, p2)
	}
	t += time.Now().Sub(start)
	a.Flush()
	return uint(t / 1e3)
}

func testFastLines(a *pix.Area, color1, color2 color.Color) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(color1)
	start := time.Now()
	var p0, p1 image.Point
	p1.X = r.Max.X - 1
	for p0.Y = 0; p0.Y < r.Max.Y; p0.Y += 5 {
		p1.Y = p0.Y
		a.Line(p0, p1)
	}
	a.SetColor(color2)
	p0.Y = 0
	for p0.X = 0; p0.X < r.Max.X; p0.X += 5 {
		p1.X = p0.X
		a.Line(p0, p1)
	}
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testRects(a *pix.Area, color color.Color) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(color)
	c := r.Max.Div(2)
	n := min(r.Max.X, r.Max.Y)
	start := time.Now()
	for i := 2; i < n; i += 6 {
		d := image.Pt(i/2, i/2)
		a.RoundRect(c.Sub(d), c.Add(d), 0, 0, false)
	}
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testFilledRects(a *pix.Area, color1, color2 color.Color) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	c := r.Max.Div(2).Sub(image.Pt(1, 1))
	n := min(r.Max.X, r.Max.Y)
	var t time.Duration
	for i := n; i > 0; i -= 6 {
		d := image.Pt(i/2, i/2)
		a.SetColor(color1)
		start := time.Now()
		a.RoundRect(c.Sub(d), c.Add(d), 0, 0, true)
		t += time.Now().Sub(start)
		// Outlines are not included in timing results
		a.SetColor(color2)
		a.RoundRect(c.Sub(d), c.Add(d), 0, 0, false)
	}
	a.Flush()
	return uint(t / 1e3)
}

func testFilledCircles(a *pix.Area, radius int, color color.Color) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	a.SetColor(color)
	r2 := radius * 2
	var p image.Point
	start := time.Now()
	for p.X = radius; p.X < r.Max.X; p.X += r2 {
		for p.Y = radius; p.Y < r.Max.Y; p.Y += r2 {
			a.RoundRect(p, p, radius, radius, true)
		}
	}
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testCircles(a *pix.Area, radius int, color color.Color) uint {
	// Screen is not cleared for this one -- this is
	// intentional and does not affect the reported time.
	r := a.Bounds()
	r.Max = r.Max.Add(image.Pt(radius, radius))
	a.SetColor(color)
	r2 := radius * 2
	var p image.Point
	start := time.Now()
	for p.X = 0; p.X < r.Max.X; p.X += r2 {
		for p.Y = 0; p.Y < r.Max.Y; p.Y += r2 {
			a.RoundRect(p, p, radius, radius, false)
		}
	}
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testTriangles(a *pix.Area) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	c := r.Max.Div(2).Sub(image.Pt(1, 1))
	n := min(c.X, c.Y)
	start := time.Now()
	for i := 0; i < n; i += 5 {
		a.SetColor(color.Gray{uint8(i)})
		p0 := c.Add(image.Pt(0, -i)) // peak
		p1 := c.Add(image.Pt(-i, i)) // bottom left
		p2 := c.Add(image.Pt(i, i))  // bottom right
		a.Quad(p0, p0, p1, p2, false)
	}
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testFilledTriangles(a *pix.Area) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	c := r.Max.Div(2).Sub(image.Pt(1, 1))
	var t time.Duration
	for i := min(c.X, c.Y); i > 10; i -= 5 {
		a.SetColorRGBA(0, uint8(i*10), uint8(i*10), 255)
		p0 := c.Add(image.Pt(0, -i)) // peak
		p1 := c.Add(image.Pt(-i, i)) // bottom left
		p2 := c.Add(image.Pt(i, i))  // bottom right
		start := time.Now()
		a.Quad(p0, p0, p1, p2, true)
		t += time.Now().Sub(start)
		a.SetColorRGBA(uint8(i*10), uint8(i*10), 0, 255)
		a.Quad(p0, p0, p1, p2, false)
	}
	a.Flush()
	return uint(t / 1e3)
}

func testRoundRects(a *pix.Area) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	c := r.Max.Div(2).Sub(image.Pt(1, 1))
	w := min(r.Max.X, r.Max.Y)
	start := time.Now()
	for i := 0; i < w; i += 6 {
		r := i / 8
		d := image.Pt(i/2-r, i/2-r)
		a.SetColorRGBA(uint8(i), 0, 0, 255)
		a.RoundRect(c.Sub(d), c.Add(d), r, r, false)
	}
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func testFilledRoundRects(a *pix.Area) uint {
	r := a.Bounds()
	a.SetColor(black)
	a.Fill(r)
	c := r.Max.Div(2).Sub(image.Pt(1, 1))
	w := min(r.Max.X, r.Max.Y)
	start := time.Now()
	for i := w; i > 20; i -= 6 {
		r := i / 8
		d := image.Pt(i/2-r, i/2-r)
		a.SetColorRGBA(0, uint8(i), 0, 255)
		a.RoundRect(c.Sub(d), c.Add(d), r, r, true)
	}
	a.Flush()
	return uint(time.Now().Sub(start) / 1e3)
}

func min(a, b int) int {
	if b < a {
		a = b
	}
	return a
}
