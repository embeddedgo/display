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

// magic numbers
const (
	sa     = 2 // must be smalest
	sb     = 3
	sc     = 11
	se     = 4                   // must be >= 1
	bufLen = (sa*sb*sc + se) * 3 // must be multiple of 2 and 3
)

var bufDim = [...]uint16{
	sa*sb*sc<<8 | 1, // Fill requires one row here
	1<<8 | sa*sb*sc, // Fill requires one column here
	sa<<8 | sb*sc,
	sb*sc<<8 | sa,
	sb<<8 | sa*sc,
	sa*sc<<8 | sb,
	sc<<8 | sa*sb,
	sa*sb<<8 | sc,
}

// BUG: we assume that any controller supports 24-bit pixel data format

// DriverOver implements pixd.Driver interface with the full support for
// draw.Over operator. It requires tftdrv.RDCI to read the frame memory content.
// If the display has write-only interface use Driver instead.
type DriverOver struct {
	dci        RDCI
	startRead  AccessFrame
	startWrite AccessFrame
	pixSet     PixSet
	setDir     PixSet
	w, h       uint16
	r, g, b    uint16
	cfast      uint16
	cinfo      byte
	pdf        PDF
	parg       [1]byte
	xarg       [4]byte
	buf        [bufLen]byte
} // ont 32-bit MCU the size of this struct is 253 B (bufLen=210), almost full 256 B allocation unit (see runtime/sizeclasses_mcu.go)

// NewOver returns new DriverOver.
func NewOver(dci RDCI, w, h uint16, pdf PDF, startRead, startWrite AccessFrame, pixSet, setDir PixSet) *DriverOver {
	return &DriverOver{
		dci:        dci,
		startRead:  startRead,
		startWrite: startWrite,
		pixSet:     pixSet,
		setDir:     setDir,
		w:          w,
		h:          h,
		pdf:        pdf,
	}
}

func (d *DriverOver) Err(clear bool) error { return d.dci.Err(clear) }
func (d *DriverOver) Flush()               {}
func (d *DriverOver) Size() image.Point    { return image.Pt(int(d.w), int(d.h)) }

// Init initializes display using provided initialization commands. The
// initialization commands depends on the LCD pannel. See InitGFX for working
// example.
func (d *DriverOver) Init(cmds []byte) {
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

func (d *DriverOver) SetDir(dir int) {}

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
		if d.pdf&W24 == 0 {
			// clear two LS-bits to increase the chances of Byte/WordNWriter
			r &^= 3
			g &^= 3
			b &^= 3
		}
		if d.pdf&W16 != 0 && r&7 == 0 && b&7 == 0 {
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
			d.r = uint16(h)
			d.g = uint16(l)
			return
		}
		if r == g && g == b {
			if _, ok := d.dci.(ByteNWriter); ok {
				d.cfast = uint16(r)
				d.cinfo = fastByte<<otype | 3<<osize
				return
			}
		}
	}
	d.cinfo = bufInit<<otype | 3<<osize
	d.r = uint16(r)
	d.g = uint16(g)
	d.b = uint16(b)
	d.cfast = uint16(a)
}

func (d *DriverOver) Fill(r image.Rectangle) {
	if d.cinfo == transparent {
		return
	}
	width, height := r.Dx(), r.Dy()
	n := width * height
	if n == 0 {
		return
	}
	pixSize := int(d.cinfo>>osize) & 3
	d.pixSet(d.dci, &d.parg, pixSize)
	d.startWrite(d.dci, &d.xarg, r)
	n *= pixSize
	typ := d.cinfo >> otype
	switch {
	case typ == fastWord:
		d.dci.(WordNWriter).WriteWordN(d.cfast, n)
	case typ == fastByte:
		d.dci.(ByteNWriter).WriteByteN(byte(d.cfast), n)
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
	default:
		// find the best coverage of the r area by d.buf
		var best image.Point
		if height < sa || width < sa {
			// fast path for hline and vline
			if width >= height {
				best = image.Pt(int(bufDim[0])>>8, int(bufDim[0])&0xff)
			} else {
				best = image.Pt(int(bufDim[1])>>8, int(bufDim[1])&0xff)
			}
		} else {
			bu := width * height
			for _, dim := range bufDim {
				dw := int(dim) >> 8
				dh := int(dim) & 0xff
				nx := width / dw
				ny := height / dh
				ux := width - nx*dw
				if ux != 0 {
					ux = ny // we do not pay attention to the size
				}
				uy := height - ny*dh
				if uy != 0 {
					uy = nx // we do not pay attention to the size
				}
				if uc := uy + ux; uc < bu {
					bu = uc
					best = image.Pt(dw, dh)
				}
			}
		}
		// draw
		sr := uint(d.r)
		sg := uint(d.g)
		sb := uint(d.b)
		a := 0xffff - uint(d.cfast)
		y := r.Min.Y
		for {
			height := r.Max.Y - y
			if height <= 0 {
				break
			}
			if height > best.Y {
				height = best.Y
			}
			x := r.Min.X
			for {
				width := r.Max.X - x
				if width <= 0 {
					break
				}
				if width > best.X {
					width = best.X
				}
				r1 := image.Rectangle{
					image.Pt(x, y),
					image.Pt(x+width, y+height),
				}
				x += width
				n := width*height*3 + 1
				d.startRead(d.dci, &d.xarg, r1)
				d.dci.ReadBytes(d.buf[0:n])
				d.dci.End() // required to end RAMRD (undocumented)
				for i := 1; i < n; i += 3 {
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
				d.startWrite(d.dci, &d.xarg, r1)
				d.dci.WriteBytes(d.buf[1:n])
			}
			y += height
		}
	}
	d.dci.End()
}

func (d *DriverOver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {

}
