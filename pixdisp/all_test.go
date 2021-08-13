// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp_test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/embeddedgo/display/pixdisp"
	"github.com/embeddedgo/display/pixdisp/drivers/imgdrv"
)

func TestAll(t *testing.T) {
	failErr := func(err error) {
		if err != nil {
			t.Error(err)
		}
	}

	dir := filepath.Join(os.TempDir(), "pixdisp_test")
	os.Mkdir(dir, 0755)
	failErr(os.Chdir(dir))

	fmt.Println("see images in:", dir)

	img := image.NewNRGBA(image.Rect(0, 0, 91, 400))
	disp := pixdisp.NewDisplay(imgdrv.New(img))

	a := disp.NewArea(disp.Bounds().Inset(5))
	a.SetColor(color.Gray{220})
	a.Fill(a.Bounds())

	max := a.Bounds().Max

	a.SetColor(pixdisp.RGB565(0x1234))
	for x := 0; x < max.X; x += 2 {
		a.DrawPixel(image.Pt(x, 0))
		a.DrawPixel(image.Pt(x+1, max.Y-1))
	}
	for y := 2; y < max.Y; y += 2 {
		a.DrawPixel(image.Pt(0, y))
		a.DrawPixel(image.Pt(max.X-1, y-1))
	}

	x := max.X / 2
	for r := 0; r < 19; r++ {
		y := 3 + (r+2)*r
		a.SetColorRGB(0, 100, 200)
		a.FillCircle(image.Pt(x, y), r)
		a.SetColor(color.RGBA{100, 50, 0, 255})
		a.DrawCircle(image.Pt(x, y), r)
	}

	a.SetColorRGB(250, 100, 0)
	for i := 0; i < 11; i++ {
		x := 3*i + 7
		y := i*i + 2
		a.DrawLine(image.Pt(2, y), image.Pt(x, 2))
		a.DrawLine(image.Pt(max.X-1-2, y), image.Pt(max.X-1-x, 2))
	}

	f, err := os.OpenFile("image.png", os.O_WRONLY|os.O_CREATE, 0755)
	failErr(err)
	failErr(png.Encode(f, img))
	failErr(f.Close())
}
