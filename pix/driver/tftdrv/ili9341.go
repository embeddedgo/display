// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/embeddedgo/display/pix"
)

// The commands below can be used to directly interract with a display
// controller..
const (
	NOP      = 0x00
	SWRESET  = 0x01
	SLPIN    = 0x10
	SLPOUT   = 0x11
	GAMMASET = 0x26
	DISPOFF  = 0x28
	DISPON   = 0x29
	CASET    = 0x2A
	PASET    = 0x2B
	RAMWR    = 0x2C
	MADCTL   = 0x36
	VSCRSADD = 0x37
	PIXSET   = 0x3A
	FRMCTR1  = 0xB1
	DFUNCTR  = 0xB6
	PWCTR1   = 0xC0
	PWCTR2   = 0xC1
	VMCTR1   = 0xC5
	VMCTR2   = 0xC7
	PWCTRB   = 0xCF
	GMCTRP1  = 0xE0
	GMCTRN1  = 0xE1
)

// MADCTL arguments
const (
	MH  = 1 << 2 // horizontal refresh order
	BGR = 1 << 3 // BGR color filter panel
	ML  = 1 << 4 // vertical refresh order
	MV  = 1 << 5 // row/column exchange
	MX  = 1 << 6 // column address order
	MY  = 1 << 7 // row address order
)

// PIXSET arguments
const (
	PF16 = 0x55 // 16-bit 565 pixel format.
	PF18 = 0x66 // 18-bit 666 pixel format.
)

type ILI9341 struct {
	dci        DCI
	buf        [64]uint16
	r, g, b, a uint16
}

func NewILI9341(dci DCI) *ILI9341 {
	return &ILI9341{dci: dci}
}

func (d *ILI9341) DCI() DCI { return d.dci }

func (d *ILI9341) Init(madctl uint8, swreset bool) {
	resetTime := time.Now()
	time.Sleep(5 * time.Millisecond)

	dci := d.dci

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

func (d *ILI9341) Dim() (width, height int) {
	return 200, 320
}

func (d *ILI9341) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
}

func (d *ILI9341) SetColor(c color.Color) {
	if rgb16, ok := c.(pix.RGB16); ok {
		d.r = uint16(rgb16)
		d.a = 0xffff
		return
	}
	r, g, b, a := c.RGBA()
	if a>>8 == 0xff {
		d.r = uint16((r&0x1f)<<11 | (g&0x3f)<<5 | (b & 0x1f))
		d.a = 0xffff
	} else {
		d.r = uint16(r)
		d.g = uint16(g)
		d.b = uint16(b)
		d.a = uint16(a)
	}
}

func (d *ILI9341) Fill(r image.Rectangle) {
	n := r.Dx() * r.Dy()
	dci := d.dci
	if d.a == 0xffff {
		dci.Cmd(CASET)
		dci.WriteWords(uint16(r.Min.X), uint16(r.Max.X-1))
		dci.Cmd(PASET)
		dci.WriteWords(uint16(r.Min.Y), uint16(r.Max.Y-1))
		dci.Cmd(RAMWR)
		dci.WriteWordN(d.r, n)
		return
	}
}

func (d *ILI9341) Flush()               {}
func (d *ILI9341) Err(clear bool) error { return d.dci.Err(clear) }
