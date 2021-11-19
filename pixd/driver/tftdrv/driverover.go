// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"image/color"
	"image/draw"
)

// magic numbers
const (
	sa     = 2 // must be smalest
	sb     = 5
	sc     = 7
	se     = 4                   // must be >= 1
	bufLen = (sa*sb*sc + se) * 3 // must be multiple of 2 and 3
)

var bufDim = [...]uint16{
	sa*sb*sc<<8 | 1, // bestBufSize requires one row here
	1<<8 | sa*sb*sc, // bestBufSize requires one column here
	sa<<8 | sb*sc,
	sb*sc<<8 | sa,
	sb<<8 | sa*sc,
	sa*sc<<8 | sb,
	sc<<8 | sa*sb,
	sa*sb<<8 | sc,
}

// bestBufSize finds the best buffer dimensions to cover the width x height
// rectangle
func bestBufSize(rsiz image.Point) image.Point {
	var best image.Point
	if rsiz.X < sa || rsiz.Y < sa {
		// fast path for hline and vline
		if rsiz.X >= rsiz.Y {
			best = image.Pt(int(bufDim[0])>>8, int(bufDim[0])&0xff)
		} else {
			best = image.Pt(int(bufDim[1])>>8, int(bufDim[1])&0xff)
		}
	} else {
		bu := rsiz.X * rsiz.Y
		for _, dim := range bufDim {
			dw := int(dim) >> 8
			dh := int(dim) & 0xff
			nx := rsiz.X / dw
			ny := rsiz.Y / dh
			ux := rsiz.X - nx*dw
			if ux != 0 {
				ux = ny // we do not pay attention to the size
			}
			uy := rsiz.Y - ny*dh
			if uy != 0 {
				uy = nx // we do not pay attention to the size
			}
			if uc := uy + ux; uc < bu {
				bu = uc
				best = image.Pt(dw, dh)
			}
		}
	}
	return best
}

// BUG: we assume that any controller supports 24-bit pixel data format

// DriverOver implements pixd.Driver interface with the full support for
// draw.Over operator. The DCI must fully implement ReadBytes method to read the
// frame memory content. If the display has write-only interface use Driver
// instead.
type DriverOver struct {
	dci     DCI
	c       *Ctrl
	w, h    uint16
	r, g, b uint16
	cfast   uint16
	cinfo   byte
	pf      PF
	mv      byte
	dir     [1]byte
	parg    [1]byte
	xarg    [4]byte
	buf     [bufLen]byte
} // ont 32-bit MCU the size of this struct is 255 B (bufLen=222), almost full 256 B allocation unit (see runtime/sizeclasses_mcu.go)

// NewOver returns new DriverOver.
func NewOver(dci DCI, w, h uint16, pf PF, ctrl *Ctrl) *DriverOver {
	return &DriverOver{
		dci: dci,
		c:   ctrl,
		w:   w,
		h:   h,
		pf:  pf,
	}
}

func (d *DriverOver) Err(clear bool) error { return d.dci.Err(clear) }
func (d *DriverOver) Flush()               {}
func (d *DriverOver) Size() image.Point    { return image.Pt(int(d.w), int(d.h)) }

// Init initializes the display using provided initialization commands. The
// initialization commands depends on the LCD pannel. The command that sets
// the display orientation and the color order must be the last one in the cmds
// See ili9341.GFX, ili9486.MSP4022 as examples.
func (d *DriverOver) Init(cmds []byte) {
	initialize(d.dci, cmds)
	d.dir[0] = cmds[len(cmds)-1]
}

func (d *DriverOver) SetDir(dir int) image.Rectangle {
	if d.c.SetDir != nil {
		if mv := byte(dir & 1); mv != d.mv {
			d.mv = mv
			d.w, d.h = d.h, d.w
		}
		d.c.SetDir(d.dci, &d.parg, &d.dir, dir)
	}
	return image.Rectangle{Max: image.Pt(int(d.w), int(d.h))}
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
	dstSize := r.Size()
	n := dstSize.X * dstSize.Y
	if n == 0 {
		return
	}
	pixSize := int(d.cinfo>>osize) & 3
	if d.c.SetPF != nil {
		d.c.SetPF(d.dci, &d.parg, pixSize)
	}
	if typ := d.cinfo >> otype; typ < bufInit || d.cfast >= alphaOpaque {
		// no alpha blending
		d.c.StartWrite(d.dci, &d.xarg, r)
		n *= pixSize
		switch {
		case typ == fastWord:
			d.dci.(WordNWriter).WriteWordN(d.cfast, n)
		case typ == fastByte:
			d.dci.(ByteNWriter).WriteByteN(byte(d.cfast), n)
		default:
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
		}
	} else {
		// alpha blending with the current display content
		bufSize := bestBufSize(dstSize)
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
			if height > bufSize.Y {
				height = bufSize.Y
			}
			x := r.Min.X
			for {
				width := r.Max.X - x
				if width <= 0 {
					break
				}
				if width > bufSize.X {
					width = bufSize.X
				}
				r1 := image.Rectangle{
					image.Pt(x, y),
					image.Pt(x+width, y+height),
				}
				x += width
				n := width*height*3 + 1
				d.c.Read(d.dci, &d.xarg, r1, d.buf[0:n])
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
				d.c.StartWrite(d.dci, &d.xarg, r1)
				d.dci.WriteBytes(d.buf[1:n])
			}
			y += height
		}
	}
	d.dci.End()
}

func (d *DriverOver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	getBuf := func() []byte {
		if d.cinfo&(bufFull<<otype) == bufFull<<otype {
			// inform Fill about dirty buf
			d.cinfo &^= (bufFull ^ bufInit) << otype
		}
		return d.buf[:]
	}
	dst := dst{r.Size(), 3}
	sip := imageAtPoint(src, sp)
	if op == draw.Src {
		if d.c.SetPF != nil {
			if mask == nil && sip.pixSize <= 3 {
				dst.pixSize = sip.pixSize
			}
			d.c.SetPF(d.dci, &d.parg, dst.pixSize)
		}
		d.c.StartWrite(d.dci, &d.xarg, r)
		drawSrc(d.dci, dst, src, sp, sip, mask, mp, getBuf, len(d.buf)*3/4)
	} else {
		if d.c.SetPF != nil {
			d.c.SetPF(d.dci, &d.parg, 3)
		}
		buf := getBuf()
		bufSize := bestBufSize(dst.size)
		y := 0
		for {
			height := dst.size.Y - y
			if height <= 0 {
				break
			}
			if height > bufSize.Y {
				height = bufSize.Y
			}
			var r1 image.Rectangle
			r1.Min.Y = r.Min.Y + y
			r1.Max.Y = r1.Min.Y + height
			x := 0
			for {
				width := dst.size.X - x
				if width <= 0 {
					break
				}
				if width > bufSize.X {
					width = bufSize.X
				}
				r1.Min.X = r.Min.X + x
				r1.Max.X = r1.Min.X + width
				n := width*height*3 + 1
				d.c.Read(d.dci, &d.xarg, r1, buf[0:n])
				i := 1
				for y1 := y; y1 < y+height; y1++ {
					j := y1*sip.stride + x*sip.pixSize
					for x1 := x; x1 < x+width; x1++ {
						var sr, sg, sb, sa uint32
						if sip.pixSize != 0 {
							sr, sg, sb, sa = fastRGBA(&sip, j)
							j += sip.pixSize
						} else {
							sr, sg, sb, sa = src.At(sp.X+x1, sp.Y+y1).RGBA()
						}
						ma := uint32(0xffff)
						if mask != nil {
							_, _, _, ma = mask.At(mp.X+x1, mp.Y+y1).RGBA()
						}
						a := 0xffff - (sa * ma / 0xffff)
						dr := uint32(buf[i+0])
						dg := uint32(buf[i+1])
						db := uint32(buf[i+2])
						dr = ((dr<<8|dr)*a + sr*ma) / 0xffff
						dg = ((dg<<8|dg)*a + sg*ma) / 0xffff
						db = ((db<<8|db)*a + sb*ma) / 0xffff
						buf[i+0] = uint8(dr >> 8)
						buf[i+1] = uint8(dg >> 8)
						buf[i+2] = uint8(db >> 8)
						i += 3
					}
				}
				d.c.StartWrite(d.dci, &d.xarg, r1)
				d.dci.WriteBytes(buf[1:n])
				x += width
			}
			y += height
		}
	}
	d.dci.End()
}
