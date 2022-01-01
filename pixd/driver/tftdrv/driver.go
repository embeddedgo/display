// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"image/color"
	"image/draw"
)

// Driver implements pixd.Driver interface with a limited support for draw.Over
// operation. It is designed for write-only displays (doesn't use DCI.ReadBytes
// method) so the alpha blending is slow and reduced to 1-bit resolution. Use
// DriverOver if the full-fledged Porter-Duff composition is required and the
// display supports reading from its frame memory.
type Driver struct {
	dci     DCI
	c       *Ctrl
	w, h    uint16
	r, g, b uint16
	ctyp    byte
	csiz    int8
	cnpp    int8
	pf      PF
	rdir    [1]byte // current/initial value of direction relaed register
	rpf     [1]byte // current/initial value of pixel format relaed register
	xarg    [4]byte
	buf     [74 * 3]byte // must be multiple of two and three
} // ont 32-bit MCU the size of this struct is 254 B, almost full 256 B allocation unit (see runtime/sizeclasses_mcu.go)

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

// Init initializes the display using provided initialization commands. The
// initialization commands depends on the LCD pannel. The command that sets
// the display orientation and the color order must be the last one in the cmds
// See ili9341.GFX for working example.
func (d *Driver) Init(cmds []byte) {
	initialize(d.dci, cmds)
	d.rdir[0] = cmds[len(cmds)-1]
}

func (d *Driver) SetDir(dir int) image.Rectangle {
	if d.c.SetDir != nil {
		d.c.SetDir(d.dci, &d.rpf, &d.rdir, dir)
		if dir&1 != 0 {
			return image.Rectangle{Max: image.Pt(int(d.h), int(d.w))}
		}
	}
	return image.Rectangle{Max: image.Pt(int(d.w), int(d.h))}
}

func (d *Driver) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	if a>>15 == 0 {
		d.ctyp = ctrans // only 1-bit transparency is supported
		return
	}
	r >>= 8
	g >>= 8
	b >>= 8
	if d.pf&W16 != 0 {
		x := ((r ^ r>>5) | (b ^ b>>5)) & 7
		if d.pf&W24 == 0 {
			x &= 4
		} else {
			x |= (g ^ g>>6) & 3
		}
		if x == 0 {
			r &^= 7
			g &^= 3
			b &^= 7
			rgb565 := uint16(r<<8 | g<<3 | b>>3)
			d.csiz = 2
			if _, ok := d.dci.(WordNWriter); ok {
				d.ctyp = cfast
				d.cnpp = 1
				d.r = rgb565
				return
			}
			h := rgb565 >> 8
			l := rgb565 & 0xff
			if h == l {
				if _, ok := d.dci.(ByteNWriter); ok {
					d.ctyp = cfast
					d.cnpp = 2
					d.r = h
					return
				}
			}
			d.ctyp = cslow
			d.cnpp = 2
			d.r = uint16(h)
			d.g = uint16(l)
			return
		}
	}
	if d.pf&W24 == 0 {
		r &^= 3
		g &^= 3
		b &^= 3
	}
	if r == g && g == b {
		if _, ok := d.dci.(ByteNWriter); ok {
			d.ctyp = cfast
			d.csiz = 3
			d.cnpp = 3
			d.r = uint16(r)
			return
		}
	}
	d.ctyp = cslow
	d.csiz = 3
	d.cnpp = 3
	d.r = uint16(r)
	d.g = uint16(g)
	d.b = uint16(b)
}

func (d *Driver) Fill(r image.Rectangle) {
	if d.ctyp == ctrans {
		return
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	if d.c.SetPF != nil {
		d.c.SetPF(d.dci, &d.rpf, int(d.csiz))
	}
	d.c.StartWrite(d.dci, &d.xarg, r)
	n *= int(d.cnpp)
	if d.ctyp == cfast {
		if d.cnpp == 1 {
			d.dci.(WordNWriter).WriteWordN(d.r, n)
		} else {
			d.dci.(ByteNWriter).WriteByteN(byte(d.r), n)
		}
	} else {
		if d.ctyp == cslow {
			d.ctyp = cinbuf
			for i := 0; i < len(d.buf); i += int(d.csiz) {
				d.buf[i+0] = uint8(d.r)
				d.buf[i+1] = uint8(d.g)
				if d.csiz == 3 {
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
	d.dci.End()
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	dst := dst{r.Size(), 3}
	sip := imageAtPoint(src, sp)
	if op == draw.Src {
		if d.c.SetPF != nil {
			if mask == nil && sip.pixSize < dst.pixSize {
				dst.pixSize = sip.pixSize
			}
			d.c.SetPF(d.dci, &d.rpf, dst.pixSize)
		}
		d.c.StartWrite(d.dci, &d.xarg, r)
		bufUsed := drawSrc(d.dci, dst, src, sp, sip, mask, mp, d.buf[:])
		if bufUsed && d.ctyp == cinbuf {
			d.ctyp = cslow
		}
	} else {
		if d.c.SetPF != nil {
			d.c.SetPF(d.dci, &d.rpf, dst.pixSize)
		}
		buf := d.buf[:]
		if d.ctyp == cinbuf {
			d.ctyp = cslow
		}
		i := 0
		width := dst.size.X
		height := dst.size.Y
		for y := 0; y < height; y++ {
			j := y * sip.stride
			drawing := false
			for x := 0; x < width; x++ {
				ma := uint32(0x8000)
				if mask != nil {
					_, _, _, ma = mask.At(mp.X+x, mp.Y+y).RGBA()
				}
				if ma>>15 != 0 { // only 1-bit transparency supported
					var sr, sg, sb, sa uint32
					if sip.pixSize != 0 {
						sr, sg, sb, sa = fastRGBA(&sip, j)
						j += sip.pixSize
					} else {
						sr, sg, sb, sa = src.At(sp.X+x, sp.Y+y).RGBA()
					}
					if mask != nil {
						sa = (sa * ma / 0xffff) >> 15 // 1-bit transparency
						if sa != 0 {
							sr = sr * ma / 0xffff
							sg = sg * ma / 0xffff
							sb = sb * ma / 0xffff
						}
					}
					if sa != 0 {
						// opaque pixel
						if !drawing {
							drawing = true
							if i != 0 {
								d.dci.WriteBytes(buf[:i])
								i = 0
							}
							r1 := image.Rectangle{
								image.Pt(x, y),
								image.Pt(x+width, y+1),
							}.Add(r.Min)
							d.c.StartWrite(d.dci, &d.xarg, r1)
						}
						if dst.pixSize == 2 {
							sr >>= 11
							sg >>= 10
							sb >>= 11
							buf[i+0] = uint8(sr<<3 | sg>>3)
							buf[i+1] = uint8(sg<<5 | sb)
						} else {
							buf[i+0] = uint8(sr >> 8)
							buf[i+1] = uint8(sg >> 8)
							buf[i+2] = uint8(sb >> 8)
						}
						i += dst.pixSize
						if i == len(buf) {
							d.dci.WriteBytes(buf)
							i = 0
						}
						continue
					}
				}
				// transparent pixel
				drawing = false
			}
		}
		if i != 0 {
			d.dci.WriteBytes(buf[:i])
		}
	}
	d.dci.End()
}
