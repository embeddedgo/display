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

func saveImage(name string, img image.Image) error {
	os.Mkdir(dir, 0755)
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func TestMono(t *testing.T) {
	flush := func(img *images.Mono) error { return saveImage("mono.png", img) }
	disp := pix.NewDisplay(fbdrv.NewMono(41, 81, flush))
	a := disp.NewArea(disp.Bounds().Inset(1))
	a.SetColor(white)
	r := a.Bounds()
	a.Fill(r)
	r = r.Inset(1)
	a.SetColor(black)
	p := r.Min.Add(r.Max).Div(2)
	a.RoundRect(p, p, r.Dx()/2, r.Dy()/2, true)
	a.Flush()
	failErr(t, a.Err(false))
}
