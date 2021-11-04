// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"image/color"
	"image/draw"
)

// BUG: we assume that any controller supports 24-bit pixel data format

// Driver implements pixd.Driver interface with a limited support for draw.Over
// operation. It uses write-only DCI so the alpha blending is slow and reduced
// to 1-bit resolution. Use DriverOver if the full-fledged Porter-Duff
// composition is required and the display supports reading from its frame
// memory.
type Driver struct {
	dci   DCI
	c     *Ctrl
	w, h  uint16
	cfast uint16
	cinfo byte
	pf    PF
	parg  [1]byte
	xarg  [4]byte
	buf   [52 * 3]byte // must be multiple of two and three
} // ont 32-bit MCU the size of this struct is 189 B, almost full 192 B allocation unit (see runtime/sizeclasses_mcu.go)

// New returns new Driver.
func New(dci DCI, w, h uint16, pf PF, ctrl *Ctrl) *Driver {
	return &Driver{
		dci: dci,
		c:   ctrl,
		w:   w,
		h:   h,
		pf:  pf,
	}
}

func (d *Driver) Err(clear bool) error { return d.dci.Err(clear) }
func (d *Driver) Flush()               {}
func (d *Driver) Size() image.Point    { return image.Pt(int(d.w), int(d.h)) }

// Init initializes the display using provided initialization commands. The
// initialization commands depends on the LCD pannel. The command that sets
// the display orientation and the color order must be the last one in the cmds
// See ili9341.InitGFX for working example.
func (d *Driver) Init(cmds []byte) {
	initialize(d.dci, cmds)
}

func (d *Driver) SetDir(dir int) {}

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
	if d.pf&W24 == 0 {
		// clear two LS-bits to increase the chances of Byte/WordNWriter
		r &^= 3
		g &^= 3
		b &^= 3
	}
	if d.pf&W16 != 0 && r&7 == 0 && b&7 == 0 {
		rgb565 := r<<8 | g<<3 | b>>3
		if _, ok := d.dci.(WordNWriter); ok {
			d.cinfo = fastWord<<otype | 1<<osize
			d.cfast = uint16(rgb565)
			return
		}
		h := rgb565 >> 8
		l := rgb565 & 0xff
		if h == l {
			if _, ok := d.dci.(ByteNWriter); ok {
				d.cinfo = fastByte<<otype | 2<<osize
				d.cfast = uint16(h)
				return
			}
		}
		d.cinfo = bufInit<<otype | 2<<osize
		d.buf[0] = byte(h)
		d.buf[1] = byte(l)
		return
	}
	if r == g && g == b {
		if _, ok := d.dci.(ByteNWriter); ok {
			d.cfast = uint16(r)
			d.cinfo = fastByte<<otype | 3<<osize
			return
		}
	}
	d.cinfo = bufInit<<otype | 3<<osize
	d.buf[0] = uint8(r)
	d.buf[1] = uint8(g)
	d.buf[2] = uint8(b)
}

func (d *Driver) Fill(r image.Rectangle) {
	if d.cinfo == transparent {
		return
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	pixSize := int(d.cinfo>>osize) & 3
	if d.c.SetPF != nil {
		d.c.SetPF(d.dci, &d.parg, pixSize)
	}
	d.c.StartWrite(d.dci, &d.xarg, r)
	n *= pixSize
	switch d.cinfo >> otype {
	case fastWord:
		d.dci.(WordNWriter).WriteWordN(d.cfast, n)
		goto end
	case fastByte:
		d.dci.(ByteNWriter).WriteByteN(byte(d.cfast), n)
		goto end
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
	for m := len(d.buf); ; {
		if m > n {
			m = n
		}
		d.dci.WriteBytes(d.buf[:m])
		if n -= m; n == 0 {
			break
		}
	}
end:
	d.dci.End()
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	sip := imageAtPoint(src, sp)
	dst := framePart{d.dci, r.Size(), 3}
	getBuf := func() []byte {
		if d.cinfo&(bufInit<<otype) != 0 {
			// bufInit or bufFull
			d.cinfo &^= (bufFull ^ bufInit) << otype // inform Fill about dirty buf
			return d.buf[d.cinfo>>osize&3:]
		}
		return d.buf[:]
	}
	if op == draw.Src {
		if d.c.SetPF != nil {
			if mask == nil && sip.pixSize < dst.pixSize {
				dst.pixSize = sip.pixSize
			}
			d.c.SetPF(d.dci, &d.parg, dst.pixSize)
		}
		d.c.StartWrite(d.dci, &d.xarg, r)
		drawSrc(dst, src, sp, sip, mask, mp, getBuf, len(d.buf)*3/4)
	} else {
		if d.c.SetPF != nil {
			d.c.SetPF(d.dci, &d.parg, dst.pixSize)
		}
		startWrite := func(r image.Rectangle) {
			d.c.StartWrite(d.dci, &d.xarg, r)
		}
		drawOverNoRead(dst, r.Min, src, sp, sip, mask, mp, getBuf(), startWrite)
	}
	d.dci.End()
}
