// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"image/color"
	"image/draw"
)

// Driver implements pix.Driver interface with a limited support for draw.Over
// operation. It is designed for write-only displays (doesn't use DCI.ReadBytes
// method) so the alpha blending is slow and reduced to 1-bit resolution. Use
// DriverOver if the full-fledged Porter-Duff composition is required and the
// display supports reading from its frame memory.
type Driver struct {
	dci   DCI
	ctrl  *Ctrl
	w, h  uint16
	color fillColor
	reg   Reg
	buf   [74 * 3]byte // must be multiple of two and three
} // ont 32-bit MCU the size of this struct is 256 B, a full 256 B allocation unit (see runtime/sizeclasses_mcu.go)

// New returns new Driver.
func New(dci DCI, w, h uint16, pf PF, ctrl *Ctrl) *Driver {
	d := new(Driver)
	d.dci = dci
	d.ctrl = ctrl
	d.w = w
	d.h = h
	d.color.pf = pf
	return d
}

func (d *Driver) Err(clear bool) error { return d.dci.Err(clear) }
func (d *Driver) Flush()               {}

// Init initializes the display using provided initialization commands. The
// initialization commands depends on the LCD pannel. The command that sets
// the display orientation and the color order must be the last one in the cmds
// See ili9341.GFX for working example.
func (d *Driver) Init(cmds []byte) {
	initialize(d.dci, &d.reg, cmds)
}

func (d *Driver) SetDir(dir int) image.Rectangle {
	if d.ctrl.SetDir != nil {
		d.ctrl.SetDir(d.dci, &d.reg, dir)
		if dir&1 != 0 {
			return image.Rectangle{Max: image.Pt(int(d.h), int(d.w))}
		}
	}
	return image.Rectangle{Max: image.Pt(int(d.w), int(d.h))}
}

func (d *Driver) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	if a>>15 == 0 {
		d.color.typ = ctrans // only 1-bit transparency is supported
		return
	}
	setColor(&d.color, r, g, b, alphaOpaque, d.dci)
}

func (d *Driver) Fill(r image.Rectangle) {
	if d.color.typ == ctrans {
		return
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	if d.ctrl.SetPF != nil {
		d.ctrl.SetPF(d.dci, &d.reg, int(d.color.siz))
	}
	d.ctrl.StartWrite(d.dci, &d.reg, r)
	fillOpaque(d.dci, &d.color, n, d.buf[:])
	d.dci.End()
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	dst := dst{size: r.Size(), pixSize: 3}
	if d.color.pf&X2H != 0 {
		dst.shift = 2
	}
	sip := imageAtPoint(src, sp)
	if op == draw.Src {
		if d.ctrl.SetPF != nil {
			if mask == nil && sip.pixSize != 0 && sip.pixSize < dst.pixSize {
				dst.pixSize = sip.pixSize
			}
			d.ctrl.SetPF(d.dci, &d.reg, dst.pixSize)
		}
		d.ctrl.StartWrite(d.dci, &d.reg, r)
		bufUsed := drawSrc(d.dci, dst, src, sp, sip, mask, mp, d.buf[:])
		if bufUsed && d.color.typ == cinbuf {
			d.color.typ = cslow
		}
	} else {
		if d.ctrl.SetPF != nil {
			d.ctrl.SetPF(d.dci, &d.reg, dst.pixSize)
		}
		buf := d.buf[:]
		if d.color.typ == cinbuf {
			d.color.typ = cslow
		}
		i := 0
		width := dst.size.X
		height := dst.size.Y
		shift := 8 + dst.shift
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
							r1 := r
							r1.Min.X += x
							r1.Min.Y += y
							d.ctrl.StartWrite(d.dci, &d.reg, r1)
						}
						if dst.pixSize == 2 {
							sr >>= 11
							sg >>= 10
							sb >>= 11
							buf[i+0] = uint8(sr<<3 | sg>>3)
							buf[i+1] = uint8(sg<<5 | sb)
						} else {
							buf[i+0] = uint8(sr >> shift)
							buf[i+1] = uint8(sg >> shift)
							buf[i+2] = uint8(sb >> shift)
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
