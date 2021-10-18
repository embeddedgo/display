// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epson

import (
	"image"

	"github.com/embeddedgo/display/pixd/driver/tftdrv"
)

// S1D15G00 commands
const (
	DISON   = 0xAF // Display on
	DISOFF  = 0xAE // Display off
	DISNOR  = 0xA6 // Normal display
	DISINV  = 0xA7 // Inverse display
	COMSCN  = 0xBB // Common scan direction
	DISCTL  = 0xCA // Display control
	SLPIN   = 0x95 // Sleep in
	SLPOUT  = 0x94 // Sleep out
	PASET   = 0x75 // Page address set
	CASET   = 0x15 // Column address set
	DATCTL  = 0xBC // Data scan direction, etc.
	RGBSET8 = 0xCE // 256-color position set
	RAMWR   = 0x5C // Writing to memory
	RAMRD   = 0x5D // Reading from memory
	PTLIN   = 0xA8 // Partial display in
	PTLOUT  = 0xA9 // Partial display out
	RMWIN   = 0xE0 // Read and modify write
	RMWOUT  = 0xEE // End
	ASCSET  = 0xAA // Area scroll set
	SCSTART = 0xAB // Scroll start set
	OSCON   = 0xD1 // Internal oscillation on
	OSCOFF  = 0xD2 // Internal oscillation off
	PWRCTR  = 0x20 // Power control
	VOLCTR  = 0x81 // Electronic volume control
	VOLUP   = 0xD6 // Increment electronic control by 1
	VOLDOWN = 0xD7 // Decrement electronic control by 1
	TMPGRD  = 0x82 // Temperature gradient set
	EPCTIN  = 0xCD // Control EEPROM
	EPCOUT  = 0xCC // Cancel EEPROM control
	EPMWR   = 0xFC // Write into EEPROM
	EPMRD   = 0xFD // Read from EEPROM
	EPSRRD1 = 0x7C // Read register 1
	EPSRRD2 = 0x7D // Read register 2
	NOP     = 0x25 // NOP instruction
)

func StartWrite8(dci tftdrv.DCI, xarg *[4]byte, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	xarg[0] = uint8(r.Min.X)
	xarg[1] = uint8(r.Max.X)
	dci.WriteBytes(xarg[:2])
	dci.Cmd(PASET)
	xarg[0] = uint8(r.Min.Y)
	xarg[1] = uint8(r.Max.Y)
	dci.WriteBytes(xarg[:2])
	dci.Cmd(RAMWR)
}

func StartRead8(dci tftdrv.DCI, xarg *[4]byte, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	xarg[0] = uint8(r.Min.X)
	xarg[1] = uint8(r.Max.X)
	dci.WriteBytes(xarg[:2])
	dci.Cmd(PASET)
	xarg[0] = uint8(r.Min.Y)
	xarg[1] = uint8(r.Max.Y)
	dci.WriteBytes(xarg[:2])
	dci.Cmd(RAMRD)
}

func StartWrite16(dci tftdrv.DCI, xarg *[4]byte, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	xarg[0] = uint8(r.Min.X >> 8)
	xarg[1] = uint8(r.Min.X)
	xarg[2] = uint8(r.Max.X >> 8)
	xarg[3] = uint8(r.Max.X)
	dci.WriteBytes(xarg[:])
	dci.Cmd(PASET)
	xarg[0] = uint8(r.Min.Y >> 8)
	xarg[1] = uint8(r.Min.Y)
	xarg[2] = uint8(r.Max.Y >> 8)
	xarg[3] = uint8(r.Max.Y)
	dci.WriteBytes(xarg[:])
	dci.Cmd(RAMWR)
}

func StartRead16(dci tftdrv.DCI, xarg *[4]byte, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	xarg[0] = uint8(r.Min.X >> 8)
	xarg[1] = uint8(r.Min.X)
	xarg[2] = uint8(r.Max.X >> 8)
	xarg[3] = uint8(r.Max.X)
	dci.WriteBytes(xarg[:])
	dci.Cmd(PASET)
	xarg[0] = uint8(r.Min.Y >> 8)
	xarg[1] = uint8(r.Min.Y)
	xarg[2] = uint8(r.Max.Y >> 8)
	xarg[3] = uint8(r.Max.Y)
	dci.WriteBytes(xarg[:])
	dci.Cmd(RAMRD)
}
