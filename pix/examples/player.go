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

	"github.com/embeddedgo/display/pix"
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
	d.err = jpeg.Encode(f, d.Image, nil)
	f.Close()
}

func (d *Driver) Err(clear bool) error {
	err := d.err
	if clear {
		d.err = nil
	}
	return err
}

func newDisplay(name string, width, height int) *pix.Display {
	driver := &Driver{
		Name:  name,
		Image: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	return pix.NewDisplay(driver)
}

func playerView(disp *pix.Display, title string, cover image.Image) {
	a := disp.NewArea(disp.Bounds())
	a.Draw(a.Bounds(), cover, cover.Bounds().Min, nil, image.Point{}, draw.Src)
	a.Flush()
}

func main() {
	disp := newDisplay("/tmp/player.jpg", 640, 320)

	title := "Gophers - Work Hard Play Hard"
	cover := loadImage("../../testdata/gopherbug.jpg")
	playerView(disp, title, cover)

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
