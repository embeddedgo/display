// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package displays

import (
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
	"github.com/embeddedgo/display/pix/driver/tftdrv/ssd1351"
)

func newSSD1351GFX(dci tftdrv.DCI) *pix.Display {
	drv := ssd1351.New(dci)
	drv.Init(ssd1351.GFX)
	return pix.NewDisplay(drv)
}

// Adafruit OLED Breakout Board - 16-bit Color 1.5" - UG-2828GDEDF11/SSD1351
func Adafruit_1i5_128x128_OLED_SSD1351() Param {
	return Param{
		0,
		ssd1351.MaxSPIWriteClock,
		newSSD1351GFX,
	}
}

func newUG2828GDEDF11(dci tftdrv.DCI) *pix.Display {
	drv := ssd1351.New(dci)
	drv.Init(ssd1351.UG2828GDEDF11)
	return pix.NewDisplay(drv)
}

// Waveshare 128x128, General 1.5inch OLED display Module - UG-2828GDEDF11/SSD1351
func Waveshare_1i5_128x128_OLED_SSD1351() Param {
	return Param{
		ssd1351.MaxSPIReadClock,
		ssd1351.MaxSPIWriteClock,
		newUG2828GDEDF11,
	}
}
