// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package displays

import (
	"image"

	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/fbdrv"
	"github.com/embeddedgo/display/pix/driver/fbdrv/sh1106"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
)

func newSH1106GFX128x64(dci tftdrv.DCI) *pix.Display {
	fb := sh1106.New(dci)
	fb.Init(sh1106.GFX128x64)
	disp := pix.NewDisplay(fbdrv.NewMono(fb))
	disp.SetRect(image.Rect(2, 0, 130, 64)) // adjust to the view port
	return disp
}

// Adafruit Monochrome 1.3" 128x64 OLED Graphic Display.
func Adafruit_1i3_128x64_OLED_SH1106() Param {
	return Param{
		0,
		sh1106.MaxSPIWriteClock,
		newSH1106GFX128x64,
	}
}
