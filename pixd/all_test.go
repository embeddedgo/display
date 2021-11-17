// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd_test

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/embeddedgo/display/pixd"
	"github.com/embeddedgo/display/pixd/driver/imgdrv"
	"github.com/embeddedgo/display/pixd/font"
	"github.com/embeddedgo/display/pixd/font/font9/anonpro11"
	"github.com/embeddedgo/display/pixd/font/font9/dejavu12"
	"github.com/embeddedgo/display/pixd/font/font9/vga"
)

var dir = filepath.Join(os.TempDir(), "pix_test")

func failErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestDrawGeom(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 200, 410))

	disp1 := pixd.NewDisplay(imgdrv.New(screen.SubImage(image.Rect(0, 0, 200, 200)).(*image.NRGBA)))
	disp2 := pixd.NewDisplay(imgdrv.New(screen.SubImage(image.Rect(0, 210, 200, 410)).(*image.NRGBA)))
	disp2.SetOrigin(image.Pt(0, 200))

	a := pixd.NewArea(image.Rect(0, 0, 200, 400), disp1, disp2)
	a.SetColor(color.Gray{220})
	a.Fill(a.Bounds())

	max := a.Bounds().Max

	a.SetColorRGBA(24, 46, 68, 255)
	for x := 0; x < max.X; x += 2 {
		a.Pixel(x, 0)
		a.Pixel(x+1, max.Y-1)
	}
	for y := 2; y < max.Y; y += 2 {
		a.Pixel(0, y)
		a.Pixel(max.X-1, y-1)
	}

	x := max.X / 2
	xl := max.X / 4
	xr := max.X * 3 / 4
	for r := 0; r < 19; r++ {
		y := 3 + (r+2)*r
		a.SetColorRGBA(0, 100, 200, 255)
		a.Ellipse(image.Pt(x, y), r, r, true)
		a.Ellipse(image.Pt(xl, y), r, r/2, true)
		a.Ellipse(image.Pt(xr, y), r/2, r, true)
		a.SetColorRGBA(100, 50, 0, 255)
		a.Ellipse(image.Pt(x, y), r, r, false)
		a.Ellipse(image.Pt(xl, y), r, r/2, false)
		a.Ellipse(image.Pt(xr, y), r/2, r, false)
	}

	a.SetColorRGBA(250, 100, 0, 255)
	for i := 0; i < 20; i++ {
		x := 2*i + 4
		y := i*i + 2
		a.Line(image.Pt(2, y), image.Pt(x, 2))
		a.Line(image.Pt(max.X-1-2, y), image.Pt(max.X-1-x, 2))
	}

	f, err := os.OpenFile(filepath.Join(dir, "geom.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}

func TestDrawImage(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 40, 40))
	disp := pixd.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds().Inset(4))
	a.SetColorRGBA(0, 0, 128, 255)
	a.Fill(a.Bounds())

	img := pixd.NewAlphaN(image.Rect(0, 0, 11, 11), 1)
	img.Set(0, 10, color.Alpha{1})
	img.Set(2, 8, color.Gray{1})
	img.Set(4, 6, color.Gray16{1})
	img.SetAlpha(6, 4, color.Alpha{1})
	img.Set(8, 2, color.RGBA{0, 0, 0, 1})
	img.Set(10, 0, color.RGBA64{0, 0, 0, 1})
	img.SetAlpha(4, 4, color.Alpha{1})

	a.Draw(disp.Bounds(), img, image.Point{}, nil, image.Point{}, draw.Over)
	a.Draw(disp.Bounds().Add(image.Pt(20, 25)), img, image.Point{}, nil, image.Point{}, draw.Over)

	imm := pixd.NewImmAlphaN(img.Bounds(), 1, string(img.Pix))

	a.Draw(disp.Bounds().Add(image.Pt(5, 5)),
		&image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, // source
		imm, image.Point{}, // mask
		draw.Over,
	)
	a.Draw(disp.Bounds().Add(image.Pt(10, 10)),
		&image.Uniform{color.NRGBA{255, 0, 0, 150}}, image.Point{}, // source
		imm, image.Point{}, // mask
		draw.Over,
	)

	imm = imm.SubImage(image.Rect(2, 2, 11, 11)).(*pixd.ImmAlphaN)
	a.Draw(disp.Bounds().Add(image.Pt(16, 16)), imm, image.Pt(2, 2), nil, image.Point{}, draw.Src)

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
			&dejavu12.X0101_0201,
		},
	}
	AnonPro11 = &font.Face{
		Height: anonpro11.Height,
		Ascent: anonpro11.Ascent,
		Subfonts: []*font.Subfont{
			&anonpro11.X0000_0100,
			&anonpro11.X0101_0201,
		},
	}
	VGA = &font.Face{
		Height: vga.Height,
		Ascent: vga.Ascent,
		Subfonts: []*font.Subfont{
			&vga.X0000_007f,
			&vga.X00a0_021f,
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

	screen := image.NewNRGBA(image.Rect(0, 0, 440, 780))
	disp := pixd.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(250, 250, 200, 255)
	a.Fill(a.Bounds())
	a.SetRect(a.Rect().Inset(4))
	a.SetColorRGBA(0, 0, 100, 255)

	w := a.NewTextWriter(Dejavu12)
	w.WriteString(AkermanianSteppes)

	w.Face = AnonPro11
	w.WriteString(AkermanianSteppes)

	w.Face = VGA
	w.WriteString(AkermanianSteppes)

	f, err := os.OpenFile(filepath.Join(dir, "font.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}

func TestTriangle(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 400, 800))
	disp := pixd.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(0, 0, 0, 255)
	a.Fill(a.Bounds())

	p0 := image.Pt(100, 10)
	p1 := image.Pt(380, 200)
	p2 := image.Pt(10, 210)
	a.SetColorRGBA(255, 0, 0, 255)
	a.Triangle(p0, p1, p2, true)
	a.SetColorRGBA(0, 192, 0, 192)
	a.Triangle(p0, p1, p2, false)

	p0 = image.Pt(100, 220)
	p1 = image.Pt(380, 300)
	p2 = image.Pt(10, 790)
	a.SetColorRGBA(255, 0, 0, 255)
	a.Triangle(p0, p1, p2, true)
	a.SetColorRGBA(0, 192, 0, 192)
	a.Triangle(p0, p1, p2, false)

	p0 = image.Pt(390, 320)
	p1 = image.Pt(300, 790)
	p2 = image.Pt(100, 700)
	a.SetColorRGBA(255, 0, 0, 255)
	a.Triangle(p0, p1, p2, true)
	a.SetColorRGBA(0, 192, 0, 192)
	a.Triangle(p0, p1, p2, false)

	f, err := os.OpenFile(filepath.Join(dir, "triangle.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}
