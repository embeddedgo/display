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

	"github.com/embeddedgo/display/font/subfont/font9/dejavusans18"
	"github.com/embeddedgo/display/math2d"
	"github.com/embeddedgo/display/pix"
)

var (
	bgColor   = color.Gray{50}
	textColor = color.Gray{255}
	ctrlColor = color.Gray{255}

	titleFont = dejavusans18.NewFace(
		dejavusans18.X0020_007e, // !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNO
		dejavusans18.X0101_0201, // āĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİ
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

func hms(sec int) (h, m, s int) {
	return sec / 3600, sec % 3600 / 60, sec % 60
}

func min(a, b int) int {
	if a > b {
		a = b
	}
	return a
}

func printTime(w *pix.TextWriter, t int) {
	h, m, s := hms(t)
	str := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	width := pix.StringWidth(str, titleFont)
	w.Pos.X -= width / 2
	w.WriteString(str)
}

func drawTimeDuration1(a *pix.Area, t, duration int) {
	r := a.Bounds()
	a.SetColor(bgColor)
	a.Fill(r)

	p := r.Min.Add(r.Max).Div(2)
	rmax := min(r.Dx(), r.Dy()) / 2
	rmin := rmax - 20
	th0 := math2d.RightAngle + math2d.RightAngle/3
	th1 := math2d.RightAngle - math2d.RightAngle/3
	a.SetColor(ctrlColor)
	a.Arc(p, rmin, rmin, rmax, rmax, th0, th1, false)
	if t > 0 {
		alpha := int64(math2d.FullAngle*3/4) * int64(t) / int64(duration)
		a.Arc(p, rmin, rmin, rmax, rmax, th0, th0+int32(alpha), true)
	}

	w := a.NewTextWriter(titleFont)
	w.SetColor(textColor)
	height, _ := titleFont.Size()

	w.Pos.X = p.X
	w.Pos.Y = p.Y/2 - height/2
	printTime(w, t)

	w.Pos.X = p.X
	w.Pos.Y += p.Y
	printTime(w, duration)

	// play button
	p0 := image.Pt(-14, -14).Add(p)
	p1 := image.Pt(14, 0).Add(p)
	p2 := image.Pt(-14, 14).Add(p)
	a.Quad(p0, p0, p1, p2, true)

	a.Flush()
}

func playerView(disp *pix.Display, artist, title string, cover image.Image) {
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
	r.Min.Y = r.Max.Y + margin/2
	r.Max.Y = disp.Bounds().Max.Y - margin/2
	drawTimeDuration1(disp.NewArea(r), 58, 3*60+17)

	disp.Flush()
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

func main() {
	disp := NewDisplay("/tmp/player.jpg", 640, 320)
	artist := "Gophers"
	title := "Work Hard Play Hard"
	cover := loadImage("../../testdata/gopherbug.jpg")
	//cover = images.Magnify(cover, 2, 2, images.Bilinear)
	playerView(disp, artist, title, cover)
	checkErr(disp.Err(true))
}
