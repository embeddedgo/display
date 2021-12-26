// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"image"

	"github.com/embeddedgo/display/pixd"
	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/st7789"
)

// Adafruit 1.54" 240x240 Wide Angle TFT LCD Display with MicroSD - ST7789
func Adafruit154IPS(dci tftdrv.DCI) *pixd.Display {
	drv := st7789.New(dci)
	drv.Init(st7789.GFX)

	// Move the 240 x 240 view port in the middle of the 240 x 320 frame memory
	// to simplify handling of display rotations.
	dci.Cmd(st7789.VSCSAD)
	dci.WriteBytes([]byte{0, 40})

	disp := pixd.NewDisplay(drv)
	disp.SetRect(image.Rect(0, 40, 240, 240+40)) // adjust disp to the view port

	return disp
}

// ER-TFTM1.54-1 IPS LCD Module
func ERTFTM154(dci tftdrv.DCI) *pixd.Display { return Adafruit154IPS(dci) }
