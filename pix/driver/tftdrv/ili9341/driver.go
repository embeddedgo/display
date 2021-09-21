// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
)

type Driver struct {
	dci    tftdrv.DCI
	rgb16  int32
	pixset [1]byte
	xarg   [4]byte
	rgb    [13 * 3]uint8
}

func New(dci tftdrv.DCI) *Driver {
	return &Driver{dci: dci}
}

func (d *Driver) DCI() tftdrv.DCI          { return d.dci }
func (d *Driver) Dim() (width, height int) { return 240, 320 }
func (d *Driver) Err(clear bool) error     { return d.dci.Err(clear) }
func (d *Driver) Flush()                   {}

var initCmds = [...][]byte{
	{0xEF, 0x03, 0x80, 0x02}, // {0xCA, 0xC3, 0x08, 0x50}
	{PWCTRB, 0x00, 0xC1, 0x30},
	{PONSEQ, 0x64, 0x03, 0x12, 0x81},
	{DRVTIM, 0x85, 0x00, 0x78},
	{PWCTRA, 0x39, 0x2C, 0x00, 0x34, 0x02},
	{PUMPRT, 0x20},
	{DRVTIMB, 0x00, 0x00},
	{PWCTR1, 0x23},
	{PWCTR2, 0x10},
	{VMCTR1, 0x3e, 0x28},
	{VMCTR2, 0x86},
	{VSCRSADD, 0x00},
	{FRMCTR1, 0x00, 0x18},
	{DFUNCTR, 0x08, 0x82, 0x27},
	{GAMMASET, 0x01},
	{GMCTRP1, 0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1, 0x37, 0x07, 0x10,
		0x03, 0x0E, 0x09, 0x00},
	{GMCTRN1, 0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1, 0x48, 0x08, 0x0F,
		0x0C, 0x31, 0x36, 0x0F},
}

func (d *Driver) Init(madctl uint8, swreset bool) {
	resetTime := time.Now()
	time.Sleep(5 * time.Millisecond)
	dci := d.dci
	for _, cmd := range initCmds {
		dci.Cmd(cmd[0])
		dci.WriteBytes(cmd[1:])
	}
	dci.Cmd(MADCTL)
	d.xarg[0] = madctl
	dci.WriteBytes(d.xarg[:1])
	time.Sleep(resetTime.Add(120 * time.Millisecond).Sub(time.Now()))
	dci.Cmd(SLPOUT)
	time.Sleep(5 * time.Millisecond)
	dci.Cmd(DISPON)
}

func (d *Driver) SetColor(c color.Color) {
	switch cc := c.(type) {
	case pix.RGB16:
		if _, ok := d.dci.(tftdrv.WordNWriter); ok {
			d.rgb16 = int32(cc)
			return
		}
		d.rgb[0] = uint8(cc >> 11)
		d.rgb[1] = uint8(cc >> 5 & 0x3f)
		d.rgb[2] = uint8(cc & 0x1f)
	case pix.RGB:
		d.rgb[0] = cc.R
		d.rgb[1] = cc.G
		d.rgb[2] = cc.B
	case color.RGBA:
		if cc.A>>7 == 0 {
			d.rgb16 = -1 // transparent color, only 1-bit transparency supported
			return
		}
		d.rgb[0] = cc.R
		d.rgb[1] = cc.G
		d.rgb[2] = cc.B
	default:
		r, g, b, a := c.RGBA()
		if a>>15 == 0 {
			d.rgb16 = -1 // transparent color, only 1-bit transparency supported
			return
		}
		d.rgb[0] = uint8(r >> 8)
		d.rgb[1] = uint8(g >> 8)
		d.rgb[2] = uint8(b >> 8)
	}
	for i := 3; i < len(d.rgb); i += 3 {
		d.rgb[i] = d.rgb[0]
		d.rgb[i+1] = d.rgb[1]
		d.rgb[i+2] = d.rgb[2]
	}
	d.rgb16 = -2
}

func pixset(d *Driver) {
	pixset := false
	if d.rgb16 >= 0 {
		if pixset = d.pixset[0] != PF16; pixset {
			d.pixset[0] = PF16
		}
	} else {
		if pixset = d.pixset[0] != PF18; pixset {
			d.pixset[0] = PF18
		}
	}
	if pixset {
		d.dci.Cmd(PIXSET)
		d.dci.WriteBytes(d.pixset[:])
	}
}

func (d *Driver) Fill(r image.Rectangle) {
	if d.rgb16 == -1 {
		return // transparent color
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	pixset(d)
	dci := d.dci
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	d.xarg[0] = uint8(r.Min.X >> 8)
	d.xarg[1] = uint8(r.Min.X)
	d.xarg[2] = uint8(r.Max.X >> 8)
	d.xarg[3] = uint8(r.Max.X)
	dci.WriteBytes(d.xarg[:4])
	dci.Cmd(PASET)
	d.xarg[0] = uint8(r.Min.Y >> 8)
	d.xarg[1] = uint8(r.Min.Y)
	d.xarg[2] = uint8(r.Max.Y >> 8)
	d.xarg[3] = uint8(r.Max.Y)
	dci.WriteBytes(d.xarg[:4])
	dci.Cmd(RAMWR)
	if d.rgb16 >= 0 {
		dci.(tftdrv.WordNWriter).WriteWordN(uint16(d.rgb16), n)
		return
	}
	n *= 3
	if d.rgb[0] == d.rgb[1] && d.rgb[1] == d.rgb[2] {
		if w, ok := dci.(tftdrv.ByteNWriter); ok {
			w.WriteByteN(d.rgb[0], n)
			return
		}
	}
	m := len(d.rgb)
	for {
		if m > n {
			m = n
		}
		dci.WriteBytes(d.rgb[:m])
		n -= m
		if n == 0 {
			break
		}
	}
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	if op == draw.Src {

	}
}
