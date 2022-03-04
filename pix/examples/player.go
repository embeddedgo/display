// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build ignore

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/embeddedgo/display/font/subfont/font9/dejavusans16"
	"github.com/embeddedgo/display/font/subfont/font9/dejavusans18"
	"github.com/embeddedgo/display/pix"
)

var (
	bgColor      = color.Gray{50}
	textColor    = color.Gray{255}
	buttonColor  = color.Gray{255}
	captionColor = color.Gray{50}

	titleFont = dejavusans18.NewFace(
		dejavusans18.X0020_007e, // !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNO
		dejavusans18.X0101_0201, // āĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİ
	)
	buttonFont = dejavusans16.NewFace(
		dejavusans16.X0020_007e,
	)
)

type Driver struct {
	Name  string
	Image draw.Image
	fill  image.Uniform
	err   error
}

func (d *Driver) SetDir(dir int) image.Rectangle {
	return d.Image.Bounds()
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	draw.DrawMask(d.Image, r, src, sp, mask, mp, op)
}

func (d *Driver) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	d.fill.C = color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (d *Driver) Fill(r image.Rectangle) {
	d.Draw(r, &d.fill, image.Point{}, nil, image.Point{}, draw.Over)
}

func (d *Driver) Flush() {
	if d.err != nil {
		return
	}
	f, err := os.Create(d.Name)
	if err != nil {
		d.err = err
		return
	}
	d.err = jpeg.Encode(f, d.Image, &jpeg.Options{95})
	f.Close()
}

func (d *Driver) Err(clear bool) error {
	err := d.err
	if clear {
		d.err = nil
	}
	return err
}

func NewDisplay(name string, width, height int) *pix.Display {
	driver := &Driver{
		Name:  name,
		Image: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	return pix.NewDisplay(driver)
}

func playerView1(disp *pix.Display, artist, title string, cover image.Image) {
	a := disp.NewArea(disp.Bounds())
	a.Draw(a.Bounds(), cover, cover.Bounds().Min, nil, image.Point{}, draw.Src)
	disp.Flush()
}

func playerView2(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	a.SetRect(r.Inset(20))
	a.SetColorRGBA(0, 0, 200, 255)
	a.Fill(a.Bounds())
	disp.Flush()
}

func playerView3(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	const margin = 20
	a.SetRect(r.Inset(margin))
	sr := cover.Bounds()
	r = a.Bounds()
	r.Min = r.Size().Sub(sr.Size()).Div(2)
	a.Draw(r, cover, sr.Min, nil, image.Point{}, draw.Src)
	disp.Flush()
}

func playerView5(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	const margin = 20
	a.SetRect(r.Inset(margin))
	sr := cover.Bounds()
	r = a.Bounds()
	r.Min = r.Size().Sub(sr.Size()).Div(2)
	a.Draw(r, cover, sr.Min, nil, image.Point{}, draw.Src)

	r = disp.Bounds()
	r.Min.X += r.Dy()
	a.SetRect(r.Inset(margin))
	w := a.NewTextWriter(titleFont)
	w.Break = pix.BreakSpace
	w.SetColor(textColor)
	w.WriteString(artist + " -- " + title)

	disp.Flush()
}

func playerView7(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	const margin = 20
	a.SetRect(r.Inset(margin))
	sr := cover.Bounds()
	r = a.Bounds()
	r.Min = r.Size().Sub(sr.Size()).Div(2)
	a.Draw(r, cover, sr.Min, nil, image.Point{}, draw.Src)

	height, _ := titleFont.Size()
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y += margin
	r.Max.Y = r.Min.Y + 2*height
	a.SetRect(r)
	a.SetColorRGBA(0, 0, 200, 255)
	a.Fill(a.Bounds())

	w := a.NewTextWriter(titleFont)
	w.Break = pix.BreakSpace
	w.SetColor(textColor)
	w.WriteString(artist + " -- " + title)

	disp.Flush()
}

func hms(sec int) (h, m, s int) {
	return sec / 3600, sec % 3600 / 60, sec % 60
}

func drawTimeDuration(a *pix.Area, t, duration int) {
	h, m, s := hms(t)
	str := fmt.Sprintf("%02d:%02d:%02d /", h, m, s)
	width := pix.StringWidth(str, titleFont)
	height, _ := titleFont.Size()
	a.SetColor(bgColor)
	r := a.Bounds()
	a.Fill(r)
	w := a.NewTextWriter(titleFont)
	w.SetColor(textColor)
	w.Pos.X = r.Min.X + r.Dx()/2 - width
	w.Pos.Y = r.Min.Y + (r.Dy()-height)/2
	w.WriteString(str)
	h, m, s = hms(duration)
	fmt.Fprintf(w, " %02d:%02d:%02d", h, m, s)
	a.Flush()
}

func button1(w *pix.TextWriter, r image.Rectangle, caption string) {
	w.Area.Fill(r)
	width := pix.StringWidth(caption, buttonFont)
	height, _ := buttonFont.Size()
	w.Pos.X = r.Min.X + (r.Dx()-width)/2
	w.Pos.Y = r.Min.Y + (r.Dy()-height)/2
	w.WriteString(caption)
}

func button(w *pix.TextWriter, r image.Rectangle, caption string) {
	r = r.Inset(7)
	r.Max.X--
	r.Max.Y--
	w.Area.RoundRect(r.Min, r.Max, 7, 7, true)
	width := pix.StringWidth(caption, buttonFont)
	height, _ := buttonFont.Size()
	w.Pos.X = r.Min.X + (r.Dx()-width)/2
	w.Pos.Y = r.Min.Y + (r.Dy()-height)/2
	w.WriteString(caption)
}

func drawControls(a *pix.Area) {
	a.SetColor(buttonColor)
	w := a.NewTextWriter(buttonFont)
	w.SetColor(captionColor)
	r := a.Bounds()
	width := r.Dx() / 3
	r.Max.X = r.Min.X + width - 2
	button(w, r, "Rewind")
	r.Min.X = r.Max.X + 2
	r.Max.X = r.Min.X + width - 2
	button(w, r, "Play/Pause")
	r.Min.X = r.Max.X + 2
	r.Max.X = r.Min.X + width
	button(w, r, "FastForward")
	a.Flush()
}

func playerView8(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	const margin = 20
	a.SetRect(r.Inset(margin))
	sr := cover.Bounds()
	r = a.Bounds()
	r.Min = r.Size().Sub(sr.Size()).Div(2)
	a.Draw(r, cover, sr.Min, nil, image.Point{}, draw.Src)
	height, _ := titleFont.Size()
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y += margin
	r.Max.Y = r.Min.Y + 2*height
	a.SetRect(r)
	w := a.NewTextWriter(titleFont)
	w.Break = pix.BreakSpace
	w.SetColor(textColor)
	w.WriteString(artist + " -- " + title)

	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y = (r.Min.Y+r.Max.Y)/2 - 10
	r.Max.Y = r.Min.Y + 20
	drawTimeDuration(disp.NewArea(r), 20, 3*60+17)

	disp.Flush()
}

func playerView9(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	const margin = 20
	a.SetRect(r.Inset(margin))
	sr := cover.Bounds()
	r = a.Bounds()
	r.Min = r.Size().Sub(sr.Size()).Div(2)
	a.Draw(r, cover, sr.Min, nil, image.Point{}, draw.Src)
	height, _ := titleFont.Size()
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y += margin
	r.Max.Y = r.Min.Y + 2*height
	a.SetRect(r)
	w := a.NewTextWriter(titleFont)
	w.Break = pix.BreakSpace
	w.SetColor(textColor)
	w.WriteString(artist + " -- " + title)
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y = (r.Min.Y+r.Max.Y)/2 - 20
	r.Max.Y = r.Min.Y + 40
	drawTimeDuration(disp.NewArea(r), 20, 3*60+17)

	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Max.Y -= margin
	r.Min.Y = r.Max.Y - 30
	drawControls(disp.NewArea(r))

	disp.Flush()
}

func playerView10(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	const margin = 20
	a.SetRect(r.Inset(margin))
	sr := cover.Bounds()
	r = a.Bounds()
	r.Min = r.Size().Sub(sr.Size()).Div(2)
	a.Draw(r, cover, sr.Min, nil, image.Point{}, draw.Src)
	height, _ := titleFont.Size()
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y += margin
	r.Max.Y = r.Min.Y + 2*height
	a.SetRect(r)
	w := a.NewTextWriter(titleFont)
	w.Break = pix.BreakSpace
	w.SetColor(textColor)
	w.WriteString(artist + " -- " + title)
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y = (r.Min.Y+r.Max.Y)/2 - 20
	r.Max.Y = r.Min.Y + 40
	drawTimeDuration(disp.NewArea(r), 20, 3*60+17)

	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Max.Y -= margin
	r.Min.Y = r.Max.Y - 31
	drawControls(disp.NewArea(r))

	disp.Flush()
}

func gbutton(a *pix.Area, r image.Rectangle, caption [][4]image.Point) {
	r = r.Inset(7)
	r.Max.X--
	r.Max.Y--
	a.SetColor(buttonColor)
	a.RoundRect(r.Min, r.Max, 7, 7, true)
	p := r.Min.Add(r.Max).Div(2)
	a.SetColor(captionColor)
	for _, q := range caption {
		a.Quad(q[0].Add(p), q[1].Add(p), q[2].Add(p), q[3].Add(p), true)
	}
}

var rewind = [][4]image.Point{
	{{14, -7}, {14, -7}, {0, 0}, {14, 7}},
	{{-4, -7}, {-4, -7}, {-18, 0}, {-4, 7}},
}
var playPause = [][4]image.Point{
	{{-15, -7}, {-15, -7}, {-1, 0}, {-15, 7}},
	{{6, -7}, {9, -7}, {9, 7}, {6, 7}},
	{{14, -7}, {17, -7}, {17, 7}, {14, 7}},
}
var fastForward = [][4]image.Point{
	{{-14, -7}, {-14, -7}, {0, 0}, {-14, 7}},
	{{4, -7}, {4, -7}, {18, 0}, {4, 7}},
}

func drawControls1(a *pix.Area) {
	r := a.Bounds()
	width := r.Dx() / 3
	r.Max.X = r.Min.X + width - 3
	gbutton(a, r, rewind)
	r.Min.X = r.Max.X + 3
	r.Max.X = r.Min.X + width - 3
	gbutton(a, r, playPause)
	r.Min.X = r.Max.X + 3
	r.Max.X = r.Min.X + width
	gbutton(a, r, fastForward)
	a.Flush()
}

func playerView11(disp *pix.Display, artist, title string, cover image.Image) {
	r := disp.Bounds()
	a := disp.NewArea(r)
	a.SetColor(bgColor)
	a.Fill(a.Bounds())
	r.Max.X = r.Min.X + r.Dy()
	const margin = 20
	a.SetRect(r.Inset(margin))
	sr := cover.Bounds()
	r = a.Bounds()
	r.Min = r.Size().Sub(sr.Size()).Div(2)
	a.Draw(r, cover, sr.Min, nil, image.Point{}, draw.Src)
	height, _ := titleFont.Size()
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y += margin
	r.Max.Y = r.Min.Y + 2*height
	a.SetRect(r)
	w := a.NewTextWriter(titleFont)
	w.Break = pix.BreakSpace
	w.SetColor(textColor)
	w.WriteString(artist + " -- " + title)
	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Min.Y = (r.Min.Y+r.Max.Y)/2 - 20
	r.Max.Y = r.Min.Y + 40
	drawTimeDuration(disp.NewArea(r), 20, 3*60+17)

	r = disp.Bounds()
	r.Min.X += r.Dy()
	r.Max.X -= margin
	r.Max.Y -= margin
	r.Min.Y = r.Max.Y - 30
	drawControls1(disp.NewArea(r))

	disp.Flush()
}

func main() {
	disp := NewDisplay("/tmp/player.jpg", 640, 320)

	artist := "Gophers"
	title := "Work Hard Play Hard"
	cover := loadImage("../../testdata/gopherbug.jpg")
	//cover = images.Magnify(cover, 2, 2, images.Bilinear)
	playerView11(disp, artist, title, cover)

	checkErr(disp.Err(true))
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadImage(name string) image.Image {
	f, err := os.Open(name)
	checkErr(err)
	defer f.Close()
	img, _, err := image.Decode(f)
	checkErr(err)
	return img
}
