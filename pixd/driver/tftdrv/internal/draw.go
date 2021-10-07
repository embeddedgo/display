// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"image"

	"github.com/embeddedgo/display/pixd"
	"github.com/embeddedgo/display/pixd/driver/tftdrv"
)

type Image struct {
	PixSize int

	stride int
	p      []byte
	s      string
}

func ImageAtPoint(img image.Image, pt image.Point) (out Image) {
	switch img := img.(type) {
	case *pixd.RGB16:
		out.PixSize = 2
		out.stride = img.Stride
		out.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.ImmRGB16:
		out.PixSize = 2
		out.stride = img.Stride
		out.s = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.RGB:
		out.PixSize = 3
		out.stride = img.Stride
		out.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.ImmRGB:
		out.PixSize = 3
		out.stride = img.Stride
		out.s = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *image.RGBA:
		out.PixSize = 4
		out.stride = img.Stride
		out.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	}
	return
}

// DDRAM provides access to the Display Data RAM (the nomenclature used by
// Philips / Epson / Ilitek) of the specified size (in pixels) and the pixel
// size (in bytes).
type DDRAM struct {
	DCI     tftdrv.DCI
	Size    image.Point
	PixSize int
}

// DrawSrc draws masked image to the prepared region of DDRAM. dst.PixSize must
// be 3 (RGB 888) or 2 (RGB 565).
func DrawSrc(dst DDRAM, src image.Image, sp image.Point, sip Image, mask image.Image, mp image.Point, getBuf func() []byte, minChunk int) {
	var buf struct {
		p []byte
		i int
	}
	width := dst.Size.X
	height := dst.Size.Y
	if mask == nil {
		if sip.PixSize != 0 {
			// known type of source image - we can speed up access to pixel data
			width *= sip.PixSize
			if sip.PixSize == dst.PixSize {
				// can write sip directly to the graphics RAM
				if len(sip.p) != 0 {
					if width == sip.stride {
						// write the entire sip
						dst.DCI.WriteBytes(sip.p[:height*sip.stride])
						return
					}
					if width >= minChunk {
						// write line by line directly from sip
						for {
							dst.DCI.WriteBytes(sip.p[:width])
							if height--; height == 0 {
								break
							}
							sip.p = sip.p[sip.stride:]
						}
						return
					}
				} else if w, ok := dst.DCI.(tftdrv.StringWriter); ok {
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
			buf.p = getBuf()
			j := 0
			k := width
			max := height * sip.stride
			for {
				if sip.p != nil {
					buf.p[buf.i+0] = sip.p[j+0]
					buf.p[buf.i+1] = sip.p[j+1]
					if dst.PixSize == 3 {
						buf.p[buf.i+2] = sip.p[j+2]
					}
				} else {
					buf.p[buf.i+0] = sip.s[j+0]
					buf.p[buf.i+1] = sip.s[j+1]
					if dst.PixSize == 3 {
						buf.p[buf.i+2] = sip.s[j+2]
					}
				}
				buf.i += dst.PixSize
				j += sip.PixSize
				if buf.i == len(buf.p) {
					dst.DCI.WriteBytes(buf.p)
					buf.i = 0
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
			buf.p = getBuf()
			r := image.Rectangle{sp, sp.Add(dst.Size)}
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					r, g, b, _ := src.At(x, y).RGBA()
					buf.p[buf.i+0] = uint8(r >> 8)
					buf.p[buf.i+1] = uint8(g >> 8)
					buf.p[buf.i+2] = uint8(b >> 8)
					buf.i += 3
					if buf.i == len(buf.p) {
						dst.DCI.WriteBytes(buf.p)
						buf.i = 0
					}
				}
			}
		}
	} else {
		// non-nil mask
		buf.p = getBuf()
		for y := 0; y < height; y++ {
			j := y * sip.stride
			for x := 0; x < width; x++ {
				var r, g, b uint32
				if sip.PixSize != 0 {
					if sip.p != nil {
						r = uint32(sip.p[j+0])
						g = uint32(sip.p[j+1])
						if sip.PixSize != 2 {
							b = uint32(sip.p[j+2])
						}
					} else {
						r = uint32(sip.s[j+0])
						g = uint32(sip.s[j+1])
						if sip.PixSize != 2 {
							b = uint32(sip.s[j+2])
						}
					}
					if sip.PixSize != 2 {
						r |= r << 8
						g |= g << 8
						b |= b << 8
					} else {
						r, g, b = r>>3, r&7<<3|g>>5, g&0x1F
						r = r<<11 | r<<6 | r<<1
						g = g<<10 | g<<4 | g>>2
						b = b<<11 | b<<6 | b<<1
					}
					j += sip.PixSize
				} else {
					r, g, b, _ = src.At(sp.X+x, sp.Y+y).RGBA()
				}
				_, _, _, ma := mask.At(mp.X+x, mp.Y+y).RGBA()
				r = r * ma / 0xffff
				g = g * ma / 0xffff
				b = b * ma / 0xffff
				if dst.PixSize != 2 {
					buf.p[buf.i+0] = uint8(r >> 8)
					buf.p[buf.i+1] = uint8(g >> 8)
					buf.p[buf.i+2] = uint8(b >> 8)
				} else {
					r >>= 11
					g >>= 10
					b >>= 11
					buf.p[buf.i+0] = uint8(r<<3 | g>>3)
					buf.p[buf.i+1] = uint8(g<<5 | b)
				}
				buf.i += dst.PixSize
				if buf.i == len(buf.p) {
					dst.DCI.WriteBytes(buf.p)
					buf.i = 0
				}
			}
		}
	}
	if buf.i != 0 {
		dst.DCI.WriteBytes(buf.p[:buf.i])
	}
}

func DrawOverNoRead(dst DDRAM, src image.Image, sp image.Point, sip Image, mask image.Image, mp image.Point, buf []byte, capaset func(r image.Rectangle)) {
	if sip.PixSize != 0 {

	}
}
