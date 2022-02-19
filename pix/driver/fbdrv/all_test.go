// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/embeddedgo/display/images"
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/fbdrv"
)

const testDir = "testdata"

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

func saveImage(t *testing.T, img image.Image, name string) {
	f, err := os.Create(filepath.Join(testDir, name))
	failErr(t, err)
	failErr(t, png.Encode(f, img))
	failErr(t, f.Close())
}

type monofb struct {
	img *images.Mono
}

func newMonoFB(width, height int) *monofb {
	return &monofb{images.NewMono(image.Rect(0, 0, width, height))}
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

func (fb *monofb) Flush() []byte        { return fb.img.Pix }
func (fb *monofb) Err(clear bool) error { return nil }

func TestMonoGraph(t *testing.T) {
	testFile := "mono_graph_%d.png"

	fb := newMonoFB(41, 81)
	disp := pix.NewDisplay(fbdrv.NewMono(fb))
	a := disp.NewArea(disp.Bounds().Inset(1))
	dir := 0
	for {
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

		//saveImage(t, fb.img, fmt.Sprintf(testFile, dir))
		checkImage(t, fb.img, fmt.Sprintf(testFile, dir))
		if dir++; dir > 4 {
			break
		}
		disp.SetDir(dir)
		a.SetRect(disp.Bounds().Inset(1))
	}
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

func TestMonoImage(t *testing.T) {
	testFile := "mono_image_%d.png"

	fb := newMonoFB(21, 21)
	disp := pix.NewDisplay(fbdrv.NewMono(fb))
	a := disp.NewArea(disp.Bounds())
	dir := 0
	for {
		a.SetColor(black)
		a.Fill(a.Bounds())
		r := image.Rectangle{Max: img.Bounds().Size()}
		p := image.Pt(8, 1)
		a.Draw(r.Add(p), img, image.Point{}, nil, image.Point{}, draw.Src)
		p = image.Pt(1, 8)
		u := &image.Uniform{color.Gray{255}}
		a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Src)
		p = image.Pt(11, 7)
		a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Over)

		//saveImage(t, fb.img, fmt.Sprintf(testFile, dir))
		checkImage(t, fb.img, fmt.Sprintf(testFile, dir))
		if dir++; dir > 4 {
			break
		}
		disp.SetDir(dir)
		a.SetRect(disp.Bounds())
	}
}

type rgbfb struct {
	img *images.RGB
}

func newRGBFB(width, height int) *rgbfb {
	return &rgbfb{images.NewRGB(image.Rect(0, 0, width, height))}
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

func (fb *rgbfb) Flush() []byte        { return fb.img.Pix }
func (fb *rgbfb) Err(clear bool) error { return nil }

func TestRGBGraph(t *testing.T) {
	testFile := "rgb_graph_%d.png"

	fb := newRGBFB(41, 81)
	disp := pix.NewDisplay(fbdrv.NewRGB(fb))
	a := disp.NewArea(disp.Bounds().Inset(1))
	dir := 0
	for {
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

		//saveImage(t, fb.img, fmt.Sprintf(testFile, dir))
		checkImage(t, fb.img, fmt.Sprintf(testFile, dir))
		if dir++; dir > 4 {
			break
		}
		disp.SetDir(dir)
		a.SetRect(disp.Bounds().Inset(1))
	}
}

func TestRGBImage(t *testing.T) {
	testFile := "rgb_image_%d.png"

	fb := newRGBFB(64, 128)
	disp := pix.NewDisplay(fbdrv.NewRGB(fb))
	a := disp.NewArea(disp.Bounds())
	dir := 0
	for {
		a.SetColor(blue)
		a.Fill(a.Bounds())
		r := image.Rectangle{Max: img.Bounds().Size()}
		p := image.Pt(8, 1)
		a.Draw(r.Add(p), img, image.Point{}, nil, image.Point{}, draw.Src)
		p = image.Pt(1, 8)
		u := &image.Uniform{color.Gray{255}}
		a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Src)
		p = image.Pt(11, 7)
		a.Draw(r.Add(p), u, image.Point{}, img, image.Point{}, draw.Over)
		f, err := os.Open("../../../testdata/gopher.png")
		failErr(t, err)
		gopher, _, err := image.Decode(f)
		f.Close()
		failErr(t, err)
		r = gopher.Bounds()
		p = image.Pt(0, 45)
		a.Draw(r.Add(p.Sub(r.Min)), gopher, r.Min, nil, image.Point{}, draw.Src)
		p = image.Pt(20, 2)
		a.Draw(r.Add(p), gopher, image.Point{}, nil, image.Point{}, draw.Over)

		//saveImage(t, fb.img, fmt.Sprintf(testFile, dir))
		checkImage(t, fb.img, fmt.Sprintf(testFile, dir))
		if dir++; dir > 4 {
			break
		}
		disp.SetDir(dir)
		a.SetRect(disp.Bounds())
	}
}

func TestRGBImageMask(t *testing.T) {
	testFile := "rgb_imagemask_%d.png"

	f, err := os.Open("../../../testdata/gopherbug.jpg")
	failErr(t, err)
	gobug, _, err := image.Decode(f)
	f.Close()
	failErr(t, err)

	gobugMask := new(images.AlphaN)
	gobugMask.Rect.Min = image.Pt(31, 21)
	gobugMask.Rect.Max = gobugMask.Rect.Min.Add(image.Pt(211, 251))
	gobugMask.Stride = 27
	gobugMask.Pix, err = os.ReadFile("../../../testdata/gopherbug.211x251s27b1")

	fb := newRGBFB(272*2, 272*2)
	disp := pix.NewDisplay(fbdrv.NewRGB(fb))
	a := disp.NewArea(disp.Bounds())
	dir := 0
	for {
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
		for i, b := range gobugMask.Pix {
			gobugMask.Pix[i] = ^b
		}

		//saveImage(t, fb.img, fmt.Sprintf(testFile, dir))
		checkImage(t, fb.img, fmt.Sprintf(testFile, dir))
		if dir++; dir > 4 {
			break
		}
		disp.SetDir(dir)
		a.SetRect(disp.Bounds())
	}
}
