// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/embeddedgo/display/pixd"
	"github.com/embeddedgo/display/pixd/font"
	"github.com/embeddedgo/display/pixd/font/font9/dejavu12"
	"github.com/embeddedgo/display/pixd/font/font9/vga"
)

var (
	black = color.Gray{0}
	white = color.Gray{255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
)

// RotateDisplay draws a white arrow pointing to the top left corner of the
// screen, the red, green and blue bars on near top right corner of the screen
// and a purple rounded rectangle around the screen. Next it sequentially
// rotates the coordinate system 90 degrees clockwise and draws it all again.
//
// Pay attention to the position of all elements, the order of RGB bars (the
// red one should be closest to the arrow, the blue one farthest) and the
// direction of rotation.
func RotateDisplay(disp *pixd.Display, n int) {
	r := disp.Bounds()
	a := disp.NewArea(r)

	// the dimensions of arrow and bars will be constant, related to w, h, but
	// they position will adapt to current orientation of screen
	w := r.Max.X / 10
	h := r.Dy() / 3
	for i := 0; i != n; i++ {
		t1 := time.Now()

		// change display direction (orientation)
		disp.SetDir(i)
		a.SetRect(disp.Bounds()) // a.SetRect is mandatory after disp.SetDir
		r = a.Bounds()           // new dimensions of a

		// clear the screen
		a.SetColorRGBA(0, 0, 0, 255)
		a.Fill(r)

		// draw the arrow
		a.SetColorRGBA(255, 255, 255, 255)
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

// The Akkerman Steppes original Polish by Adam Mickiewicz (1798-1855).
const akkermanSteppesPL = `Wpłynąłem na suchego przestwór oceanu, Wóz nurza się w zieloność i jak łódka brodzi, Śród fali łąk szumiących, śród kwiatów powodzi, Omijam koralowe ostrowy burzanu.

Już mrok zapada, nigdzie drogi ni kurhanu; Patrzę w niebo, gwiazd szukam, przewodniczek łodzi; Tam z dala błyszczy obłok - tam jutrzenka wschodzi; To błyszczy Dniestr, to weszła lampa Akermanu.

Stójmy! - jak cicho! - słyszę ciągnące żurawie, Których by nie dościgły źrenice sokoła; Słyszę, kędy się motyl kołysa na trawie,

Kędy wąż śliską piersią dotyka się zioła. W takiej ciszy - tak ucho natężam ciekawie, Że słyszałbym głos z Litwy. - Jedźmy, nikt nie woła.`

// The Akkerman Steppes translated to English by Leo Yankevich.
const akkermanSteppesEN = `I launch myself across the dry and open narrows, My carriage plunging into green as if a ketch, Floundering through the meadow flowers in the stretch. I pass an archipelago of coral yarrows.

It's dusk now, not a road in sight, nor ancient barrows. I look up at the sky and look for stars to catch. There distant clouds glint-there tomorrow starts to etch; The Dnieper glimmers; Akkerman's lamp shines and harrows.

I stand in stillness, hear the migratory cranes, Their necks and wings beyond the reach of preying hawks; Hear where the sooty copper glides across the plains,

Where on its underside a viper writhes through stalks. Amid the hush I lean my ears down grassy lanes And listen for a voice from home. Nobody talks.`

var fdejavu = &font.Face{
	Height: dejavu12.Height,
	Ascent: dejavu12.Ascent,
	Subfonts: []*font.Subfont{
		&dejavu12.X0000_0100,
		&dejavu12.X0101_0201,
	},
}

var fvga = &font.Face{
	Height: vga.Height,
	Ascent: vga.Ascent,
	Subfonts: []*font.Subfont{
		&vga.X0000_007f,
		&vga.X00a0_021f,
	},
}

func rndPt(r image.Rectangle) (p image.Point) {
	p.X = r.Min.X + rand.Int()%r.Dx()
	p.Y = r.Min.Y + rand.Int()%r.Dy()
	return
}

func clearAndPrint(a *pixd.Area, face *font.Face, s string) {
	t1 := time.Now()
	a.SetColor(white)
	r := a.Bounds()
	a.Fill(r)
	a.SetColorRGBA(80, 160, 240, 255)
	p0, p1, p2 := rndPt(r), rndPt(r), rndPt(r)
	a.Quad(p0, p0, p1, p2, true)
	a.SetColor(black)
	w := a.NewTextWriter(face)
	w.WriteString(s)
	t2 := time.Now()
	println("DrawText:", t2.Sub(t1).String())
	time.Sleep(10 * time.Second)
}

func DrawText(disp *pixd.Display, n int) {
	a := disp.NewArea(disp.Bounds())
	for i := 0; i != n; i++ {
		clearAndPrint(a, fdejavu, akkermanSteppesEN)
		clearAndPrint(a, fdejavu, akkermanSteppesPL)
		clearAndPrint(a, fvga, akkermanSteppesEN)
		clearAndPrint(a, fvga, akkermanSteppesPL)
	}
}
