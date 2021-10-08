// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/internal"
)

type Driver struct {
	dci   tftdrv.DCI
	xarg  [4]byte
	pf    [1]byte
	cinfo byte
	cfast uint16
	w, h  uint16
	buf   [32 * 3]byte // must be multiple of two and three
}

// Max. SPI clock speed: write: 100 ns (10 MHz), read: 150 ns (6.7 MHz).
// It seems that it works well even with 20 MHz clock.

func New(dci tftdrv.DCI) *Driver {
	return &Driver{dci: dci, w: 240, h: 320}
}

func (d *Driver) DCI() tftdrv.DCI      { return d.dci }
func (d *Driver) Err(clear bool) error { return d.dci.Err(clear) }
func (d *Driver) Flush()               {}

func (d *Driver) Size() image.Point {
	return image.Point{int(d.w), int(d.h)}
}

var initCmds = [...][]byte{
	{SFD, 0x03, 0x80, 0x02}, // {0xCA, 0xC3, 0x08, 0x50}
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
	{PGAMCTRL, 0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, 0x4E, 0xF1, 0x37, 0x07, 0x10,
		0x03, 0x0E, 0x09, 0x00},
	{NGAMCTRL, 0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, 0x31, 0xC1, 0x48, 0x08, 0x0F,
		0x0C, 0x31, 0x36, 0x0F},
}

func (d *Driver) Init(swreset bool) {
	if swreset {
		d.dci.Cmd(SWRESET)
	}
	resetTime := time.Now()
	time.Sleep(5 * time.Millisecond)
	for _, cmd := range initCmds {
		d.dci.Cmd(cmd[0])
		d.dci.WriteBytes(cmd[1:])
	}
	time.Sleep(resetTime.Add(120 * time.Millisecond).Sub(time.Now()))
	d.dci.Cmd(SLPOUT)
	time.Sleep(5 * time.Millisecond)
	d.dci.Cmd(DISPON)
}

func (d *Driver) SetMADCTL(madctl byte) {
	d.dci.Cmd(MADCTL)
	d.xarg[0] = madctl
	d.dci.WriteBytes(d.xarg[:1])
}

const (
	transparent = 0

	osize    = 4
	otype    = 6 // Fill relies on the type field takes up two MSbits
	fastByte = 0
	fastWord = 1
	bufInit  = 2 // getBuf relies on the one bit difference to the bufFull
	bufFull  = 3 // Fill relies on the both bits set
)

func (d *Driver) SetColor(c color.Color) {
	var r, g, b uint32
	switch cc := c.(type) {
	case color.RGBA:
		if cc.A>>7 == 0 {
			d.cinfo = transparent // only 1-bit transparency is supported
			return
		}
		r = uint32(cc.R)
		g = uint32(cc.G)
		b = uint32(cc.B)
	default:
		var a uint32
		r, g, b, a = c.RGBA()
		if a>>15 == 0 {
			d.cinfo = transparent // only 1-bit transparency is supported
			return
		}
		r >>= 8
		g >>= 8
		b >>= 8
	}
	// best color format supported is 18-bit RGB 666
	r &^= 3
	g &^= 3
	b &^= 3
	if r&7 == 0 && b&7 == 0 {
		rgb565 := r<<8 | g<<3 | b>>3
		if _, ok := d.dci.(tftdrv.WordNWriter); ok {
			d.cinfo = fastWord<<otype | 1<<osize | MCU16
			d.cfast = uint16(rgb565)
			return
		}
		h := rgb565 >> 8
		l := rgb565 & 0xff
		if h == l {
			if _, ok := d.dci.(tftdrv.ByteNWriter); ok {
				d.cinfo = fastByte<<otype | 2<<osize | MCU16
				d.cfast = uint16(h)
				return
			}
		}
		d.cinfo = bufInit<<otype | 2<<osize | MCU16
		d.buf[0] = byte(h)
		d.buf[1] = byte(l)
		return
	}
	if r == g && g == b {
		if _, ok := d.dci.(tftdrv.ByteNWriter); ok {
			d.cfast = uint16(r)
			d.cinfo = fastByte<<otype | 3<<osize | MCU18
			return
		}
	}
	d.cinfo = bufInit<<otype | 3<<osize | MCU18
	d.buf[0] = uint8(r)
	d.buf[1] = uint8(g)
	d.buf[2] = uint8(b)
}

func capaset(d *Driver, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	d.dci.Cmd(CASET)
	d.xarg[0] = uint8(r.Min.X >> 8)
	d.xarg[1] = uint8(r.Min.X)
	d.xarg[2] = uint8(r.Max.X >> 8)
	d.xarg[3] = uint8(r.Max.X)
	d.dci.WriteBytes(d.xarg[:4])
	d.dci.Cmd(PASET)
	d.xarg[0] = uint8(r.Min.Y >> 8)
	d.xarg[1] = uint8(r.Min.Y)
	d.xarg[2] = uint8(r.Max.Y >> 8)
	d.xarg[3] = uint8(r.Max.Y)
	d.dci.WriteBytes(d.xarg[:4])
}

func pixset(d *Driver, pf byte) {
	if d.pf[0] != pf {
		d.pf[0] = pf
		d.dci.Cmd(PIXSET)
		d.dci.WriteBytes(d.pf[:])
	}
}

func (d *Driver) Fill(r image.Rectangle) {
	if d.cinfo == transparent {
		return
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	capaset(d, r)
	pixset(d, d.cinfo&0xf)
	d.dci.Cmd(RAMWR)

	pixSize := int(d.cinfo>>osize) & 3
	n *= pixSize
	switch d.cinfo >> otype {
	case fastWord:
		d.dci.(tftdrv.WordNWriter).WriteWordN(d.cfast, n)
		return
	case fastByte:
		d.dci.(tftdrv.ByteNWriter).WriteByteN(byte(d.cfast), n)
		return
	case bufInit:
		d.cinfo |= bufFull << otype
		for i := pixSize; i < len(d.buf); i += pixSize {
			d.buf[i+0] = d.buf[0]
			d.buf[i+1] = d.buf[1]
			if pixSize == 3 {
				d.buf[i+2] = d.buf[2]
			}
		}
	}
	m := len(d.buf)
	for {
		if m > n {
			m = n
		}
		d.dci.WriteBytes(d.buf[:m])
		n -= m
		if n == 0 {
			break
		}
	}
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	sip := internal.ImageAtPoint(src, sp)
	dst := internal.DDRAM{d.dci, r.Size(), 3}
	getBuf := func() []byte {
		if d.cinfo&(bufInit<<otype) != 0 {
			d.cinfo &^= (bufFull ^ bufInit) << otype // inform Fill about dirty buf
			return d.buf[d.cinfo>>osize&3:]
		}
		return d.buf[:]
	}
	if op == draw.Src {
		capaset(d, r)
		pf := byte(MCU18)
		if mask == nil && sip.PixSize == 2 {
			pf = MCU16
			dst.PixSize = 2
		}
		pixset(d, pf)
		d.dci.Cmd(RAMWR)
		internal.DrawSrc(dst, src, sp, sip, mask, mp, getBuf, len(d.buf)*3/4)
	} else {
		pixset(d, MCU18)
		wrRect := func(r image.Rectangle) { capaset(d, r); d.dci.Cmd(RAMWR) }
		internal.DrawOverNoRead(dst, r.Min, src, sp, sip, mask, mp, getBuf(), wrRect)
	}

}
