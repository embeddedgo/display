// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package displays contains initialization functions for the most popular
// displays.
//
// If you do not found a function for your display here do not worry. If your
// display uses a display controller that had a driver in ../driver/tftdrv you
// can easily initialize your display this way:
//
//	drv := driver.New(dci) // or NewOver(dci) if your display supports reading
//	drv.Init(driver.InitCommands)
//	disp := pix.NewDisplay(drv)
//
// which is basically what all functions in this package do internally.
package displays

import (
	"image"

	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/fbdrv"
	"github.com/embeddedgo/display/pix/driver/fbdrv/ssd1306"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
	"github.com/embeddedgo/display/pix/driver/tftdrv/ili9341"
	"github.com/embeddedgo/display/pix/driver/tftdrv/ili9486"
	"github.com/embeddedgo/display/pix/driver/tftdrv/ssd1351"
	"github.com/embeddedgo/display/pix/driver/tftdrv/st7789"
)

type Param struct {
	MaxReadClk  int
	MaxWriteClk int
	New         func(dci tftdrv.DCI) *pix.Display
}

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

func newMSP4022(dci tftdrv.DCI) *pix.Display {
	drv := ili9486.NewOver(dci)
	drv.Init(ili9486.MSP4022)
	return pix.NewDisplay(drv)
}

// MSP4022 4.0" TFT LCD SPI Module - ILI9486
func MSP4022_4i0_320x480_TFT_ILI9486() Param {
	return Param{
		ili9486.MaxSPIReadClock,
		ili9486.MaxSPIWriteClock,
		newMSP4022,
	}
}

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

func newSSD1306GFX128x64(dci tftdrv.DCI) *pix.Display {
	fb := ssd1306.New(dci)
	fb.Init(ssd1306.GFX128x64)
	return pix.NewDisplay(fbdrv.NewMono(fb))
}

// Adafruit Monochrome 0.96" 128x64 OLED Graphic Display
func Adafruit_0i96_128x64_OLED_SSD1306() Param {
	return Param{
		0,
		ssd1306.MaxSPIWriteClock,
		newSSD1306GFX128x64,
	}
}
