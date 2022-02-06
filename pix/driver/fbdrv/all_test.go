// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv_test

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/embeddedgo/display/images"
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/fbdrv"
)

var workDir = filepath.Join(os.TempDir(), "pix_test")

var (
	white = color.Gray{255}
	black = color.Gray{0}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
)

func failErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

type monofb struct {
	img  *images.Mono
	path string
}

func (fb *monofb) SetDir(dir int) (pix []byte, width, height, stride int, shift, mvxy uint8) {
	pix = fb.img.Pix
	width = fb.img.Rect.Dx()
	height = fb.img.Rect.Dy()
	stride = fb.img.Stride
	shift = uint8(fb.img.Shift)
	switch dir {
	case 1:
		mvxy = fbdrv.MV | fbdrv.MX
	case 2:
		mvxy = fbdrv.MX | fbdrv.MY
	case 3:
		mvxy = fbdrv.MV | fbdrv.MY
	}
	return
}

func (fb *monofb) Flush() ([]byte, error) {
	os.Mkdir(workDir, 0755)
	f, err := os.OpenFile(fb.path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return fb.img.Pix, err
	}
	defer f.Close()
	return fb.img.Pix, png.Encode(f, fb.img)
}

func TestMonoGraph(t *testing.T) {
	dir := 0 // change this to test different directions
	fb := &monofb{
		img:  images.NewMono(image.Rect(0, 0, 41, 81)),
		path: filepath.Join(workDir, "mono.png"),
	}
	disp := pix.NewDisplay(fbdrv.NewMono(fb))
	disp.SetDir(dir)
	a := disp.NewArea(disp.Bounds().Inset(1))
	a.SetColor(white)
	r := a.Bounds()
	a.Fill(r)
	a.SetColor(black)
	p0 := image.Pt(5, 10)
	p1 := image.Pt(25, 5)
	p2 := image.Pt(35, 10)
	p3 := image.Pt(25, 15)
	a.Quad(p0, p1, p2, p3, true)
	p := r.Min.Add(r.Max).Mul(3).Div(4)
	a.RoundRect(p, p, 5, 8, false)
	p.X = p.X / 3
	a.RoundRect(p, p, 4, 6, true)
	a.Flush()
	failErr(t, a.Err(false))
}

var img = &image.Alpha{
	Pix: []uint8{
		0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00,
		0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00,
		0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
		0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00,
		0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff,
		0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00,
		0xff, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	},
	Stride: 9,
	Rect:   image.Rectangle{Max: image.Point{9, 12}},
}

func TestMonoDraw(t *testing.T) {
	dir := 3 // change this to test different directions
	fb := &monofb{
		img:  images.NewMono(image.Rect(0, 0, 21, 21)),
		path: filepath.Join(workDir, "monodraw.png"),
	}
	disp := pix.NewDisplay(fbdrv.NewMono(fb))
	disp.SetDir(dir)
	a := disp.NewArea(disp.Bounds())
	r := image.Rectangle{Max: img.Bounds().Size()}
	p := image.Pt(8, 1)
	a.Draw(r.Add(p), img, image.Point{}, nil, image.Point{}, draw.Src)
	p = image.Pt(1, 8)
	u := &image.Uniform{color.Gray{255}}
	a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Src)
	p = image.Pt(11, 7)
	a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Over)
	a.Flush()
	failErr(t, a.Err(false))
}

type rgbfb struct {
	img  *images.RGB
	path string
}

func (fb *rgbfb) SetDir(dir int) (pix []byte, width, height, stride int, shift, mvxy uint8) {
	pix = fb.img.Pix
	width = fb.img.Rect.Dx()
	height = fb.img.Rect.Dy()
	stride = fb.img.Stride
	switch dir {
	case 1:
		mvxy = fbdrv.MV | fbdrv.MX
	case 2:
		mvxy = fbdrv.MX | fbdrv.MY
	case 3:
		mvxy = fbdrv.MV | fbdrv.MY
	}
	return
}

func (fb *rgbfb) Flush() ([]byte, error) {
	os.Mkdir(workDir, 0755)
	f, err := os.OpenFile(fb.path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return fb.img.Pix, err
	}
	defer f.Close()
	return fb.img.Pix, png.Encode(f, fb.img)
}

func TestRGBGraph(t *testing.T) {
	dir := 3 // change this to test different directions
	fb := &rgbfb{
		img:  images.NewRGB(image.Rect(0, 0, 41, 81)),
		path: filepath.Join(workDir, "rgb.png"),
	}
	disp := pix.NewDisplay(fbdrv.NewRGB(fb))
	disp.SetDir(dir)
	a := disp.NewArea(disp.Bounds().Inset(1))
	a.SetColor(white)
	r := a.Bounds()
	a.Fill(r)
	p0 := image.Pt(5, 10)
	p1 := image.Pt(25, 5)
	p2 := image.Pt(35, 10)
	p3 := image.Pt(25, 15)
	a.SetColor(red)
	a.Quad(p0, p1, p2, p3, true)
	p := r.Min.Add(r.Max).Mul(3).Div(4)
	a.SetColor(green)
	a.RoundRect(p, p, 5, 8, false)
	p.X = p.X / 3
	a.SetColor(blue)
	a.RoundRect(p, p, 4, 6, true)
	a.Flush()
	failErr(t, a.Err(false))
}

func TestRGBDraw(t *testing.T) {
	dir := 4 // change this to test different directions
	fb := &rgbfb{
		img:  images.NewRGB(image.Rect(0, 0, 21, 21)),
		path: filepath.Join(workDir, "rgbdraw.png"),
	}
	disp := pix.NewDisplay(fbdrv.NewRGB(fb))
	disp.SetDir(dir)
	a := disp.NewArea(disp.Bounds())
	r := image.Rectangle{Max: img.Bounds().Size()}
	p := image.Pt(8, 1)
	a.Draw(r.Add(p), img, image.Point{}, nil, image.Point{}, draw.Src)
	p = image.Pt(1, 8)
	u := &image.Uniform{color.Gray{255}}
	a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Src)
	p = image.Pt(11, 7)
	a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Over)
	a.Flush()
	failErr(t, a.Err(false))
}
