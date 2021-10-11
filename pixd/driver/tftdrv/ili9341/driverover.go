// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/embeddedgo/display/pixd/driver/tftdrv"
)

// DriverOver implements pixd.Driver interface with the full support for
// draw.Over operator. It requires tftdrv.RDCI to read the frame memory content.
// If the display has write-only interface use Driver instead.
type DriverOver struct {
	dci     tftdrv.RDCI
	xarg    [4]byte
	pf      [1]byte
	cinfo   byte
	cfast   uint16
	r, g, b uint16
	w, h    uint16
	buf     [32 * 3]byte // must be multiple of two and three
}

// NewOver returns new DriverOver.
func NewOver(dci tftdrv.RDCI) *DriverOver {
	return &DriverOver{dci: dci, w: 240, h: 320}
}

func (d *DriverOver) DCI() tftdrv.RDCI     { return d.dci }
func (d *DriverOver) Err(clear bool) error { return d.dci.Err(clear) }
func (d *DriverOver) Flush()               {}

func (d *DriverOver) Size() image.Point {
	return image.Point{int(d.w), int(d.h)}
}

// Init initializes display using provided initialization commands. The
// initialization commands depends on the LCD pannel. See InitGFX and InitST for
// working examples.
func (d *DriverOver) Init(cmds []byte, swreset bool) {
	initialize(d.dci, cmds, swreset)
}

func (d *DriverOver) SetMADCTL(madctl byte) {
	d.dci.Cmd(MADCTL)
	d.xarg[0] = madctl
	d.dci.WriteBytes(d.xarg[:1])
}

const alphaOpaque = 0xfb00

func (d *DriverOver) SetColor(c color.Color) {
	var r, g, b, a uint32
	switch cc := c.(type) {
	case color.RGBA:
		if cc.A < 4 {
			d.cinfo = transparent
			return
		}
		r = uint32(cc.R) | uint32(cc.R)<<8
		g = uint32(cc.G) | uint32(cc.G)<<8
		b = uint32(cc.B) | uint32(cc.B)<<8
		a = uint32(cc.A) | uint32(cc.A)<<8
	default:
		var a uint32
		r, g, b, a = c.RGBA()
		if a < 0x0404 {
			d.cinfo = transparent
			return
		}
	}
	if a >= alphaOpaque {
		r >>= 8
		g >>= 8
		b >>= 8
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
			d.r = uint16(h)
			d.g = uint16(l)
			return
		}
		if r == g && g == b {
			if _, ok := d.dci.(tftdrv.ByteNWriter); ok {
				d.cfast = uint16(r)
				d.cinfo = fastByte<<otype | 3<<osize | MCU18
				return
			}
		}
	}
	d.cinfo = bufInit<<otype | 3<<osize | MCU18
	d.r = uint16(r)
	d.g = uint16(g)
	d.b = uint16(b)
	d.cfast = uint16(a)
}

func (d *DriverOver) Fill(r image.Rectangle) {
	if d.cinfo == transparent {
		return
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	pixset(d.dci, &d.pf, d.cinfo&0xf)
	capaset(d.dci, &d.xarg, r)
	d.dci.Cmd(RAMWR)
	pixSize := int(d.cinfo>>osize) & 3
	n *= pixSize
	typ := d.cinfo >> otype
	switch {
	case typ == fastWord:
		d.dci.(tftdrv.WordNWriter).WriteWordN(d.cfast, n)
		print("w")
	case typ == fastByte:
		d.dci.(tftdrv.ByteNWriter).WriteByteN(byte(d.cfast), n)
		print("b")
	case d.cfast >= alphaOpaque:
		if typ == bufInit {
			d.cinfo |= bufFull << otype
			for i := 0; i < len(d.buf); i += pixSize {
				d.buf[i+0] = uint8(d.r)
				d.buf[i+1] = uint8(d.g)
				if pixSize == 3 {
					d.buf[i+2] = uint8(d.b)
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
		print("o")
	default:
		sr := uint(d.r)
		sg := uint(d.g)
		sb := uint(d.b)
		a := 0xffff - uint(d.cfast)
		for y := r.Min.Y; y < r.Max.Y; y++ {
			x := r.Min.X
			for {
				width := r.Max.X - x
				if width <= 0 {
					break
				}
				if width > len(d.buf)/3 {
					width = len(d.buf) / 3
				}
				r1 := image.Rectangle{
					image.Point{x, y},
					image.Point{x + width, y + 1},
				}
				x += width
				width *= 3
				capaset(d.dci, &d.xarg, r1)
				d.dci.Cmd(RAMRD)
				d.dci.ReadBytes(d.buf[:width])
				d.dci.End() // required to end RAMRD (undocumented)
				for i := 0; i < width; i += 3 {
					r := uint(d.buf[i+0])
					g := uint(d.buf[i+1])
					b := uint(d.buf[i+2])
					r = (r<<8|r)*a/0xffff + sr
					g = (g<<8|g)*a/0xffff + sg
					b = (b<<8|b)*a/0xffff + sb
					d.buf[i+0] = uint8(r >> 8)
					d.buf[i+1] = uint8(g >> 8)
					d.buf[i+2] = uint8(b >> 8)
				}
				capaset(d.dci, &d.xarg, r1)
				d.dci.Cmd(RAMWR)
				d.dci.WriteBytes(d.buf[:width])
			}
		}
	}
	d.dci.End()
}

func (d *DriverOver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
}
