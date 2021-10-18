// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"github.com/embeddedgo/display/pixd/driver/tftdrv"
)

const (
	transparent = 0

	osize    = 4
	otype    = 6 // Fill relies on the type field takes up two MSbits
	fastByte = 0
	fastWord = 1
	bufInit  = 2 // getBuf relies on the one bit difference to the bufFull
	bufFull  = 3 // Fill relies on the both bits set
)

// InitGFX contains initialization commands taken from Adafruit GFX library.
var InitGFX = []byte{
	SWRESET, 0,
	5, 255, // wait 5 ms
	SFD, 3, 0x03, 0x80, 0x02,
	PWCTRB, 3, 0x00, 0xC1, 0x30,
	PONSEQ, 4, 0x64, 0x03, 0x12, 0x81,
	DRVTIM, 3, 0x85, 0x00, 0x78,
	PWCTRA, 5, 0x39, 0x2C, 0x00, 0x34, 0x02,
	PUMPRT, 1, 0x20,
	DRVTIMB, 2, 0x00, 0x00,
	PWCTR1, 1, 0x23,
	PWCTR2, 1, 0x10,
	VMCTR1, 2, 0x3e, 0x28,
	VMCTR2, 1, 0x86,
	VSCRSADD, 1, 0x00,
	FRMCTR1, 2, 0x00, 0x18,
	DFUNCTR, 3, 0x08, 0x82, 0x27,
	GAMMASET, 1, 0x01,
	PGAMCTRL, 15, 0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00,
	NGAMCTRL, 15, 0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F,
	120, 255, // wait 120 ms
	SLPOUT, 0,
	5, 255,
	DISPON, 0,
	MADCTL, 1, BGR | MX, // default display orientation, must be the last one
}

func pixSet(dci tftdrv.DCI, oldpf *[1]byte, newpf byte) {
	if oldpf[0] != newpf {
		oldpf[0] = newpf
		dci.Cmd(PIXSET)
		dci.WriteBytes(oldpf[:])
	}
}
