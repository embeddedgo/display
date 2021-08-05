// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tft

import (
	"image"
	"image/color"
	"time"
)

type ILI9341 struct{ DCI DCI }

func (tft ILI9341) Init(madctl uint8, swreset bool) {
	resetTime := time.Now()
	time.Sleep(5 * time.Millisecond)

	dci := tft.DCI

	dci.Cmd(0xEF)
	dci.WriteBytes(0x03, 0x80, 0x02)
	//dci.Cmd(0xCA)
	//dci.WriteBytes(0xC3, 0x08, 0x50)

	dci.Cmd(PWCTRB)
	dci.WriteBytes(0x00, 0xC1, 0x30)

	dci.Cmd(0xED)
	dci.WriteBytes(0x64, 0x03, 0x12, 0x81)

	dci.Cmd(0xE8)
	dci.WriteBytes(0x85, 0x00, 0x78)

	dci.Cmd(0xCB)
	dci.WriteBytes(0x39, 0x2C, 0x00, 0x34, 0x02)

	dci.Cmd(0xF7)
	dci.WriteBytes(0x20)

	dci.Cmd(0xEA)
	dci.WriteBytes(0x00, 0x00)

	dci.Cmd(PWCTR1)
	dci.WriteBytes(0x23)

	dci.Cmd(PWCTR2)
	dci.WriteBytes(0x10)

	dci.Cmd(VMCTR1)
	dci.WriteBytes(0x3e, 0x28)

	dci.Cmd(VMCTR2)
	dci.WriteBytes(0x86)

	dci.Cmd(VSCRSADD)
	dci.WriteBytes(0x00)

	dci.Cmd(FRMCTR1)
	dci.WriteBytes(0x00, 0x18)

	dci.Cmd(DFUNCTR)
	dci.WriteBytes(0x08, 0x82, 0x27)

	dci.Cmd(GAMMASET)
	dci.WriteBytes(0x01)

	dci.Cmd(GMCTRP1)
	dci.WriteBytes(0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00)

	dci.Cmd(GMCTRN1)
	dci.WriteBytes(0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F)

	dci.Cmd(MADCTL)
	dci.WriteBytes(madctl)

	dci.Cmd(PIXSET)
	dci.WriteBytes(PF16)

	time.Sleep(resetTime.Add(120 * time.Millisecond).Sub(time.Now()))

	dci.Cmd(SLPOUT)

	time.Sleep(5 * time.Millisecond)

	dci.Cmd(DISPON)
	dci.Cmd(RAMWR)
}

func (tft ILI9341) DriverColor(c color.Color) uint64 {
	r, g, b, _ := c.RGBA()
	return uint64((r&0x1f)<<11 | (g&0x3f)<<5 | (b & 0x1f))
}

func (tft ILI9341) FillRect(r image.Rectangle, color uint64) {
	dci.Cmd(CASET)
	dci.WriteWords(uint16(r.Min.X), uint16(r.Max.X-1))
	dci.Cmd(PASET)
	dci.WriteWords(uint16(r.Min.Y), uint16(r.Max.Y-1))
	dci.Cmd(RAMWR)
	dci.WriteWordN(uint16(color), r.Dx()*r.Dy())
}

func (tft ILI9341) Draw(at image.Point, img image.Image) {
	// TODO
}
