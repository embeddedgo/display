// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"time"

	"github.com/embeddedgo/display/images"
)

// PF describes pixel data formats supported by display cpntroller.
type PF byte

const (
	R16 PF = 1 << iota // Read  2 bytes/pixel, RGB 565
	W16                // Write 2 bytes/pixel, RGB 565
	R24                // Read  3 bytes/pixel, RGB nnn,
	W24                // Write 3 bytes/pixel, RGB nnn,
	X2L                // Unused 2 low  bits per subpixel
	X2H                // Unused 2 high bits per subpixel

	R18L = R24 | X2H // 3 bytes/pixel, 6 bits/subpixel in low  part of byte
	W18L = W24 | X2H // 3 bytes/pixel, 6 bits/subpixel in low  part of byte
	R18H = R24 | X2L // 3 bytes/pixel, 6 bits/subpixel in high part of byte
	W18H = W24 | X2L // 3 bytes/pixel, 6 bits/subpixel in high part of byte
)

// Reg contains local copy of some controller registers to allow working with
// write-only displays.
type Reg struct {
	PF   [1]byte // pixel format relaed register
	Dir  [1]byte // direction/orientation related register
	Xarg [4]byte // scratch buffer to avoid allocation
}

// Ctrl contains display specific control functions.
type Ctrl struct {
	StartWrite func(dci DCI, reg *Reg, r image.Rectangle)
	Read       func(dci DCI, reg *Reg, r image.Rectangle, buf []byte)
	SetPF      func(dci DCI, reg *Reg, size int)
	SetDir     func(dci DCI, reg *Reg, dir int)
}

func initialize(dci DCI, reg *Reg, cmds []byte) {
	i := 0
	for i < len(cmds) {
		cmd := cmds[i]
		n := int(cmds[i+1])
		i += 2
		if n == 255 {
			time.Sleep(time.Duration(cmd) * time.Millisecond)
			continue
		}
		dci.Cmd(cmd)
		if n != 0 {
			k := i + n
			data := cmds[i:k]
			i = k
			dci.WriteBytes(data)
		}
	}
	dci.End()
	reg.Dir[0] = cmds[len(cmds)-1]
}

const (
	alphaTrans  = 0x0404
	alphaOpaque = 0xfc00
)

const (
	ctrans = iota
	cfast
	cslow
	cinbuf
)

type fillColor struct {
	r, g, b, a uint16
	typ        byte
	siz        int8
	npp        int8
	pf         PF
}

func setColor(c *fillColor, r, g, b, a uint32, dci DCI) {
	c.a = uint16(a)
	if a >= alphaOpaque {
		r >>= 8
		g >>= 8
		b >>= 8
		if c.pf&W16 != 0 {
			x := ((r ^ r>>5) | (b ^ b>>5)) & 7
			if c.pf&W24 == 0 {
				x &= 4
			} else {
				x |= (g ^ g>>6) & 3
			}
			if x == 0 {
				r &^= 7
				g &^= 3
				b &^= 7
				rgb565 := uint16(r<<8 | g<<3 | b>>3)
				c.siz = 2
				if _, ok := dci.(WordNWriter); ok {
					c.typ = cfast
					c.npp = 1
					c.r = rgb565
					return
				}
				h := rgb565 >> 8
				l := rgb565 & 0xff
				if h == l {
					if _, ok := dci.(ByteNWriter); ok {
						c.typ = cfast
						c.npp = 2
						c.r = h
						return
					}
				}
				if c.typ == cinbuf && c.r == h && c.g == l && c.npp == 2 {
					return // avoiding refilling the buffer with the same color
				}
				c.typ = cslow
				c.npp = 2
				c.r = uint16(h)
				c.g = uint16(l)
				return
			}
		}
		if c.pf&X2L != 0 {
			// only to increase probability of (r == g && g == b) below
			r &^= 3
			g &^= 3
			b &^= 3
		} else if c.pf&X2H != 0 {
			// required by display controller
			r >>= 2
			g >>= 2
			b >>= 2
		}
		if r == g && g == b {
			if _, ok := dci.(ByteNWriter); ok {
				c.typ = cfast
				c.siz = 3
				c.npp = 3
				c.r = uint16(r)
				return
			}
		}
		if c.typ == cinbuf && uint32(c.r) == r && uint32(c.g) == g &&
			uint32(c.b) == b && c.npp == 3 {
			return // avoiding refilling the buffer with the same color
		}
	}
	c.typ = cslow
	c.siz = 3
	c.npp = 3
	c.r = uint16(r)
	c.g = uint16(g)
	c.b = uint16(b)
}

func fillOpaque(dci DCI, c *fillColor, n int, buf []byte) {
	n *= int(c.npp)
	if c.typ == cfast {
		if c.npp == 1 {
			dci.(WordNWriter).WriteWordN(c.r, n)
		} else {
			dci.(ByteNWriter).WriteByteN(byte(c.r), n)
		}
	} else {
		if c.typ == cslow {
			c.typ = cinbuf
			for i := 0; i < len(buf); i += int(c.siz) {
				buf[i+0] = uint8(c.r)
				buf[i+1] = uint8(c.g)
				if c.siz == 3 {
					buf[i+2] = uint8(c.b)
				}
			}
		}
		m := len(buf)
		for {
			if m > n {
				m = n
			}
			dci.WriteBytes(buf[:m])
			n -= m
			if n == 0 {
				break
			}
		}
	}
}

type fastImage struct {
	stride  int
	pixSize int
	p       []byte
	s       string
}

func imageAtPoint(img image.Image, pt image.Point) (fi fastImage) {
	switch img := img.(type) {
	case *images.RGB16:
		fi.pixSize = 2
		fi.stride = img.Stride
		fi.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *images.ImmRGB16:
		fi.pixSize = 2
		fi.stride = img.Stride
		fi.s = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *images.RGB:
		fi.pixSize = 3
		fi.stride = img.Stride
		fi.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *images.ImmRGB:
		fi.pixSize = 3
		fi.stride = img.Stride
		fi.s = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *image.RGBA:
		fi.pixSize = 4
		fi.stride = img.Stride
		fi.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	}
	return
}

func fastRGBA(fi *fastImage, j int) (r, g, b, a uint32) {
	a = 255
	if fi.p != nil {
		r = uint32(fi.p[j+0])
		g = uint32(fi.p[j+1])
		if fi.pixSize != 2 {
			b = uint32(fi.p[j+2])
			if fi.pixSize == 4 {
				a = uint32(fi.p[j+3])
			}
		}
	} else {
		r = uint32(fi.s[j+0])
		g = uint32(fi.s[j+1])
		if fi.pixSize != 2 {
			b = uint32(fi.s[j+2])
			if fi.pixSize == 4 {
				a = uint32(fi.s[j+3])
			}
		}
	}
	a |= a << 8
	if fi.pixSize == 2 {
		r, g, b = r>>3, r&7<<3|g>>5, g&0x1f
		r = r<<11 | r<<6 | r<<1
		g = g<<10 | g<<4 | g>>2
		b = b<<11 | b<<6 | b<<1
	} else {
		r |= r << 8
		g |= g << 8
		b |= b << 8
	}
	return
}

type dst struct {
	size    image.Point
	pixSize int  // 2 for W16, 3 for W24, W18L, W18H
	shift   uint // noin-zero for W18L
}

// drawSrc draws masked image to the prepared region of frame memory.
func drawSrc(dci DCI, dst dst, src image.Image, sp image.Point, sip fastImage, mask image.Image, mp image.Point, buf []byte) (bufUsed bool) {
	i := 0
	width := dst.size.X
	height := dst.size.Y
	if mask == nil {
		if sip.pixSize != 0 {
			// known type of source image - we can speed up access to pixel data
			width *= sip.pixSize
			if dst.shift == 0 && dst.pixSize == sip.pixSize {
				// can write sip directly to the frame memory
				minChunk := len(buf) * 3 / 4
				if len(sip.p) != 0 {
					if width == sip.stride {
						// write the entire sip
						dci.WriteBytes(sip.p[:height*sip.stride])
						return
					}
					if width >= minChunk {
						// write line by line directly from sip
						for {
							dci.WriteBytes(sip.p[:width])
							if height--; height == 0 {
								break
							}
							sip.p = sip.p[sip.stride:]
						}
						return
					}
				} else if w, ok := dci.(StringWriter); ok {
					if width == sip.stride {
						// write the entire sip
						w.WriteString(sip.s[:height*sip.stride])
						return
					}
					if width > minChunk {
						// write line by line directly from sip
						for {
							w.WriteString(sip.s[:width])
							if height--; height == 0 {
								break
							}
							sip.s = sip.s[sip.stride:]
						}
						return
					}
				}
			}
			// buffered write
			bufUsed = true
			j := 0
			k := width
			max := height * sip.stride
			for {
				if sip.p != nil {
					buf[i+0] = sip.p[j+0] >> dst.shift
					buf[i+1] = sip.p[j+1] >> dst.shift
					if dst.pixSize == 3 {
						buf[i+2] = sip.p[j+2] >> dst.shift
					}
				} else {
					buf[i+0] = sip.s[j+0] >> dst.shift
					buf[i+1] = sip.s[j+1] >> dst.shift
					if dst.pixSize == 3 {
						buf[i+2] = sip.s[j+2] >> dst.shift
					}
				}
				i += dst.pixSize
				j += sip.pixSize
				if i == len(buf) {
					dci.WriteBytes(buf)
					i = 0
				}
				if j == k {
					k += sip.stride
					if k > max {
						break
					}
					j = k - width
				}
			}
		} else {
			// unknown type of source image - generic algorithm
			bufUsed = true
			r := image.Rectangle{sp, sp.Add(dst.size)}
			shift := 8 + dst.shift
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					r, g, b, _ := src.At(x, y).RGBA()
					buf[i+0] = uint8(r >> shift)
					buf[i+1] = uint8(g >> shift)
					buf[i+2] = uint8(b >> shift)
					i += 3
					if i == len(buf) {
						dci.WriteBytes(buf)
						i = 0
					}
				}
			}
		}
	} else {
		// non-nil mask
		bufUsed = true
		shift := 8 + dst.shift
		for y := 0; y < height; y++ {
			j := y * sip.stride
			for x := 0; x < width; x++ {
				var r, g, b uint32
				if sip.pixSize != 0 {
					r, g, b, _ = fastRGBA(&sip, j)
					j += sip.pixSize
				} else {
					r, g, b, _ = src.At(sp.X+x, sp.Y+y).RGBA()
				}
				_, _, _, ma := mask.At(mp.X+x, mp.Y+y).RGBA()
				r = r * ma / 0xffff
				g = g * ma / 0xffff
				b = b * ma / 0xffff
				if dst.pixSize != 2 {
					buf[i+0] = uint8(r >> shift)
					buf[i+1] = uint8(g >> shift)
					buf[i+2] = uint8(b >> shift)
				} else {
					r >>= 11
					g >>= 10
					b >>= 11
					buf[i+0] = uint8(r<<3 | g>>3)
					buf[i+1] = uint8(g<<5 | b)
				}
				i += dst.pixSize
				if i == len(buf) {
					dci.WriteBytes(buf)
					i = 0
				}
			}
		}
	}
	if i != 0 {
		dci.WriteBytes(buf[:i])
	}
	return
}
