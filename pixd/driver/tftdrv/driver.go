// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"image/color"
	"image/draw"
	"time"
)

// Driver implements pixd.Driver interface with a limited support for draw.Over
// operation. It uses write-only DCI so the alpha blending is slow and reduced
// to 1-bit resolution. Use DriverOver if the full-fledged Porter-Duff
// composition is required and the display supports reading from its frame
// memory.
type Driver struct {
	dci        DCI
	startWrite AccessRAM
	pixSet     PixSet
	w, h       uint16
	cfast      uint16
	cinfo      byte
	pf         [1]byte
	xarg       [4]byte
	pf16       byte
	pf18       byte
	buf        [38 * 3]byte // must be multiple of two and three
} // ont 32-bit MCU the size of this struct is 144 B, exactly full allocation unit (see runtime/sizeclasses_mcu.go)

// New returns new Driver.
func New(dci DCI, w, h uint16, startWrite AccessRAM, pixSet PixSet, pf16, pf18 byte) *Driver {
	return &Driver{
		dci:        dci,
		startWrite: startWrite,
		pixSet:     pixSet,
		w:          w,
		h:          h,
		pf16:       pf16,
		pf18:       pf18,
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
	i := 0
	for i < len(cmds) {
		cmd := cmds[i]
		n := int(cmds[i+1])
		i += 2
		if n == 255 {
			time.Sleep(time.Duration(cmd) * time.Millisecond)
			continue
		}
		d.dci.Cmd(cmd)
		if n != 0 {
			k := i + n
			data := cmds[i:k]
			i = k
			d.dci.WriteBytes(data)
		}
	}
	d.dci.End()
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
	// best color format supported is 18-bit RGB 666
	r &^= 3
	g &^= 3
	b &^= 3
	if d.pf16 != 0 && r&7 == 0 && b&7 == 0 {
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
	pf := d.pf16
	if d.cinfo>>osize&3 == 3 {
		pf = d.pf18
	}
	d.pixSet(d.dci, &d.pf, pf)
	d.startWrite(d.dci, &d.xarg, r)

	pixSize := int(d.cinfo>>osize) & 3
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
	dst := ddram{d.dci, r.Size(), 3}
	getBuf := func() []byte {
		if d.cinfo&(bufInit<<otype) != 0 {
			d.cinfo &^= (bufFull ^ bufInit) << otype // inform Fill about dirty buf
			return d.buf[d.cinfo>>osize&3:]
		}
		return d.buf[:]
	}
	if op == draw.Src {
		pf := d.pf18
		if mask == nil && sip.pixSize == 2 {
			pf = d.pf16
			dst.pixSize = 2
		}
		d.pixSet(d.dci, &d.pf, pf)
		d.startWrite(d.dci, &d.xarg, r)
		drawSrc(dst, src, sp, sip, mask, mp, getBuf, len(d.buf)*3/4)
	} else {
		d.pixSet(d.dci, &d.pf, d.pf18)
		startWrite := func(r image.Rectangle) {
			d.startWrite(d.dci, &d.xarg, r)
		}
		drawOverNoRead(dst, r.Min, src, sp, sip, mask, mp, getBuf(), startWrite)
	}
	d.dci.End()
}
