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
	xarg  [4]byte
	pf    [1]byte
	w, h  uint16
	cinfo byte
	cfast uint16
	buf   [24 * 3]uint8 // must be multiple of two and three
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
	transparent = 0
	fastByte    = 0 << 6
	fastWord    = 1 << 6
	bufInit     = 2 << 6
	bufFull     = 3 << 6
)

func (d *Driver) SetColor(c color.Color) {
	var r, g, b uint32
	switch cc := c.(type) {
	case color.RGBA:
		if cc.A>>7 == 0 {
			d.cinfo = transparent // only 1-bit transparency supported
			return
		}
		r = cc.R
		g = cc.G
		b = cc.B
	default:
		var a uint32
		r, g, b, a = c.RGBA()
		if a>>15 == 0 {
			d.cinfo = transparent // only 1-bit transparency supported
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
	if r&4 == 0 && b&4 == 0 {
		rgb565 := r<<8 | g<<3 | b>>3
		h := rgb565 >> 8
		l := rgb565 && 0xff
		if h == l {
			if _, ok := d.dci.(tftdrv.ByteNWriter); ok {
				d.cinfo = fastByte | 2<<4 | MCU16
				d.cfast = uint16(h)
				return
			}
		}
		if _, ok := d.dci.(tftdrv.WordNWriter); ok {
			d.cinfo = fastWord | 2<<4 | MCU16
			d.cfast = uint16(rgb565)
			return
		}
		d.cinfo = bufInit | 2<<4 | MCU16
		d.buf[0] = byte(h)
		d.buf[1] = byte(l)
		return
	}
	if r == g && g == b {
		if _, ok := d.dci.(tftdrv.ByteNWriter); ok {
			d.cfast = uint16(r)
			d.cinfo = fastByte | 3<<4 | MCU18
			return
		}
	}
	d.cinfo = bufInit | 3<<4 | MCU18
	d.buf[0] = uint8(r)
	d.buf[1] = uint8(g)
	d.buf[2] = uint8(b)
}

func pixset(d *Driver) {
	pf = d.cinfo & 0xf
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
	pixset(d)
	d.dci.Cmd(RAMWR)

	if d.rgb16 >= 0 {
		d.dci.(tftdrv.WordNWriter).WriteWordN(uint16(d.rgb16), n)
		return
	}
	if d.rgb16 >= fullRGB16 {
		n *= 2
		if d.buf[0] == d.buf[1] {
			if w, ok := d.dci.(tftdrv.ByteNWriter); ok {
				w.WriteByteN(d.buf[0], n)
				return
			}
		}
		if d.rgb16 == initRGB16 {
			d.rgb16 = fullRGB16
			for i := 2; i < len(d.rgb); i += 2 {
				d.buf[i+0] = d.buf[0]
				d.buf[i+1] = d.buf[1]
			}
		}
	} else {
		n *= 3
		if d.buf[0] == d.buf[1] && d.buf[1] == d.buf[2] {
			if w, ok := d.dci.(tftdrv.ByteNWriter); ok {
				w.WriteByteN(d.buf[0], n)
				return
			}
		}
		if d.rgb16 == initRGB24 {
			d.rgb16 = fullRGB24
			for i := 3; i < len(d.rgb); i += 3 {
				d.buf[i+0] = d.buf[0]
				d.buf[i+1] = d.buf[1]
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
	var (
		srcBytes   []byte
		srcString  string
		srcStride  int
		srcPixSize int
	)
	switch img := src.(type) {
	case *pix.ImageRGB16:
		srcBytes = img.Pix[img.PixOffset(sp.X, sp.Y):]
		srcStride = img.Stride
		srcPixSize = 2
	case *pix.ImmRGB16:
		srcString = img.Pix[img.PixOffset(sp.X, sp.Y):]
		srcStride = img.Stride
		srcPixSize = 2
	case *pix.ImageRGB:
		srcBytes = img.Pix[img.PixOffset(sp.X, sp.Y):]
		srcStride = img.Stride
		srcPixSize = 3
	case *pix.ImmRGB:
		srcString = img.Pix[img.PixOffset(sp.X, sp.Y):]
		srcStride = img.Stride
		srcPixSize = 3
	case *image.RGBA:
		srcBytes = img.Pix[img.PixOffset(sp.X, sp.Y):]
		srcStride = img.Stride
		srcPixSize = 4
	}
	buf := d.buf[:]
	if d.rgb16 <= initRGB24 {
		if d.rgb16 == initRGB24 {
			d.rgb16 = initRGB24
		}
		buf = d.buf[3:]
	} else if d.rgb16 <= initRGB16 {
		if d.rgb16 == initRGB16 {
			d.rgb16 = initRGB16
		}
		buf = d.buf[2:]
	}
	i := 0
	if op == draw.Src {
		capaset(d, r)
		pf := byte(PF18)
		if srcPixSize == 2 {
			pf = PF16
		}
		pixset(d, pf)
		d.dci.Cmd(RAMWR)
		if mask == nil {
			if srcPixSize != 0 {
				// known image type
				width := r.Dx() * srcPixSize
				height := r.Dy()
				if srcPixSize != 4 {
					// RGB or RGB16
					if len(srcBytes) != 0 {
						if width == srcStride {
							// write the entire src
							d.dci.WriteBytes(srcBytes[:height*srcStride])
							return
						}
						if width*4 > len(buf)*3 {
							// write line by line directly from src
							for {
								d.dci.WriteBytes(srcBytes[:width])
								if height--; height == 0 {
									break
								}
								srcBytes = srcBytes[srcStride:]
							}
							return
						}
					} else if w, ok := d.dci.(tftdrv.StringWriter); ok {
						if width == srcStride {
							// write the entire src
							w.WriteString(srcString[:height*srcStride])
							return
						}
						if width*4 > len(buf)*3 {
							// write line by line directly from src
							for {
								w.WriteString(srcString[:width])
								if height--; height == 0 {
									break
								}
								srcString = srcString[srcStride:]
							}
							return
						}
					}
				}
				// buffered write
				j := 0
				k := width
				max := height * srcStride
				for {
					var r, g, b uint8
					if srcBytes != nil {
						r = srcBytes[j+0]
						g = srcBytes[j+1]
						b = srcBytes[j+2]
					} else {
						r = srcString[j+0]
						g = srcString[j+1]
						b = srcString[j+2]
					}
					buf[i+0] = r
					buf[i+1] = g
					buf[i+2] = b
					i += 3
					j += srcPixSize
					if i == len(buf) {
						d.dci.WriteBytes(buf)
						i = 0
					}
					if j == k {
						k += srcStride
						if k > max {
							break
						}
						j = k - width
					}
				}
			} else {
				// unknown image type
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
		} else {
			// mask != nil
			var (
				maskBytes   []byte
				maskString  string
				maskShift   uint
				maskStride  int
				maskPixBitN uint
			)
			switch img := mask.(type) {
			case *pix.ImageAlphaN:
				offset, shift := img.PixOffset(sp.X, sp.Y)
				maskBytes = img.Pix[offset:]
				maskShift = shift
				maskStride = img.Stride
				maskPixBitN = 1 << img.LogN
			case *pix.ImmAlphaN:
				offset, shift := img.PixOffset(sp.X, sp.Y)
				maskString = img.Pix[offset:]
				maskShift = shift
				maskStride = img.Stride
				maskPixBitN = 1 << img.LogN
			case *image.Alpha:
				maskBytes = img.Pix[img.PixOffset(sp.X, sp.Y):]
				maskStride = img.Stride
				maskPixBitN = 8
			}

			_ = maskBytes
			_ = maskString
			_ = maskShift
			_ = maskStride
			_ = maskPixBitN
		}
	}
	if i != 0 {
		d.dci.WriteBytes(buf[:i])
	}
	return

}
