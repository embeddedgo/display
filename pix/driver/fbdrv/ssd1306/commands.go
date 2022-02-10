// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1306

const (
	LCSA   = 0x00 // Lower column start address for page addressing mode
	HCSA   = 0x10 // Higher column start address for page addressing mode
	MAM    = 0x20 // Memory addressing mode
	CA     = 0x21 // Column start and end address
	PA     = 0x22 // Page start and end address
	CSR    = 0x26 // Continuous right horizontal scroll
	CSL    = 0x27 // Continuous left horizontal scroll
	CSVR   = 0x29 // Continuous vertical and right horizontal scroll
	CSVL   = 0x2A // Continuous vertical and left horizontal scroll
	CSD    = 0x2E // Deactivate scroll
	CSA    = 0x2F // Activate scroll
	DSL    = 0x40 // Display start line
	CONTR  = 0x81 // Set contrast control
	CPUMP  = 0x8D // Charge pump settings
	SRM    = 0xA0 // Segment re-map
	VSA    = 0xA3 // Vertical Scroll Area
	DOUT   = 0xA4 // Output follows RAM content
	DWHT   = 0xA5 // Set entire display white
	DNOR   = 0xA6 // Set normal display
	DINV   = 0xA7 // Set inverse display
	MR     = 0xA8 // Multiplex ratio
	SLPIN  = 0xAE // Set sleep mode ON
	SLPOUT = 0xAF // Set sleep mode OFF
	PSA    = 0xB0 // Page start address (0..7) for page addressing mode
	COMSD  = 0xC0 // COM output scan direction
	VOFF   = 0xD3 // Vertical display offset (0..63)
	CLKDIV = 0xD5 // Display clock divide ratio / Oscillator frequency
	PRE    = 0xD9 // Pre-charge period
	COMPHC = 0xDA // COM pins hardware configuration
	VCOMH  = 0xDB // VCOMH deselect level
	NOP    = 0xE3 // No operation
)

const ms = 255

var GFX128x64 = []byte{
	1, ms,
	SLPIN, 0,
	CLKDIV, 1, 0x80,
	MR, 1, 63,
	VOFF, 1, 0,
	DSL, 0,
	CPUMP, 1, 0x14, // enable charge pump during display on
	MAM, 1, 0x01, // vertical addressing mode
	SRM | 0x1, 0,
	COMSD | 0x8, 0,
	COMPHC, 1, 0x12, // alternative configuration (8 row display)
	CONTR, 1, 0xCF,
	PRE, 1, 0xF1,
	VCOMH, 1, 0x40,
	DOUT, 0,
	DNOR, 0,
	CSD, 0,
	SLPOUT, 0,
}

var GFX128x64ExtVcc = []byte{
	1, ms,
	SLPIN, 0,
	CLKDIV, 1, 0x80,
	MR, 1, 63,
	VOFF, 1, 0,
	DSL, 0,
	CPUMP, 1, 0x10,
	MAM, 1, 0x01, // vertical addressing mode
	SRM | 0x1, 0,
	COMSD | 0x8, 0,
	COMPHC, 1, 0x12, // alternative configuration (8 row display)
	CONTR, 1, 0x9F,
	PRE, 1, 0x22,
	VCOMH, 1, 0x40,
	DOUT, 0,
	DNOR, 0,
	CSD, 0,
	SLPOUT, 0,
}
