// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1351

import "github.com/embeddedgo/display/pixd/driver/tftdrv/internal/epson"

// SSD1351 commands
const (
	CASET  = epson.CASET  // Set Column Address
	PASET  = epson.PASET  // Set Row Address
	RAMWR  = epson.RAMWR  // Write RAM Command
	RAMRD  = epson.RAMRD  // Read RAM Command
	_      = 0xA0         // Set Re-map & Dual COM Line Mode
	_      = 0xA1         // Set Display Start Line
	_      = 0xA2         // Set Display Offset
	DISBLK = 0xA4         // Set Entire Display Black
	DISWHT = 0xA5         // Set Entire Display White
	DISNOR = epson.DISNOR // Set Normal Display
	DISINV = epson.DISINV // Set Inverse Display
	_      = 0xAB         // Set Function selection
	SLPIN  = epson.DISOFF // Set Sleep mode ON
	SLPOUT = epson.DISON  // Set Sleep mode OFF
	_      = 0xB1         // Set Phase Length
	_      = 0xB2         // Display Enhancement
	_      = 0xB3         // Set Front Clock Divider / Oscillator Freq
	_      = 0xB5         // Set GPIO
	_      = 0xB6         // Set Second Pre-charge period
	_      = 0xB8         // Look Up Table for Gray Scale Pulse width
	_      = 0xB9         // Use Built-in Linear LUT
	_      = 0xBB         // Set Pre-charge voltage
	VCOMH  = 0xBE         // Set VCOMH Voltage
	_      = 0xC1         // Set Contrast Current for Color A,B,C
	MCCC   = 0xC7         // Master Contrast Current Control
	_      = 0xCA         // Set Multiplex Ratio
	CMDLCK = 0xFD         // Set Command Lock
)
