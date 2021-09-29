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
	dci   tftdrv.DCI
	rgb16 int32
	pf    [1]byte
	xarg  [4]byte
	rgb   [25 * 3]uint8
	w, h  uint16
}

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

func (d *Driver) Init(swreset bool) {
	resetTime := time.Now()
	time.Sleep(5 * time.Millisecond)
	dci := d.dci
	for _, cmd := range initCmds {
		dci.Cmd(cmd[0])
		dci.WriteBytes(cmd[1:])
	}
	time.Sleep(resetTime.Add(120 * time.Millisecond).Sub(time.Now()))
	dci.Cmd(SLPOUT)
	time.Sleep(5 * time.Millisecond)
	dci.Cmd(DISPON)
}

func (d *Driver) SetMADCTL(madctl byte) {
	d.dci.Cmd(MADCTL)
	d.xarg[0] = madctl
	d.dci.WriteBytes(d.xarg[:1])
}

const (
	transparent = -1
	initRGB24   = -2
	fullRGB24   = -3
)

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
			d.rgb16 = transparent // only 1-bit transparency supported
			return
		}
		d.rgb[0] = cc.R
		d.rgb[1] = cc.G
		d.rgb[2] = cc.B
	default:
		r, g, b, a := c.RGBA()
		if a>>15 == 0 {
			d.rgb16 = transparent // only 1-bit transparency supported
			return
		}
		d.rgb[0] = uint8(r >> 8)
		d.rgb[1] = uint8(g >> 8)
		d.rgb[2] = uint8(b >> 8)
	}
	d.rgb16 = initRGB24
}

func pixset(d *Driver, pf byte) {
	if d.pf[0] != pf {
		d.pf[0] = pf
		d.dci.Cmd(PIXSET)
		d.dci.WriteBytes(d.pf[:])
	}
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

func (d *Driver) Fill(r image.Rectangle) {
	if d.rgb16 == transparent {
		return
	}
	n := r.Dx() * r.Dy()
	if n == 0 {
		return
	}
	capaset(d, r)
	pf := byte(PF18)
	if d.rgb16 >= 0 {
		pf = PF16
	}
	pixset(d, pf)
	d.dci.Cmd(RAMWR)
	if d.rgb16 >= 0 {
		d.dci.(tftdrv.WordNWriter).WriteWordN(uint16(d.rgb16), n)
		return
	}
	n *= 3
	if d.rgb[0] == d.rgb[1] && d.rgb[1] == d.rgb[2] {
		if w, ok := d.dci.(tftdrv.ByteNWriter); ok {
			w.WriteByteN(d.rgb[0], n)
			return
		}
	}
	if d.rgb16 == initRGB24 {
		d.rgb16 = fullRGB24
		for i := 3; i < len(d.rgb); i += 3 {
			d.rgb[i+0] = d.rgb[0]
			d.rgb[i+1] = d.rgb[1]
			d.rgb[i+2] = d.rgb[2]
		}
	}
	m := len(d.rgb)
	for {
		if m > n {
			m = n
		}
		d.dci.WriteBytes(d.rgb[:m])
		n -= m
		if n == 0 {
			break
		}
	}
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	var (
		bpix   []byte
		spix   string
		stride int
		ps     int
	)
	switch img := src.(type) {
	case *pix.ImageRGB16:
		bpix = img.Pix[img.PixOffset(sp.X, sp.Y):]
		stride = img.Stride
		ps = 2
	case *pix.ImmRGB16:
		spix = img.Pix[img.PixOffset(sp.X, sp.Y):]
		stride = img.Stride
		ps = 2
	case *pix.ImageRGB:
		bpix = img.Pix[img.PixOffset(sp.X, sp.Y):]
		stride = img.Stride
		ps = 3
	case *pix.ImmRGB:
		spix = img.Pix[img.PixOffset(sp.X, sp.Y):]
		stride = img.Stride
		ps = 3
	case *image.RGBA:
		bpix = img.Pix[img.PixOffset(sp.X, sp.Y):]
		stride = img.Stride
		ps = 4
	}
	buf := d.rgb[:]
	if d.rgb16 <= initRGB24 {
		d.rgb16 = initRGB24
		buf = d.rgb[3:]
	}
	i := 0
	if op == draw.Src {
		if mask == nil {
			capaset(d, r)
			pf := byte(PF18)
			if ps == 2 {
				pf = PF16
			}
			pixset(d, pf)
			d.dci.Cmd(RAMWR)
			if stride != 0 {
				width := r.Dx() * ps
				height := r.Dy()
				if ps != 4 {
					// RGB or RGB16
					if len(bpix) != 0 {
						if width == stride {
							// write the entire src
							d.dci.WriteBytes(bpix[:height*stride])
							return
						}
						if width*2 > len(buf) {
							// write line by line directly from src
							for {
								d.dci.WriteBytes(bpix[:width])
								if height--; height == 0 {
									break
								}
								bpix = bpix[stride:]
							}
							return
						}
					} else if w, ok := d.dci.(tftdrv.StringWriter); ok {
						if width == stride {
							// write the entire src
							w.WriteString(spix[:height*stride])
							return
						}
						if width*2 > len(buf) {
							// write line by line directly from src
							for {
								w.WriteString(spix[:width])
								if height--; height == 0 {
									break
								}
								spix = spix[stride:]
							}
							return
						}
					}
				}
				// buffered write
				j := 0
				k := width
				max := height * stride
				for {
					var r, g, b uint8
					if bpix != nil {
						r = bpix[j+0]
						g = bpix[j+1]
						b = bpix[j+2]
					} else {
						r = spix[j+0]
						g = spix[j+1]
						b = spix[j+2]
					}
					buf[i+0] = r
					buf[i+1] = g
					buf[i+2] = b
					i += 3
					j += ps
					if i == len(buf) {
						d.dci.WriteBytes(buf)
						i = 0
					}
					if j == k {
						k += stride
						if k > max {
							break
						}
						j = k - width
					}
				}
			} else {
				r = r.Add(sp.Sub(r.Min))
				for y := r.Min.Y; y < r.Max.Y; y++ {
					for x := r.Min.X; x < r.Max.X; x++ {
						r, g, b, _ := src.At(x, y).RGBA()
						buf[i+0] = uint8(r >> 8)
						buf[i+1] = uint8(g >> 8)
						buf[i+2] = uint8(b >> 8)
						i += 3
						if i == len(buf) {
							d.dci.WriteBytes(buf)
							i = 0
						}
					}
				}
			}
		}
	}
	if i != 0 {
		d.dci.WriteBytes(buf[:i])
	}
	return

}
