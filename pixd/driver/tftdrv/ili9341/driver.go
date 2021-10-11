// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/internal"
)

// Driver implements pixd.Driver interface with the limited support for
// draw.Over operator. It uses write-only tftdrv.DCI so the alpha blending is
// slow and reduced to 1-bit resolution. Use DriverOver if the full-fledged
// Porter-Duff composition is required and the display supports reading from
// its frame memory.
type Driver struct {
	dci   tftdrv.DCI
	xarg  [4]byte
	pf    [1]byte
	cinfo byte
	cfast uint16
	w, h  uint16
	buf   [32 * 3]byte // must be multiple of two and three
}

// New returns new Driver.
func New(dci tftdrv.DCI) *Driver {
	return &Driver{dci: dci, w: 240, h: 320}
}

func (d *Driver) DCI() tftdrv.DCI      { return d.dci }
func (d *Driver) Err(clear bool) error { return d.dci.Err(clear) }
func (d *Driver) Flush()               {}

func (d *Driver) Size() image.Point {
	return image.Point{int(d.w), int(d.h)}
}

// Init initializes display using provided initialization commands. The
// initialization commands depends on the LCD pannel. See InitGFX and InitST for
// working examples.
func (d *Driver) Init(cmds []byte, swreset bool) {
	initialize(d.dci, cmds, swreset)
}

func (d *Driver) SetMADCTL(madctl byte) {
	d.dci.Cmd(MADCTL)
	d.xarg[0] = madctl
	d.dci.WriteBytes(d.xarg[:1])
}

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

func (d *Driver) startWrite(r image.Rectangle) {
	capaset(d.dci, &d.xarg, r)
	d.dci.Cmd(RAMWR)
}

func (d *Driver) Fill(r image.Rectangle) {
	if d.cinfo == transparent {
		return
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	pixset(d.dci, &d.pf, d.cinfo&0xf)
	d.startWrite(r)

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
		pf := byte(MCU18)
		if mask == nil && sip.PixSize == 2 {
			pf = MCU16
			dst.PixSize = 2
		}
		pixset(d.dci, &d.pf, pf)
		d.startWrite(r)
		internal.DrawSrc(dst, src, sp, sip, mask, mp, getBuf, len(d.buf)*3/4)
	} else {
		pixset(d.dci, &d.pf, MCU18)
		internal.DrawOverNoRead(dst, r.Min, src, sp, sip, mask, mp, getBuf(), d.startWrite)
	}

}
