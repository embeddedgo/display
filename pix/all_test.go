// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix_test

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/embeddedgo/display/font"
	"github.com/embeddedgo/display/font/subfont/font9/anonpro11"
	"github.com/embeddedgo/display/font/subfont/font9/dejavu12"
	"github.com/embeddedgo/display/font/subfont/font9/dejavu14"
	"github.com/embeddedgo/display/font/subfont/font9/dejavusans14"
	"github.com/embeddedgo/display/font/subfont/font9/terminus12"
	"github.com/embeddedgo/display/font/subfont/font9/vga"
	"github.com/embeddedgo/display/images"
	"github.com/embeddedgo/display/math2d"
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/imgdrv"
)

const testDir = "testdata"

var (
	white = color.Gray{255}
	black = color.Gray{0}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}

	AnonPro11 = anonpro11.NewFace(
		anonpro11.X0000_0100,
		anonpro11.X0101_0201,
	)
	Dejavu12 = dejavu12.NewFace(
		dejavu12.X0000_0100,
		dejavu12.X0101_0201,
		dejavu12.X0404_0504,
	)
	Dejavu14 = dejavu14.NewFace(
		dejavu14.X0000_0100,
		dejavu14.X0101_0201,
		dejavusans14.X0404_0504,
	)
	Terminus12 = terminus12.NewFace(
		terminus12.X0020_007e,
		terminus12.X00a0_0175,
		terminus12.X0178_017f,
	)
	VGA = vga.NewFace(
		vga.X0000_007f,
		vga.X00a0_021f,
	)
)

func failErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func newDisplay(width, height int) *pix.Display {
	screen := image.NewNRGBA(image.Rect(0, 0, width, height))
	return pix.NewDisplay(imgdrv.New(screen))
}

func loadFile(t *testing.T, name string) []byte {
	data, err := os.ReadFile(name)
	failErr(t, err)
	return data
}

func loadImage(t *testing.T, name string) image.Image {
	f, err := os.Open(name)
	failErr(t, err)
	img, _, err := image.Decode(f)
	f.Close()
	failErr(t, err)
	return img
}

func checkImage(t *testing.T, img image.Image, name string) {
	saved := loadImage(t, filepath.Join(testDir, name))
	r0, r1 := saved.Bounds(), img.Bounds()
	if r0.Size() != r1.Size() {
		t.Error(name, "different sizes:", r0.Size(), r1.Size())
		return
	}
	size := r0.Size()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			c0 := img.At(r0.Min.X+x, r0.Min.Y+y)
			c1 := saved.At(r1.Min.X+x, r1.Min.Y+y)
			r0, g0, b0, a0 := c0.RGBA()
			r1, g1, b1, a1 := c1.RGBA()
			if r0 != r1 || g0 != g1 || b0 != b1 || a0 != a1 {
				t.Error(name, "different pixel color:", c0, c1, "at", x, y)
				return
			}
		}
	}
}

func checkDisplay(t *testing.T, disp *pix.Display, name string) {
	checkImage(t, disp.Driver().(*imgdrv.Driver).Image, name)
}

func saveImage(t *testing.T, img image.Image, name string) {
	f, err := os.Create(filepath.Join(testDir, name))
	failErr(t, err)
	failErr(t, png.Encode(f, img))
	failErr(t, f.Close())
}

func saveDisplay(t *testing.T, disp *pix.Display, name string) {
	saveImage(t, disp.Driver().(*imgdrv.Driver).Image, name)
}

func TestLineEllipse(t *testing.T) {
	testFile := "line_ellipse.png"

	img := image.NewNRGBA(image.Rect(0, 0, 300, 400))

	disp1 := pix.NewDisplay(imgdrv.New(img.SubImage(image.Rect(0, 0, 300, 200)).(*image.NRGBA)))
	disp2 := pix.NewDisplay(imgdrv.New(img.SubImage(image.Rect(0, 210, 300, 400)).(*image.NRGBA)))
	disp2.SetOrigin(image.Pt(0, 200))

	a := pix.NewArea(image.Rect(0, 0, 300, 390), disp1, disp2)
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
		p := image.Pt(x-r-1, y)
		a.RoundRect(p, p, r, r, true)
		p = image.Pt(xl-r-1, y)
		a.RoundRect(p, p, r, r/2, true)
		p = image.Pt(xr-r-1, y)
		a.RoundRect(p, p, r/2, r, true)

		a.SetColorRGBA(100, 50, 0, 255)
		p = image.Pt(x+r+1, y)
		a.RoundRect(p, p, r, r, false)
		p = image.Pt(xl+r+1, y)
		a.RoundRect(p, p, r, r/2, false)
		p = image.Pt(xr+r+1, y)
		a.RoundRect(p, p, r/2, r, false)
	}

	a.SetColorRGBA(250, 100, 0, 255)
	for i := 0; i < 20; i++ {
		x := 2*i + 4
		y := i*i + 2
		a.Line(image.Pt(2, y), image.Pt(x, 2))
		a.Line(image.Pt(max.X-1-2, y), image.Pt(max.X-1-x, 2))
	}

	//saveImage(t, img, testFile)
	checkImage(t, img, testFile)
}

func TestImage(t *testing.T) {
	testFile := "image.png"

	gobug := loadImage(t, "../testdata/gopherbug.jpg")
	gobugMask := new(images.AlphaN)
	gobugMask.Rect.Min = image.Pt(31, 21)
	gobugMask.Rect.Max = gobugMask.Rect.Min.Add(image.Pt(211, 251))
	gobugMask.Stride = 27
	gobugMask.Pix = loadFile(t, "../testdata/gopherbug.211x251s27b1")

	disp := newDisplay(272*2, 272*2)
	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(0, 0, 128, 255)
	a.Fill(a.Bounds())

	r := gobug.Bounds()
	p00 := image.Point{}
	r = r.Sub(r.Min)
	p := image.Pt(0, 0)
	a.Draw(r.Add(p), gobug, p00, nil, p00, draw.Src)
	p = image.Pt(272, 0)
	a.Draw(r.Add(p), gobugMask, p00, nil, p00, draw.Src)
	p = image.Pt(272, 272)
	a.Draw(r.Add(p), gobug, p00, gobugMask, p00, draw.Over)
	for i, b := range gobugMask.Pix {
		gobugMask.Pix[i] = ^b
	}
	p = image.Pt(0, 272)
	a.Draw(r.Add(p), gobug, p00, gobugMask, p00, draw.Over)

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

// The Akkerman Steppe Original Polish by Adam Mickiewicz (1798-1855)
const AkermanianSteppePL = `` +
	`Wpłynąłem na suchego przestwór oceanu,
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

// The Akkerman Steppe translated to English by Leo Yankevich
const AkermanianSteppeEN = `
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

const AkermanianSteppeDE = `
In kargen Raum gedrungen, ozeanisch weiten,
taucht mein Wagen ein, ein schwerer Kahn, gezogen
durch das Blütenmeer, rauschende Wiesenwogen,
weicht Inseln, Riffen aus, muß mit den Stürmen streiten.

Ein Dämmern senkt sich, doch kein Stern will mich geleiten
hab ich den Himmel nach Vertrautem überflogen.
Was glimmt dort? Zieht der Morgenstern schon seinen Bogen?
Als Morgengruß woll'n Lichter übern Dnjestr gleiten.

Sei still! Hoch ziehn die Kraniche in langen Ketten,
So hoch daß sie selbst wachem Falkenblick entgehen,
Ich hör, sich Schmetterlinge in die Winde betten.

Schlüpft hier ` + "`" + `ne Schlange durch das Gras? ... Die Ähren wehen...
Ich horche weit, auf Stimmen aus vertrauten Stätten;
aus Litau'n...  Weiter! ... niemand ist zu hör'n, zu sehen.
`

func TestTextWriter(t *testing.T) {
	testFile := "text_writer.png"

	disp := newDisplay(460, 1300)

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(250, 250, 200, 255)
	a.Fill(a.Bounds())
	a.SetRect(a.Rect().Inset(1))
	a.SetColorRGBA(0, 0, 100, 255)

	w := a.NewTextWriter(Dejavu12)
	w.WriteString(AkermanianSteppePL)

	w.Face = VGA
	w.WriteString(AkermanianSteppeEN)

	w.Face = font.Magnify(Terminus12, 2, 2, font.Nearest)
	w.WriteString(AkermanianSteppeDE)

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

func TestTriangle(t *testing.T) {
	testFile := "triangle.png"

	disp := newDisplay(400, 710)

	a := disp.NewArea(disp.Bounds())
	a.SetColor(black)
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
		a.Quad(tr[0], tr[1], tr[1], tr[2], true)
		a.SetColorRGBA(192, 0, 0, 192)
		a.FillQuad(tr[0], tr[1], tr[2], tr[0])
	}

	triangles = [][3]image.Point{
		{{100, 40}, {100, 240}, {10, 250}},
		{{100, 40}, {380, 260}, {100, 240}},
		{{100, 260}, {380, 310}, {10, 590}},
		{{390, 500}, {380, 310}, {10, 590}},
	}

	for i, tr := range triangles {
		if i&1 == 0 {
			a.SetColorRGBA(192, 0, 0, 192)
		} else {
			a.SetColorRGBA(0, 0, 192, 192)
		}
		a.FillQuad(tr[0], tr[0], tr[1], tr[2])
	}

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

func TestRoundRect(t *testing.T) {
	testFile := "roundrect.png"

	disp := newDisplay(320, 420)

	a := disp.NewArea(disp.Bounds())
	a.SetOrigin(image.Pt(60, 50))
	a.SetColor(black)
	a.Fill(a.Bounds())

	o := image.Pt(20, 20)
	p0 := image.Pt(150, 100)
	p1 := image.Pt(250, 400)
	ra, rb := 30, 40
	a.SetColorRGBA(0, 0, 128, 128)
	a.RoundRect(p0, p1, ra, rb, true)
	a.SetColorRGBA(0, 128, 0, 128)
	a.RoundRect(p0.Add(o), p1.Add(o), ra, rb, false)

	p0 = image.Pt(100, 300)
	p1 = image.Pt(350, 400)
	ra, rb = 15, 15
	a.SetColorRGBA(128, 0, 0, 128)
	a.RoundRect(p0, p1, ra, rb, true)
	a.SetColorRGBA(0, 128, 0, 128)
	a.RoundRect(p0.Sub(o), p1.Sub(o), ra, rb, false)

	r := a.Bounds().Inset(10)
	r.Max.X--
	r.Max.Y--
	a.SetColorRGBA(200, 0, 200, 200)
	a.RoundRect(r.Min, r.Max, 10, 10, false)

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

func TestQuad(t *testing.T) {
	testFile := "quad.png"

	disp := newDisplay(400, 420)

	a := disp.NewArea(disp.Bounds())
	a.SetColor(black)
	a.Fill(a.Bounds())

	quads := [][4]image.Point{
		{{40, 240}, {230, 300}, {300, 20}, {10, 10}},
		{{320, 10}, {380, 10}, {390, 300}, {250, 300}},
		//{{20, 10}, {30, 30}, {35, 60}, {25, 35}},
	}

	for _, q := range quads {
		a.SetColorRGBA(0, 255, 0, 255)
		a.Quad(q[0], q[1], q[2], q[3], false)
		a.SetColorRGBA(191, 0, 0, 191)
		a.Quad(q[0], q[1], q[2], q[3], true)
	}

	quads = [][4]image.Point{
		{{10, 310}, {100, 310}, {110, 400}, {30, 400}},
		{{120, 310}, {200, 320}, {200, 400}, {120, 390}},
		{{220, 320}, {320, 320}, {320, 400}, {220, 400}},
	}
	for _, q := range quads {
		a.SetColorRGBA(0, 200, 0, 200)
		a.Quad(q[0], q[1], q[2], q[3], true)
		a.SetColorRGBA(192, 0, 0, 192)
		a.FillQuad(q[0], q[1], q[2], q[3])
	}
	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

func testArc(t *testing.T, fill bool) {
	testFile := "arc0.png"
	if fill {
		testFile = "arc1.png"
	}

	disp := newDisplay(400, 900)

	a := disp.NewArea(disp.Bounds())
	a.SetColor(black)
	a.Fill(a.Bounds())
	a.SetColor(white)
	a.RoundRect(a.Bounds().Min, a.Bounds().Max.Sub(image.Pt(1, 1)), 0, 0, false)

	th0 := math2d.RightAngle / 3
	th1 := math2d.RightAngle

	a.SetColorRGBA(210, 0, 0, 210)
	a.Arc(image.Pt(200, 10), 70, 50, 140, 100, th0, th1, fill)
	a.SetColorRGBA(0, 0, 150, 150)
	a.Arc(image.Pt(200, 10), 70, 50, 140, 100, th1, th1+1e8, fill)

	th1 += math2d.RightAngle

	a.SetColorRGBA(210, 0, 0, 210)
	a.Arc(image.Pt(200, 150), 70, 50, 140, 100, th0, th1, fill)
	a.SetColorRGBA(0, 0, 150, 150)
	a.Arc(image.Pt(200, 150), 70, 50, 140, 100, th0-1e8, th0, fill)

	th1 += math2d.RightAngle - 1e8

	a.SetColorRGBA(210, 0, 0, 210)
	a.Arc(image.Pt(200, 390), 75, 55, 135, 95, th0, th1, fill)
	a.SetColorRGBA(0, 0, 150, 150)
	a.Arc(image.Pt(200, 390), 75, 55, 135, 95, th1, th0, fill)
	a.SetColorRGBA(0, 100, 0, 100)
	a.Arc(image.Pt(200, 390), 70, 50, 75, 55, 0, 0, fill)
	a.Arc(image.Pt(200, 390), 135, 95, 140, 100, 0, 0, fill)

	th1 += math2d.RightAngle / 2

	a.SetColorRGBA(210, 0, 0, 210)
	a.Arc(image.Pt(200, 630), 70, 50, 140, 100, th0, th1, fill)
	a.SetColorRGBA(0, 0, 150, 150)
	a.Arc(image.Pt(200, 630), 70, 50, 140, 100, th1, th0, fill)
	a.SetColorRGBA(200, 200, 0, 255)
	a.Arc(image.Pt(200, 630), 60, 40, 61, 41, 0, 0, fill)
	a.Arc(image.Pt(200, 630), 150, 110, 151, 111, 0, 0, fill)

	th0 = -math2d.FullAngle / 2

	a.SetColorRGBA(210, 0, 0, 210)
	a.Arc(image.Pt(200, 870), 70, 50, 140, 100, th0, th1, fill)
	a.SetColorRGBA(0, 0, 150, 150)
	a.Arc(image.Pt(200, 870), 70, 50, 140, 100, th1, th1+1e8, fill)

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

func TestArc(t *testing.T) {
	testArc(t, true)
	testArc(t, false)
}

func TestDispRect(t *testing.T) {
	testFile := "disp_rect.png"

	disp := newDisplay(240, 320)
	disp.SetRect(image.Rect(0, 40, 240, 240+40))

	a := disp.NewArea(disp.Bounds())
	a.SetColorRGBA(0, 100, 0, 255)
	a.Fill(a.Bounds())

	r := a.Bounds()
	r.Max = r.Max.Sub(image.Pt(1, 1))
	a.SetColorRGBA(255, 0, 0, 255)
	a.RoundRect(r.Min, r.Max, 0, 0, false)
	r = r.Inset(1)
	a.SetColor(white)
	a.Line(r.Min, r.Max)
	p0 := r.Min
	p1 := p0.Add(image.Pt(10, 6))
	p2 := p0.Add(image.Pt(6, 10))
	a.Quad(p0, p0, p1, p2, true)

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

func button(a *pix.Area, center image.Point, s string, f font.Face) {
	width := pix.StringWidth(s, f)
	height, ascent := f.Size()
	p0 := center
	p0.X -= (width + 1) / 2
	p0.Y -= (height + 1) / 2
	p1 := p0.Add(image.Pt(width-1, height-1))
	a.SetColorRGBA(200, 200, 200, 200)
	a.RoundRect(p0, p1, 5, 5, true)
	a.SetColor(black)
	a.RoundRect(p0, p1, 5, 5, false)
	p0.Y += ascent
	w := &pix.TextWriter{
		Area:  a,
		Face:  f,
		Color: &image.Uniform{black},
		Pos:   p0,
	}
	w.WriteString(s)
}

func TestButtons(t *testing.T) {
	testFile := "buttons.png"

	disp := newDisplay(100, 200)

	a := disp.NewArea(disp.Bounds())
	r := a.Bounds()
	a.SetColorRGBA(200, 220, 250, 255)
	a.Fill(r)
	a.SetColorRGBA(240, 180, 140, 255)
	a.Quad(image.Pt(20, 80), image.Pt(70, 10), image.Pt(80, 100), image.Pt(40, 190), true)

	p := image.Pt(r.Max.X/2, r.Max.Y/4)
	button(a, p, "Accept", AnonPro11)
	p.Y = r.Max.Y / 2
	button(a, p, "Witaj Świecie!", Dejavu12)
	p.Y = r.Max.Y * 3 / 4
	button(a, p, " OK ", VGA)

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}

var helloWorld = []string{
	"Hello, Wolrd!\n",
	"Witaj, Świecie!\n",
	"Привет, Мир!\n",
}

func TestTextRotation(t *testing.T) {
	testFile := "text_rotation.png"

	img := loadImage(t, "../testdata/gopherbug.jpg")
	img = images.Magnify(img, 1, 2, images.Bilinear)
	size := img.Bounds().Size()
	disp := newDisplay(size.X, size.Y)
	a := disp.NewArea(disp.Bounds())
	a.Draw(a.Bounds(), img, img.Bounds().Min, nil, image.Point{}, draw.Src)

	a.SetColor(red)
	r := a.Bounds()
	hc := (r.Min.X + r.Max.X - 1) / 2
	vc := (r.Min.Y + r.Max.Y - 1) / 2
	a.Line(image.Pt(hc, r.Min.Y), image.Pt(hc, r.Max.Y-1))
	a.Line(image.Pt(r.Min.X, vc), image.Pt(r.Max.X-1, vc))

	a.SetColor(blue)
	w := a.NewTextWriter(Dejavu14)

	markNextGlyph := func() {
		p0 := w.Pos.Add(w.Offset)
		p1 := p0.Add(image.Pt(6, 6))
		p2 := p0.Add(image.Pt(0, 9))
		c := a.Color()
		a.SetColor(red)
		a.Quad(p0, p0, p1, p2, true)
		a.SetColor(c)
	}

	r = disp.Bounds()
	r.Min.X = (r.Min.X+r.Max.X)/2 + 4
	a.SetRect(r)
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()
	r = disp.Bounds()
	r.Max.X = (r.Min.X+r.Max.X)/2 - 5
	a.SetRect(r)
	a.SetMirror(pix.MX)
	w.Pos = a.Bounds().Min
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()

	a.SetMirror(pix.MY)
	w.Pos = a.Bounds().Min
	r = disp.Bounds()
	r.Min.X = (r.Min.X+r.Max.X)/2 + 4
	a.SetRect(r)
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()
	r = disp.Bounds()
	r.Max.X = (r.Min.X+r.Max.X)/2 - 5
	a.SetRect(r)
	a.SetMirror(pix.MX | pix.MY)
	w.Pos = a.Bounds().Min
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()

	r = disp.Bounds()
	r.Min.Y = (r.Min.Y+r.Max.Y)/2 + 4
	a.SetRect(r)
	a.SetMirror(pix.MV)
	w.Pos = a.Bounds().Min
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()
	r = disp.Bounds()
	r.Max.Y = (r.Min.Y+r.Max.Y)/2 - 5
	a.SetRect(r)
	a.SetMirror(pix.MV | pix.MY)
	w.Pos = a.Bounds().Min
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()

	r = disp.Bounds()
	r.Min.Y = (r.Min.Y+r.Max.Y)/2 + 4
	a.SetRect(r)
	a.SetMirror(pix.MV | pix.MX)
	w.Pos = a.Bounds().Min
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()
	r = disp.Bounds()
	r.Max.Y = (r.Min.Y+r.Max.Y)/2 - 5
	a.SetRect(r)
	a.SetMirror(pix.MV | pix.MY | pix.MX)
	w.Pos = a.Bounds().Min
	for _, s := range helloWorld {
		w.WriteString(s)
	}
	markNextGlyph()

	//saveDisplay(t, disp, testFile)
	checkDisplay(t, disp, testFile)
}
