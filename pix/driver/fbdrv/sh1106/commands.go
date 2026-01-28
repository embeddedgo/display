// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh1106

const (
	LCSA   = 0x00 // Lower column start address for page addressing mode
	HCSA   = 0x10 // Higher column start address for page addressing mode
	CPMV   = 0x30 // Charge Pump voltage.
	DSL    = 0x40 // Display start line
	CONTR  = 0x81 // Set contrast control
	SRM    = 0xA0 // Segment re-map
	DOUT   = 0xA4 // Output follows RAM content
	DWHT   = 0xA5 // Set entire display white
	DNOR   = 0xA6 // Set normal display
	DINV   = 0xA7 // Set inverse display
	MR     = 0xA8 // Multiplex ratio
	DCDC   = 0xAD // DC-DC ON/OFF
	SLPIN  = 0xAE // Set sleep mode ON
	SLPOUT = 0xAF // Set sleep mode OFF
	PA     = 0xB0 // Page start and end address
	COMSD  = 0xC0 // COM output scan direction
	VOFF   = 0xD3 // Vertical display offset (0..63)
	CLKDIV = 0xD5 // Display clock divide ratio / Oscillator frequency
	PRE    = 0xD9 // Pre-charge period
	COMPHC = 0xDA // COM pins hardware configuration
	VCOMH  = 0xDB // VCOMH deselect level
	RMW    = 0xE0 // Read-Modify-Write start
	END    = 0xEE // Read-Modify-Write end
	NOP    = 0xE3 // No operation
)

var GFX128x64 = []byte{
	SLPIN,
	CLKDIV, 0x80,
	MR, 64 - 1,
	VOFF, 0,
	DSL | 0,
	DCDC, 0x8B,
	SRM | 0x01,
	COMSD | 0x08,
	COMPHC, 0x12, // alternative configuration (8 row display)
	CONTR, 0x80, // default contrest
	PRE, 0x1F,
	VCOMH, 0x40,
	CPMV | 3,
	LCSA | 0,
	DOUT,
	DNOR,
	SLPOUT,
}

var _GFX128x64 = []byte{
	SLPIN,
	MR, 64 - 1,
	VOFF, 0,
	DSL | 0,
	SRM | 0,
	COMSD | 0,
	COMPHC, 0x12,
	CONTR, 0x0F,
	CPMV | 0,
	DOUT,
	DNOR,
	CLKDIV, 0xF0,
	SLPOUT,
}
