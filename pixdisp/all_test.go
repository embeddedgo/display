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
	"github.com/embeddedgo/display/pixdisp/imgdrv"
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

	img := image.NewNRGBA(image.Rect(0, 0, 90, 400))
	disp := pixdisp.NewDisplay(imgdrv.Driver{img})

	a := disp.NewArea(disp.Bounds().Inset(5))
	a.SetColor(color.Gray{220})
	a.FillRect(a.Bounds())

	a.SetColor(pixdisp.RGB565(0x1234))
	for x := 0; x < a.Bounds().Max.X; x += 2 {
		a.DrawPixel(image.Pt(x, 0))
		a.DrawPixel(image.Pt(x+1, a.Bounds().Max.Y-1))
	}
	for y := 2; y < a.Bounds().Max.Y; y += 2 {
		a.DrawPixel(image.Pt(0, y))
		a.DrawPixel(image.Pt(a.Bounds().Max.X-1, y-1))
	}

	a.SetColorRGB(255, 0, 0)
	a.FillRect(image.Rect(1, 1, 2, 2))

	a.SetColor(pixdisp.RGB565(0x3f << 5))
	a.FillRect(image.Rect(1, 2, 4, 4))

	for r := 0; r < 19; r++ {
		y := 3 + (r+2)*r
		a.SetColorRGB(0, 100, 200)
		a.FillCircle(image.Pt(20, y), r)
		a.FillCircle1(image.Pt(60, y), r)

		a.SetColorRGB(100, 50, 0)
		a.DrawCircle(image.Pt(20, y), r)
		a.DrawCircle1(image.Pt(60, y), r)
	}

	f, err := os.OpenFile("image.png", os.O_WRONLY|os.O_CREATE, 0755)
	failErr(err)
	failErr(png.Encode(f, img))
	failErr(f.Close())
}
