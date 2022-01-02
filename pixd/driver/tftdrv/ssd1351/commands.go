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
	RMDCLM  = 0xA0         // Set Re-map & Dual COM Line Mode
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

const ms = 255

// GFX contains initialization commands taken from Adafruit GFX library.
var GFX = []byte{
	300, ms, // wait 300 ms after reset
	CMDLCK, 1, 0x12,
	CMDLCK, 1, 0xB1,
	SLPIN, 0,
	CLKDIV, 1, 0xF1,
	MUXR, 1, 127,
	DOFFSET, 1, 0,
	GPIO, 1, 0,
	FUNCSEL, 1, 1,
	1,
	PHALEN, 1, 0x32,
	VCOMH, 1, 0x05,
	DISNOR, 0,
	CCABC, 3, 0xC8, 0x80, 0xC8,
	MCCC, 1, 0x0F,
	VSL, 3, 0xA0, 0xB5, 0x55,
	PHA3LEN, 1, 0x01,
	SLPOUT, 0,
	RMDCLM, 1, 0x20,
}

/*
   SSD1351_CMD_COMMANDLOCK,
   1, // Set command lock, 1 arg
   0x12,
   SSD1351_CMD_COMMANDLOCK,
   1, // Set command lock, 1 arg
   0xB1,
   SSD1351_CMD_DISPLAYOFF,
   0, // Display off, no args
   SSD1351_CMD_CLOCKDIV,
   1,
   0xF1, // 7:4 = Oscillator Freq, 3:0 = CLK Div Ratio (A[3:0]+1 = 1..16)
   SSD1351_CMD_MUXRATIO,
   1,
   127,
   SSD1351_CMD_DISPLAYOFFSET,
   1,
   0x0,
   SSD1351_CMD_SETGPIO,
   1,
   0x00,
   SSD1351_CMD_FUNCTIONSELECT,
   1,
   0x01, // internal (diode drop)
   SSD1351_CMD_PRECHARGE,
   1,
   0x32,
   SSD1351_CMD_VCOMH,
   1,
   0x05,
   SSD1351_CMD_NORMALDISPLAY,
   0,
   SSD1351_CMD_CONTRASTABC,
   3,
   0xC8,
   0x80,
   0xC8,
   SSD1351_CMD_CONTRASTMASTER,
   1,
   0x0F,
   SSD1351_CMD_SETVSL,
   3,
   0xA0,
   0xB5,
   0x55,
   SSD1351_CMD_PRECHARGE2,
   1,
   0x01,
   SSD1351_CMD_DISPLAYON,
   0,  // Main screen turn on
   0}; // END OF COMMAND LIST

#define SSD1351_CMD_SETCOLUMN 0x15      ///< See datasheet
#define SSD1351_CMD_SETROW 0x75         ///< See datasheet
#define SSD1351_CMD_WRITERAM 0x5C       ///< See datasheet
#define SSD1351_CMD_READRAM 0x5D        ///< Not currently used
#define SSD1351_CMD_SETREMAP 0xA0       ///< See datasheet
#define SSD1351_CMD_STARTLINE 0xA1      ///< See datasheet
#define SSD1351_CMD_DISPLAYOFFSET 0xA2  ///< See datasheet
#define SSD1351_CMD_DISPLAYALLOFF 0xA4  ///< Not currently used
#define SSD1351_CMD_DISPLAYALLON 0xA5   ///< Not currently used
#define SSD1351_CMD_NORMALDISPLAY 0xA6  ///< See datasheet
#define SSD1351_CMD_INVERTDISPLAY 0xA7  ///< See datasheet
#define SSD1351_CMD_FUNCTIONSELECT 0xAB ///< See datasheet
#define SSD1351_CMD_DISPLAYOFF 0xAE     ///< See datasheet
#define SSD1351_CMD_DISPLAYON 0xAF      ///< See datasheet
#define SSD1351_CMD_PRECHARGE 0xB1      ///< See datasheet
#define SSD1351_CMD_DISPLAYENHANCE 0xB2 ///< Not currently used
#define SSD1351_CMD_CLOCKDIV 0xB3       ///< See datasheet
#define SSD1351_CMD_SETVSL 0xB4         ///< See datasheet
#define SSD1351_CMD_SETGPIO 0xB5        ///< See datasheet
#define SSD1351_CMD_PRECHARGE2 0xB6     ///< See datasheet
#define SSD1351_CMD_SETGRAY 0xB8        ///< Not currently used
#define SSD1351_CMD_USELUT 0xB9         ///< Not currently used
#define SSD1351_CMD_PRECHARGELEVEL 0xBB ///< Not currently used
#define SSD1351_CMD_VCOMH 0xBE          ///< See datasheet
#define SSD1351_CMD_CONTRASTABC 0xC1    ///< See datasheet
#define SSD1351_CMD_CONTRASTMASTER 0xC7 ///< See datasheet
#define SSD1351_CMD_MUXRATIO 0xCA       ///< See datasheet
#define SSD1351_CMD_COMMANDLOCK 0xFD    ///< See datasheet
#define SSD1351_CMD_HORIZSCROLL 0x96    ///< Not currently used
#define SSD1351_CMD_STOPSCROLL 0x9E     ///< Not currently used
#define SSD1351_CMD_STARTSCROLL 0x9F    ///< Not currently used

*/
