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

	"github.com/embeddedgo/display/math2d"
	"github.com/embeddedgo/display/pixd"
	"github.com/embeddedgo/display/pixd/driver/imgdrv"
	"github.com/embeddedgo/display/pixd/font"
	"github.com/embeddedgo/display/pixd/font/font9/anonpro11"
	"github.com/embeddedgo/display/pixd/font/font9/dejavu12"
	"github.com/embeddedgo/display/pixd/font/font9/vga"
)

var dir = filepath.Join(os.TempDir(), "pixd_test")

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
		a.RoundRect(image.Pt(x, y), image.Pt(x, y), r, r, true)
		a.RoundRect(image.Pt(xl, y), image.Pt(xl, y), r, r/2, true)
		a.RoundRect(image.Pt(xr, y), image.Pt(xr, y), r/2, r, true)
		a.SetColorRGBA(100, 50, 0, 255)
		a.RoundRect(image.Pt(x, y), image.Pt(x, y), r, r, false)
		a.RoundRect(image.Pt(xl, y), image.Pt(xl, y), r, r/2, false)
		a.RoundRect(image.Pt(xr, y), image.Pt(xr, y), r/2, r, false)
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

const AkermanianSteppesPL = `
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

const AkermanianSteppesEN = `
The Akkerman Steppes Original Polish (above) by Adam Mickiewicz
(1798-1855) translated to English by Leo Yankevich (below).

I launch myself across the dry and open narrows,
My carriage plunging into green as if a ketch,
Floundering through the meadow flowers in the stretch.
I pass an archipelago of coral yarrows.

It's dusk now, not a road in sight, nor ancient barrows.
I look up at the sky and look for stars to catch.
There distant clouds glint-there tomorrow starts to etch;
The Dnieper glimmers; Akkerman's lamp shines and harrows.

I stand in stillness, hear the migratory cranes,
Their necks and wings beyond the reach of preying hawks;
Hear where the sooty copper glides across the plains,

Where on its underside a viper writhes through stalks.
Amid the hush I lean my ears down grassy lanes
And listen for a voice from home. Nobody talks.
`

func TestFont(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 440, 810))
	disp := pixd.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(250, 250, 200, 255)
	a.Fill(a.Bounds())
	a.SetRect(a.Rect().Inset(4))
	a.SetColorRGBA(0, 0, 100, 255)

	w := a.NewTextWriter(Dejavu12)
	w.WriteString(AkermanianSteppesPL)

	w.Face = AnonPro11
	w.WriteString(AkermanianSteppesEN)

	w.Face = VGA
	w.WriteString(AkermanianSteppesPL)

	f, err := os.OpenFile(filepath.Join(dir, "font.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}

func TestRectTriangle(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 400, 800))
	disp := pixd.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(0, 0, 0, 255)
	a.Fill(a.Bounds())

	triangles := [][3]image.Point{
		{{10, 10}, {300, 8}, {390, 10}},
		{{10, 20}, {300, 19}, {390, 20}},
		{{10, 30}, {300, 30}, {390, 30}},

		{{150, 40}, {250, 120}, {350, 200}},
		{{160, 40}, {240, 80}, {320, 120}},

		{{190, 40}, {190, 40}, {190, 40}},
		{{200, 40}, {202, 41}, {201, 42}},
		{{210, 40}, {212, 40}, {211, 41}},

		{{10, 600}, {200, 600}, {390, 600}},
		{{10, 610}, {200, 611}, {390, 610}},
		{{10, 620}, {200, 622}, {390, 620}},
		{{10, 630}, {200, 633}, {390, 630}},
		{{10, 640}, {200, 644}, {390, 640}},
		{{10, 650}, {200, 655}, {390, 650}},
		{{10, 660}, {200, 666}, {390, 660}},
		{{10, 670}, {200, 677}, {390, 670}},
		{{10, 680}, {200, 688}, {390, 680}},
		{{10, 690}, {200, 699}, {390, 690}},
	}

	for _, tr := range triangles {
		a.SetColorRGBA(0, 192, 0, 192)
		a.Quad(tr[0], tr[1], tr[2], tr[0], true)
		a.SetColorRGBA(192, 0, 0, 192)
		a.FillQuad(tr[0], tr[1], tr[2], tr[0])
	}

	triangles = [][3]image.Point{
		{{100, 40}, {380, 240}, {10, 250}},
		{{100, 260}, {380, 310}, {10, 590}},
		{{390, 500}, {380, 310}, {10, 590}},
	}

	for _, tr := range triangles {
		//a.SetColorRGBA(255, 0, 0, 255)
		//a.Triangle(tr[0], tr[1], tr[2], true)
		a.SetColorRGBA(192, 0, 0, 192)
		a.FillQuad(tr[0], tr[1], tr[2], tr[2])
	}

	a.SetColorRGBA(0, 0, 128, 128)
	a.RoundRect(image.Pt(200, 100), image.Pt(300, 400), 30, 40, true)
	a.SetColorRGBA(0, 128, 0, 128)
	a.RoundRect(image.Pt(200, 100), image.Pt(300, 400), 30, 40, false)

	a.SetColorRGBA(0, 0, 128, 128)
	a.RoundRect(image.Pt(100, 300), image.Pt(350, 400), 20, 20, true)
	a.SetColorRGBA(0, 128, 0, 128)
	a.RoundRect(image.Pt(100, 300), image.Pt(350, 400), 20, 20, false)

	f, err := os.OpenFile(filepath.Join(dir, "triangle.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}

func TestQuad(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 400, 800))
	disp := pixd.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(0, 0, 0, 255)
	a.Fill(a.Bounds())

	quads := [][4]image.Point{
		{{40, 240}, {230, 300}, {300, 20}, {10, 10}},
		{{320, 10}, {380, 10}, {390, 300}, {250, 300}},
	}

	for _, q := range quads {
		a.SetColorRGBA(192, 0, 0, 192)
		a.Quad(q[0], q[1], q[2], q[3], true)
		a.SetColorRGBA(0, 192, 0, 192)
		a.Quad(q[0], q[1], q[2], q[3], false)
	}

	quads = [][4]image.Point{
		{{10, 310}, {100, 310}, {110, 400}, {30, 400}},
		{{120, 310}, {200, 320}, {200, 400}, {120, 390}},
		{{220, 320}, {320, 320}, {320, 400}, {220, 400}},
	}
	for _, q := range quads {
		a.SetColorRGBA(0, 192, 0, 192)
		a.Quad(q[0], q[1], q[2], q[3], true)
		a.SetColorRGBA(192, 0, 0, 192)
		a.FillQuad(q[0], q[1], q[2], q[3])
	}

	f, err := os.OpenFile(filepath.Join(dir, "quad.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}

func TestArc(t *testing.T) {
	os.Mkdir(dir, 0755)

	screen := image.NewNRGBA(image.Rect(0, 0, 400, 800))
	disp := pixd.NewDisplay(imgdrv.New(screen))

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(0, 0, 0, 255)
	a.Fill(a.Bounds())

	a.SetColorRGBA(255, 255, 255, 255)
	th0 := math2d.RightAngle / 3
	th1 := math2d.RightAngle
	a.Arc(image.Pt(200, 200), 50, 50, 100, 100, th0, th1, true)

	f, err := os.OpenFile(filepath.Join(dir, "arc.png"), os.O_WRONLY|os.O_CREATE, 0755)
	failErr(t, err)
	failErr(t, png.Encode(f, screen))
	failErr(t, f.Close())
}
