// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/embeddedgo/display/font"
	"github.com/embeddedgo/display/font/subfont"
	"github.com/embeddedgo/display/font/subfont/font9/dejavu12"
	"github.com/embeddedgo/display/font/subfont/font9/terminus12"
	"github.com/embeddedgo/display/pix"
)

var (
	black   = color.Gray{0}
	white   = color.Gray{255}
	red     = color.RGBA{255, 0, 0, 255}
	green   = color.RGBA{0, 255, 0, 255}
	blue    = color.RGBA{0, 0, 255, 255}
	yellow  = color.RGBA{255, 255, 0, 255}
	cyan    = color.RGBA{0, 255, 255, 255}
	magenta = color.RGBA{255, 0, 255, 255}
)

func Colors(disp *pix.Display) {
	a := disp.NewArea(disp.Bounds())
	a.SetColor(black)
	a.Fill(a.Bounds())
	max := a.Bounds().Max
	w := max.X / 4
	div := max.Y - 1
	c := color.RGBA64{A: 0xffff}
	var r image.Rectangle
	t1 := time.Now()
	for y := 0; y < max.Y; y++ {
		r.Min.Y = y
		r.Max.Y = y + 1
		c.R = uint16((y*0xffff + div/2) / div)
		c.G = 0
		c.B = 0
		a.SetColor(c)
		r.Min.X = 0
		r.Max.X = w
		a.Fill(r)
		c.R, c.G = 0, c.R
		a.SetColor(c)
		r.Min.X += w
		r.Max.X += w
		a.Fill(r)
		c.G, c.B = 0, c.G
		a.SetColor(c)
		r.Min.X += w
		r.Max.X += w
		a.Fill(r)
		c.R, c.G = c.B, c.B
		a.SetColor(c)
		r.Min.X += w
		r.Max.X = max.X
		a.Fill(r)
	}
	t2 := time.Now()
	println("Colors:", t2.Sub(t1).String())
	time.Sleep(2 * time.Second)
	t1 = time.Now()
	for y := 0; y < max.Y; y++ {
		r.Min.Y = y
		r.Max.Y = y + 1
		c.R = uint16((y*0xffff + div/2) / div)
		c.G, c.B = c.R, 0
		a.SetColor(c)
		r.Min.X = 0
		r.Max.X = w
		a.Fill(r)
		c.R, c.B = 0, c.G
		a.SetColor(c)
		r.Min.X += w
		r.Max.X += w
		a.Fill(r)
		c.G, c.R = 0, c.B
		a.SetColor(c)
		r.Min.X += w
		r.Max.X += w
		a.Fill(r)
		c.R, c.G = c.B, c.B
		a.SetColor(c)
		r.Min.X += w
		r.Max.X = max.X
		a.Fill(r)
	}
	t2 = time.Now()
	println("Colors:", t2.Sub(t1).String())
	time.Sleep(2 * time.Second)
}

// RotateDisplay draws a white arrow pointing to the top left corner of the
// screen, the red, green and blue bars near top right corner of the screen
// and a purple rounded rectangle around the screen. Next it sequentially
// rotates the coordinate system 90 degrees clockwise and draws it all again.
//
// Pay attention to the position of all elements, the order of RGB bars (the
// red one should be closest to the arrow, the blue one farthest) and the
// direction of rotation.
func RotateDisplay(disp *pix.Display) {
	r := disp.Bounds()
	a := disp.NewArea(r)

	// the dimensions of arrow and bars will be constant, related to w, h, but
	// they position will adapt to current orientation of screen
	w := r.Max.X / 10
	h := r.Dy() / 3
	for i := 0; i < 5; i++ {
		t1 := time.Now()

		// change display direction (orientation)
		disp.SetDir(i)
		a.SetRect(disp.Bounds()) // a.SetRect is mandatory after disp.SetDir
		r = a.Bounds()           // new dimensions of a

		// clear the screen
		a.SetColor(black)
		a.Fill(r)

		// draw the arrow
		a.SetColor(white)
		p0 := image.Pt(w/2, w/2)
		p1 := image.Pt(w, w*3/4)
		p2 := image.Pt(p1.Y, p1.X)
		a.Quad(p0, p0, p1, p2, true)
		a.Line(p0, p1.Add(p2))

		// draw the RGB bars
		var bar image.Rectangle
		bar.Max.X = r.Max.X - w/2
		bar.Min.X = bar.Max.X - w
		bar.Min.Y = w / 2
		bar.Max.Y = bar.Min.Y + h
		o := w * 5 / 4
		a.SetColor(blue)
		a.Fill(bar)
		bar.Min.X -= o
		bar.Max.X -= o
		a.SetColor(green)
		a.Fill(bar)
		bar.Min.X -= o
		bar.Max.X -= o
		a.SetColor(red)
		a.Fill(bar)

		// draw the rounded rectangle around the screen
		r = r.Inset(8)
		r.Max.X--
		r.Max.Y--
		a.SetColorRGBA(128, 0, 128, 255)
		a.RoundRect(r.Min, r.Max, 8, 8, false)

		t2 := time.Now()
		println("RotateDisplay:", t2.Sub(t1).String())
		time.Sleep(2 * time.Second)
	}

}

// The Akkerman Steppe original Polish by Adam Mickiewicz (1798-1855).
const akkermanSteppePL = `Wpłynąłem na suchego przestwór oceanu, Wóz nurza się w zieloność i jak łódka brodzi, Śród fali łąk szumiących, śród kwiatów powodzi, Omijam koralowe ostrowy burzanu.

Już mrok zapada, nigdzie drogi ni kurhanu; Patrzę w niebo, gwiazd szukam, przewodniczek łodzi; Tam z dala błyszczy obłok - tam jutrzenka wschodzi; To błyszczy Dniestr, to weszła lampa Akermanu.

Stójmy! - jak cicho! - słyszę ciągnące żurawie, Których by nie dościgły źrenice sokoła; Słyszę, kędy się motyl kołysa na trawie,

Kędy wąż śliską piersią dotyka się zioła. W takiej ciszy - tak ucho natężam ciekawie, Że słyszałbym głos z Litwy. - Jedźmy, nikt nie woła.`

// The Akkerman Steppe translated to English by Leo Yankevich.
const akkermanSteppeEN = `I launch myself across the dry and open narrows, My carriage plunging into green as if a ketch, Floundering through the meadow flowers in the stretch. I pass an archipelago of coral yarrows.

It's dusk now, not a road in sight, nor ancient barrows. I look up at the sky and look for stars to catch. There distant clouds glint-there tomorrow starts to etch; The Dnieper glimmers; Akkerman's lamp shines and harrows.

I stand in stillness, hear the migratory cranes, Their necks and wings beyond the reach of preying hawks; Hear where the sooty copper glides across the plains,

Where on its underside a viper writhes through stalks. Amid the hush I lean my ears down grassy lanes And listen for a voice from home. Nobody talks.`

var fdejavu = &subfont.Face{
	Height: dejavu12.Height,
	Ascent: dejavu12.Ascent,
	Subfonts: []*subfont.Subfont{
		&dejavu12.X0000_0100,
		&dejavu12.X0101_0201,
	},
}

var fterm = &subfont.Face{
	Height: terminus12.Height,
	Ascent: terminus12.Ascent,
	Subfonts: []*subfont.Subfont{
		&terminus12.X0020_007e,
		&terminus12.X00a0_0175,
		&terminus12.X0178_017f,
	},
}

func randPoint(r image.Rectangle) (p image.Point) {
	p.X = r.Min.X + rand.Int()%r.Dx()
	p.Y = r.Min.Y + rand.Int()%r.Dy()
	return
}

func randQuad(a *pix.Area, r image.Rectangle) {
	p0, p1, p2, p3 := randPoint(r), randPoint(r), randPoint(r), randPoint(r)
	if !pix.IsConvex(p0, p1, p2, p3) {
		p0 = p1
	}
	a.Quad(p0, p1, p2, p3, true)
}

func clearAndPrint(a *pix.Area, face font.Face, s string) {
	t1 := time.Now()
	r := a.Bounds()
	a.SetColor(white)
	a.Fill(r)
	a.SetColorRGBA(70, 140, 210, 210)
	randQuad(a, r)
	a.SetColorRGBA(210, 140, 70, 210)
	randQuad(a, r)
	a.SetColor(black)
	w := a.NewTextWriter(face)
	w.WriteString(s)
	t2 := time.Now()
	println("DrawText:", t2.Sub(t1).String())
	time.Sleep(5 * time.Second)
}

// DrawText draws English and Polish text of the Akkerman Steppe poem on the
// whole screen using two different font faces. The background is white with two
// random quadrilaterals or triangles with different semi-transparent colors.
// The text is black. Dejavu is a proportional font with glyph data encoded as
// 2 bpp image (anti-aliased, 4 transparency levels). VGA is a monospace font
// with glyph data encoded as 1 bpp image (no anti-aliasing).
//
// Pay attention to the background visibility, font anti-aliasing (display
// dependent), non-ASCII letters in Polish text.
func DrawText(disp *pix.Display) {
	a := disp.NewArea(disp.Bounds())
	clearAndPrint(a, fdejavu, akkermanSteppeEN)
	clearAndPrint(a, fdejavu, akkermanSteppePL)
	clearAndPrint(a, fterm, akkermanSteppeEN)
	clearAndPrint(a, fterm, akkermanSteppePL)
}
