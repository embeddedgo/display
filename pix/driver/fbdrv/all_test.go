// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fbdrv_test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/embeddedgo/display/images"
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/fbdrv"
)

var dir = filepath.Join(os.TempDir(), "pix_test")

var (
	white = color.Gray{255}
	black = color.Gray{0}
)

func failErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

type frameBuffer struct {
	img  *images.Mono
	path string
}

func (fb *frameBuffer) SetDir(dir int) (pix []byte, width, height, stride int, shift, mvxy uint8) {
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

func (fb *frameBuffer) Flush() ([]byte, error) {
	os.Mkdir(dir, 0755)
	f, err := os.OpenFile(fb.path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return fb.img.Pix, err
	}
	defer f.Close()
	return fb.img.Pix, png.Encode(f, fb.img)
}

func TestMono(t *testing.T) {
	fb := &frameBuffer{
		img:  images.NewMono(image.Rect(0, 0, 41, 81)),
		path: filepath.Join(dir, "mono.png"),
	}
	disp := pix.NewDisplay(fbdrv.NewMono(fb))
	for i := 0; i < 5; i++ {
		println(i)
		disp.SetDir(i)
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
		time.Sleep(5 * time.Second)
	}
}
