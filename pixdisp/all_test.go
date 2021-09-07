// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp_test

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/embeddedgo/display/pixdisp"
	"github.com/embeddedgo/display/pixdisp/driver/imgdrv"
	"github.com/embeddedgo/display/pixdisp/font"
	"github.com/embeddedgo/display/pixdisp/font/font9/anonpro11"
	"github.com/embeddedgo/display/pixdisp/font/font9/dejavu12"
	"github.com/embeddedgo/display/pixdisp/font/font9/vga"
)

var dir = filepath.Join(os.TempDir(), "pixdisp_test")

func failErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestDrawGeom(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 91, 400))
	disp := pixdisp.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds().Inset(5))
	a.SetColor(color.Gray{220})
	a.Fill(a.Bounds())

	max := a.Bounds().Max

	a.SetColor(pixdisp.RGB16(0x1234))
	for x := 0; x < max.X; x += 2 {
		a.DrawPoint(image.Pt(x, 0), 0)
		a.DrawPoint(image.Pt(x+1, max.Y-1), 0)
	}
	for y := 2; y < max.Y; y += 2 {
		a.DrawPoint(image.Pt(0, y), 0)
		a.DrawPoint(image.Pt(max.X-1, y-1), 0)
	}

	x := max.X / 2
	for r := 0; r < 19; r++ {
		y := 3 + (r+2)*r
		a.SetColorRGB(0, 100, 200)
		a.DrawPoint(image.Pt(x, y), r)
		a.SetColor(color.RGBA{100, 50, 0, 255})
		a.DrawCircle(image.Pt(x, y), r)
	}

	a.SetColorRGB(250, 100, 0)
	for i := 0; i < 17; i++ {
		x := 2*i + 4
		y := i*i + 2
		a.DrawLine(image.Pt(2, y), image.Pt(x, 2))
		a.DrawLine(image.Pt(max.X-1-2, y), image.Pt(max.X-1-x, 2))
	}

	f, err := os.OpenFile(filepath.Join(dir, "geom.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}

func TestDrawImage(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 40, 40))
	disp := pixdisp.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds().Inset(4))
	a.SetColor(pixdisp.RGB{0, 0, 128})
	a.Fill(a.Bounds())

	img := pixdisp.NewAlphaN(image.Rect(0, 0, 11, 11), 1)
	img.Set(0, 10, color.Alpha{1})
	img.Set(2, 8, color.Gray{1})
	img.Set(4, 6, color.Gray16{1})
	img.SetAlpha(6, 4, color.Alpha{1})
	img.Set(8, 2, color.RGBA{0, 0, 0, 1})
	img.Set(10, 0, color.RGBA64{0, 0, 0, 1})
	img.SetAlpha(4, 4, color.Alpha{1})

	a.Draw(disp.Bounds(), img, image.Pt(0, 0), draw.Over)
	a.Draw(disp.Bounds().Add(image.Pt(20, 25)), img, image.Pt(0, 0), draw.Over)

	imm := pixdisp.NewImmAlphaN(img.Bounds(), 1, string(img.Pix))

	a.DrawMask(disp.Bounds().Add(image.Pt(5, 5)),
		&image.Uniform{pixdisp.RGB{255, 0, 0}}, image.Pt(0, 0), // source
		imm, image.Pt(0, 0), // mask
		draw.Over,
	)
	a.DrawMask(disp.Bounds().Add(image.Pt(10, 10)),
		&image.Uniform{color.NRGBA{255, 0, 0, 150}}, image.Pt(0, 0), // source
		imm, image.Pt(0, 0), // mask
		draw.Over,
	)

	//width := float32(7)
	//a.PlotLine(0, 0, a.Bounds().Max.X, a.Bounds().Max.Y/2, width)
	//a.PlotLine(0, a.Bounds().Max.Y-4, a.Bounds().Max.X, a.Bounds().Max.Y-4, width)

	imm = imm.SubImage(image.Rect(2, 2, 11, 11)).(*pixdisp.ImmAlphaN)
	a.Draw(disp.Bounds().Add(image.Pt(16, 16)), imm, image.Pt(2, 2), draw.Src)

	f, err := os.OpenFile(filepath.Join(dir, "image.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}

var (
	Dejavu12 = &font.Face{
		Height: dejavu12.Height,
		Ascent: dejavu12.Ascent,
		Subfonts: []*font.Subfont{
			&dejavu12.X0000_0100,
		},
	}
	AnonPro11 = &font.Face{
		Height: anonpro11.Height,
		Ascent: anonpro11.Ascent,
		Subfonts: []*font.Subfont{
			&anonpro11.X0000_0100,
		},
	}
	VGA = &font.Face{
		Height: vga.Height,
		Ascent: vga.Ascent,
		Subfonts: []*font.Subfont{
			&vga.X0000_007f,
		},
	}
)

const AkermanianSteppes = `
Wpłynąłem na suchego przestwór oceanu,
Wóz nurza się w zieloność i jak łódka brodzi,
Śród fali łąk szumiących, śród kwiatów powodzi,
Omijam koralowe ostrowy burzanu.

Już mrok zapada, nigdzie drogi ni kurhanu;
Patrzę w niebo, gwiazd szukam, przewodniczek łodzi;
Tam z dala błyszczy obłok - tam jutrzenka wschodzi;
To błyszczy Dniestr, to weszła lampa Akermanu.

Stójmy! - jak cicho! - słyszę ciągnące żurawie,
Których by nie dościgły źrenice sokoła;
Słyszę, kędy się motyl kołysa na trawie,

Kędy wąż śliską piersią dotyka się zioła.
W takiej ciszy - tak ucho natężam ciekawie,
Że słyszałbym głos z Litwy. - Jedźmy, nikt nie woła.
`

func TestFont(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 480, 800))
	disp := pixdisp.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds())
	a.SetColor(pixdisp.RGB{250, 250, 200})
	a.Fill(a.Bounds())
	a.SetRect(a.Rect().Inset(4))
	a.SetColorRGB(0, 0, 100)

	w := a.TextWriter(Dejavu12)
	w.WriteString(AkermanianSteppes)

	f, err := os.OpenFile(filepath.Join(dir, "font.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}
