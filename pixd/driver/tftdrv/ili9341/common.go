// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"image"
	"time"

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
}

// InitST contains initialization commands taken from STM32Cube library.
var InitST = []byte{
	0xCA, 3, 0xC3, 0x08, 0x50,
	PWCTRB, 3, 0x00, 0xC1, 0x30,
	PONSEQ, 4, 0x64, 0x03, 0x12, 0x81,
	DRVTIM, 3, 0x85, 0x00, 0x78,
	PWCTRA, 5, 0x39, 0x2C, 0x00, 0x34, 0x02,
	PUMPRT, 1, 0x20,
	DRVTIMB, 2, 0x00, 0x00,
	FRMCTR1, 2, 0x00, 0x1B,
	DFUNCTR, 2, 0x0A, 0xA2,
	PWCTR1, 1, 0x10,
	PWCTR2, 1, 0x10,
	VMCTR1, 2, 0x45, 0x15,
	VMCTR2, 1, 0x90,
	MADCTL, 1, 0xC8,
	EN3G, 1, 0x00,
	IFMODE, 1, 0xC2,
	DFUNCTR, 4, 0x0A, 0xA7, 0x27, 0x04,
	IFCTL, 3, 0x01, 0x00, 0x06,
	GAMMASET, 1, 0x01,
	PGAMCTRL, 15, 0x0F, 0x29, 0x24, 0x0C, 0x0E, 0x09, 0x4E, 0x78, 0x3C, 0x09, 0x13, 0x05, 0x17, 0x11, 0x00,
	NGAMCTRL, 15, 0x00, 0x16, 0x1B, 0x04, 0x11, 0x07, 0x31, 0x33, 0x42, 0x05, 0x0C, 0x0A, 0x28, 0x2F, 0x0F,
}

func initialize(dci tftdrv.DCI, cmds []byte, swreset bool) {
	if swreset {
		dci.Cmd(SWRESET)
	}
	resetTime := time.Now()
	time.Sleep(5 * time.Millisecond)
	i := 0
	for i < len(cmds) {
		dci.Cmd(cmds[i])
		n := int(cmds[i+1])
		i += 2
		if n != 0 {
			k := i + n
			data := cmds[i:k]
			i = k
			dci.WriteBytes(data)
		}
	}
	time.Sleep(resetTime.Add(120 * time.Millisecond).Sub(time.Now()))
	dci.Cmd(SLPOUT)
	time.Sleep(5 * time.Millisecond)
	dci.Cmd(DISPON)
}

func capaset(dci tftdrv.DCI, xarg *[4]byte, r image.Rectangle) {
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
}

func pixset(dci tftdrv.DCI, oldpf *[1]byte, newpf byte) {
	if oldpf[0] != newpf {
		oldpf[0] = newpf
		dci.Cmd(PIXSET)
		dci.WriteBytes(oldpf[:])
	}
}
