// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package displays

import (
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
	"github.com/embeddedgo/display/pix/driver/tftdrv/ili9341"
)

func newILI9341GFX(dci tftdrv.DCI) *pix.Display {
	drv := ili9341.NewOver(dci)
	drv.Init(ili9341.GFX)
	return pix.NewDisplay(drv)
}

// Adafruit 2.8" TFT LCD with Touchscreen Breakout Board with MicroSD - ILI9341
func Adafruit_2i8_240x320_TFT_ILI9341() Param {
	return Param{
		ili9341.MaxSPIReadClock,
		ili9341.MaxSPIWriteClock,
		newILI9341GFX,
	}
}
