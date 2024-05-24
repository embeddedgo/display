// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package displays

import (
	"image"

	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
	"github.com/embeddedgo/display/pix/driver/tftdrv/st7789"
)

var init_240x240_IPS_ST7789 = [...]byte{st7789.VSCSAD, 0, 40}

func new240x240_IPS_ST7789(dci tftdrv.DCI) *pix.Display {
	drv := st7789.New(dci)
	drv.Init(st7789.GFX)

	// Move the 240 x 240 view port in the middle of the 240 x 320 frame
	// memory to simplify handling of display rotations.
	dci.Cmd(init_240x240_IPS_ST7789[:1], tftdrv.Write)
	dci.WriteBytes(init_240x240_IPS_ST7789[1:3])

	disp := pix.NewDisplay(drv)
	disp.SetRect(image.Rect(0, 40, 240, 240+40)) // adjust to the view port
	return disp
}

// Adafruit 1.54" 240x240 Wide Angle TFT LCD Display with MicroSD - ST7789
func Adafruit_1i54_240x240_IPS_ST7789() Param {
	return Param{
		0,
		st7789.MaxSPIWriteClock,
		new240x240_IPS_ST7789,
	}
}

// ER-TFTM1.54-1 IPS LCD Module - ST7789
func ERTFTM_1i54_240x240_IPS_ST7789() Param {
	return Adafruit_1i54_240x240_IPS_ST7789()
}

