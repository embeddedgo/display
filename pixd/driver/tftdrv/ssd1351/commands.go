// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1351

import "github.com/embeddedgo/display/pixd/driver/tftdrv/internal/epson"

// SSD1351 commands
const (
	CASET   = epson.CASET  // Set Column Address
	PASET   = epson.PASET  // Set Row Address
	RAMWR   = epson.RAMWR  // Write RAM Command
	RAMRD   = epson.RAMRD  // Read RAM Command
	RMCD    = 0xA0         // Set Re-map & Color Depth
	VSCSAD  = 0xA1         // Set Display Start Line
	DOFFSET = 0xA2         // Set Display Offset
	DISBLK  = 0xA4         // Set Entire Display Black
	DISWHT  = 0xA5         // Set Entire Display White
	DISNOR  = epson.DISNOR // Set Normal Display
	DISINV  = epson.DISINV // Set Inverse Display
	FUNCSEL = 0xAB         // Set Function selection
	SLPIN   = epson.DISOFF // Set Sleep mode ON
	SLPOUT  = epson.DISON  // Set Sleep mode OFF
	PHALEN  = 0xB1         // Set Phase Length (pase 1 and 2)
	DENH    = 0xB2         // Display Enhancement
	CLKDIV  = 0xB3         // Set Front Clock Divider / Oscillator Freq
	VSL     = 0xB4         // Set Segment Low Voltage
	GPIO    = 0xB5         // Set GPIO
	PHA3LEN = 0xB6         // Set Second Pre-charge period (phase 3)
	GRAYLUT = 0xB8         // Look Up Table for Gray Scale Pulse width
	LINLUT  = 0xB9         // Use Built-in Linear LUT
	VPREC   = 0xBB         // Set Pre-charge voltage
	VCOMH   = 0xBE         // Set VCOMH Voltage
	CCABC   = 0xC1         // Set Contrast Current for Color A,B,C
	MCCC    = 0xC7         // Master Contrast Current Control
	MUXR    = 0xCA         // Set Multiplex Ratio
	CMDLCK  = 0xFD         // Set Command Lock
)

// RMCD arguments
const (
	VAI       = 1 << 0 // Vertical (instead of horizontal) address increment
	C127_SEG0 = 1 << 1 // Column address 127 (instead of 0) is mapped to SEG0
	CBA       = 1 << 2 // CBA (intead of ABC) color sequence
	COMn_COM0 = 1 << 4 // COM[n-1] to COM0 scan (instead of COM0 to COM[n-1])
	COMSplit  = 1 << 5 // Enable COM Split Odd Even
	CF2       = 1 << 6 // color format 2
	RGB18     = 1 << 7 // 18-bit 666 pixel format (instead of 16-bit 565)
)

const ms = 255

// GFX contains initialization commands taken from Adafruit GFX library.
var GFX = []byte{
	100, ms, // wait 100 ms after reset
	CMDLCK, 1, 0x12,
	CMDLCK, 1, 0xB1,
	SLPIN, 0,
	CLKDIV, 1, 0xF1,
	MUXR, 1, 127,
	DOFFSET, 1, 0,
	GPIO, 1, 0,
	FUNCSEL, 1, 0x01,
	PHALEN, 1, 0x32,
	VCOMH, 1, 0x05,
	DISNOR, 0,
	CCABC, 3, 0xC8, 0x80, 0xC8,
	MCCC, 1, 0x0F,
	VSL, 3, 0xA0, 0xB5, 0x55,
	PHA3LEN, 1, 0x01,
	SLPOUT, 0,
	RMCD, 1, COMSplit | COMn_COM0 | CBA, // default display orientation, must be the last one
}

// UG2828GDEDF11 contains the initialization commands taken from the
// documentation of UG-2828GDEDF11 1.5" 128x128 RGB OLED. This display is used
// by Waveshare and probably by Adafruit 1.5" RGB OLED modules.
var UG2828GDEDF11 = []byte{
	100, ms, // wait 100 ms after reset
	CMDLCK, 1, 0x12,
	CMDLCK, 1, 0xB1,
	SLPIN, 0,
	CLKDIV, 1, 0xF1,
	MUXR, 1, 0x7F,
	DOFFSET, 1, 0x00,
	VSCSAD, 1, 0x00,
	// RMCD, 1, 0xB4 moved to the end
	GPIO, 1, 0x00,
	FUNCSEL, 1, 0x01,
	VSL, 3, 0xA0, 0xB5, 0x55,
	CCABC, 3, 0xC8, 0x80, 0xC8,
	MCCC, 1, 0x0F,
	GRAYLUT, 63,
	0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11,
	0x12, 0x13, 0x15, 0x17, 0x19, 0x1B, 0x1D, 0x1F,
	0x21, 0x23, 0x25, 0x27, 0x2A, 0x2D, 0x30, 0x33,
	0x36, 0x39, 0x3C, 0x3F, 0x42, 0x45, 0x48, 0x4C,
	0x50, 0x54, 0x58, 0x5C, 0x60, 0x64, 0x68, 0x6C,
	0x70, 0x74, 0x78, 0x7D, 0x82, 0x87, 0x8C, 0x91,
	0x96, 0x9B, 0xA0, 0xA5, 0xAA, 0xAF, 0xB4,
	PHALEN, 1, 0x32,
	DENH, 3, 0xA4, 0x00, 0x00,
	VPREC, 1, 0x17,
	PHA3LEN, 1, 0x01,
	VCOMH, 1, 0x05,
	DISNOR, 0,
	SLPOUT, 0,
	RMCD, 1, COMSplit | COMn_COM0 | CBA, // default display orientation, must be the last one
}
