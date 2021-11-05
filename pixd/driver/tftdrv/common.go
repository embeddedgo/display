// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"
	"time"

	"github.com/embeddedgo/display/pixd"
)

// PF describes supported pixel data formats
type PF byte

const (
	R16 PF = 1 << iota // Read  RGB 565, 2 bytes/pixel
	W16                // Write RGB 565, 2 bytes/pixel
	R18                // Read  RGB 666, 3 bytes/pixel
	W18                // Write RGB 666, 3 bytes/pixel
	R24                // Read  RGB 888, 3 bytes/pixel
	W24                // Write RGB 888, 3 bytes/pixel
)

type Ctrl struct {
	StartWrite func(dci DCI, xarg *[4]byte, r image.Rectangle)
	Read       func(dci RDCI, xarg *[4]byte, r image.Rectangle, buf []byte)
	SetPF      func(dci DCI, parg *[1]byte, size int)
	SetDir     func(dci DCI, parg *[1]byte, def byte, dir int)
}

const (
	transparent = 0

	osize = 0
	otype = 6 // Fill relies on the type field takes two MSbits

	// do not reorder
	fastByte = 0
	fastWord = 1
	bufInit  = 2 // getBuf relies on the one bit difference to the bufFull
	bufFull  = 3 // Fill relies on the both bits set
)

func initialize(dci DCI, cmds []byte) {
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
}

type fastImage struct {
	stride  int
	pixSize int
	p       []byte
	s       string
}

func imageAtPoint(img image.Image, pt image.Point) (fi fastImage) {
	switch img := img.(type) {
	case *pixd.RGB16:
		fi.pixSize = 2
		fi.stride = img.Stride
		fi.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.ImmRGB16:
		fi.pixSize = 2
		fi.stride = img.Stride
		fi.s = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.RGB:
		fi.pixSize = 3
		fi.stride = img.Stride
		fi.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.ImmRGB:
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
	pixSize int
}

// drawSrc draws masked image to the prepared region of GRAM. dst.PixSize must
// be 3 (RGB 888) or 2 (RGB 565).
func drawSrc(dci DCI, dst dst, src image.Image, sp image.Point, sip fastImage, mask image.Image, mp image.Point, getBuf func() []byte, minChunk int) {
	var buf []byte
	i := 0
	width := dst.size.X
	height := dst.size.Y
	if mask == nil {
		if sip.pixSize != 0 {
			// known type of source image - we can speed up access to pixel data
			width *= sip.pixSize
			if sip.pixSize == dst.pixSize {
				// can write sip directly to the graphics RAM
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
			buf = getBuf()
			j := 0
			k := width
			max := height * sip.stride
			for {
				if sip.p != nil {
					buf[i+0] = sip.p[j+0]
					buf[i+1] = sip.p[j+1]
					if dst.pixSize == 3 {
						buf[i+2] = sip.p[j+2]
					}
				} else {
					buf[i+0] = sip.s[j+0]
					buf[i+1] = sip.s[j+1]
					if dst.pixSize == 3 {
						buf[i+2] = sip.s[j+2]
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
			buf = getBuf()
			r := image.Rectangle{sp, sp.Add(dst.size)}
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					r, g, b, _ := src.At(x, y).RGBA()
					buf[i+0] = uint8(r >> 8)
					buf[i+1] = uint8(g >> 8)
					buf[i+2] = uint8(b >> 8)
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
		buf = getBuf()
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
					buf[i+0] = uint8(r >> 8)
					buf[i+1] = uint8(g >> 8)
					buf[i+2] = uint8(b >> 8)
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
}
