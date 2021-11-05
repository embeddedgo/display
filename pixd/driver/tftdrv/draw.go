// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdrv

import (
	"image"

	"github.com/embeddedgo/display/pixd"
)

type fastImage struct {
	stride  int
	pixSize int
	p       []byte
	s       string
}

type dst struct {
	size    image.Point
	pixSize int
}

func imageAtPoint(img image.Image, pt image.Point) (out fastImage) {
	switch img := img.(type) {
	case *pixd.RGB16:
		out.pixSize = 2
		out.stride = img.Stride
		out.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.ImmRGB16:
		out.pixSize = 2
		out.stride = img.Stride
		out.s = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.RGB:
		out.pixSize = 3
		out.stride = img.Stride
		out.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *pixd.ImmRGB:
		out.pixSize = 3
		out.stride = img.Stride
		out.s = img.Pix[img.PixOffset(pt.X, pt.Y):]
	case *image.RGBA:
		out.pixSize = 4
		out.stride = img.Stride
		out.p = img.Pix[img.PixOffset(pt.X, pt.Y):]
	}
	return
}

// drawSrc draws masked image to the prepared region of GRAM. dst.PixSize must
// be 3 (RGB 888) or 2 (RGB 565).
func drawSrc(dci DCI, dst dst, src image.Image, sp image.Point, sip fastImage, mask image.Image, mp image.Point, getBuf func() []byte, minChunk int) {
	var buf struct {
		p []byte
		i int
	}
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
			buf.p = getBuf()
			j := 0
			k := width
			max := height * sip.stride
			for {
				if sip.p != nil {
					buf.p[buf.i+0] = sip.p[j+0]
					buf.p[buf.i+1] = sip.p[j+1]
					if dst.pixSize == 3 {
						buf.p[buf.i+2] = sip.p[j+2]
					}
				} else {
					buf.p[buf.i+0] = sip.s[j+0]
					buf.p[buf.i+1] = sip.s[j+1]
					if dst.pixSize == 3 {
						buf.p[buf.i+2] = sip.s[j+2]
					}
				}
				buf.i += dst.pixSize
				j += sip.pixSize
				if buf.i == len(buf.p) {
					dci.WriteBytes(buf.p)
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
			r := image.Rectangle{sp, sp.Add(dst.size)}
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					r, g, b, _ := src.At(x, y).RGBA()
					buf.p[buf.i+0] = uint8(r >> 8)
					buf.p[buf.i+1] = uint8(g >> 8)
					buf.p[buf.i+2] = uint8(b >> 8)
					buf.i += 3
					if buf.i == len(buf.p) {
						dci.WriteBytes(buf.p)
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
				if sip.pixSize != 0 {
					if sip.p != nil {
						r = uint32(sip.p[j+0])
						g = uint32(sip.p[j+1])
						if sip.pixSize != 2 {
							b = uint32(sip.p[j+2])
						}
					} else {
						r = uint32(sip.s[j+0])
						g = uint32(sip.s[j+1])
						if sip.pixSize != 2 {
							b = uint32(sip.s[j+2])
						}
					}
					if sip.pixSize == 2 {
						r, g, b = r>>3, r&7<<3|g>>5, g&0x1F
						r = r<<11 | r<<6 | r<<1
						g = g<<10 | g<<4 | g>>2
						b = b<<11 | b<<6 | b<<1
					} else {
						r |= r << 8
						g |= g << 8
						b |= b << 8
					}
					j += sip.pixSize
				} else {
					r, g, b, _ = src.At(sp.X+x, sp.Y+y).RGBA()
				}
				_, _, _, ma := mask.At(mp.X+x, mp.Y+y).RGBA()
				r = r * ma / 0xffff
				g = g * ma / 0xffff
				b = b * ma / 0xffff
				if dst.pixSize != 2 {
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
				buf.i += dst.pixSize
				if buf.i == len(buf.p) {
					dci.WriteBytes(buf.p)
					buf.i = 0
				}
			}
		}
	}
	if buf.i != 0 {
		dci.WriteBytes(buf.p[:buf.i])
	}
}

// drawOverNoRead draws masked image. It cannot read the content of frame memory
// so it reduces the alpha channel to one bit and draws only opaque parts of
// the masked image. dst.PixSize must be 3 (RGB 888) or 2 (RGB 565).
func drawOverNoRead(dci DCI, dst dst, dmin image.Point, src image.Image, sp image.Point, sip fastImage, mask image.Image, mp image.Point, buffer []byte, startWrite func(r image.Rectangle)) {
	var buf struct {
		p []byte
		i int
	}
	buf.p = buffer
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
				var r, g, b, a uint32
				if sip.pixSize != 0 {
					a = 0xff
					if sip.p != nil {
						r = uint32(sip.p[j+0])
						g = uint32(sip.p[j+1])
						if sip.pixSize != 2 {
							b = uint32(sip.p[j+2])
							if sip.pixSize == 4 {
								a = uint32(sip.p[j+3])
							}
						}
					} else {
						r = uint32(sip.s[j+0])
						g = uint32(sip.s[j+1])
						if sip.pixSize != 2 {
							b = uint32(sip.s[j+2])
							if sip.pixSize == 4 {
								a = uint32(sip.s[j+3])
							}
						}
					}
					a |= a << 8
					if sip.pixSize == 2 {
						r, g, b = r>>3, r&7<<3|g>>5, g&0x1f
						r = r<<11 | r<<6 | r<<1
						g = g<<10 | g<<4 | g>>2
						b = b<<11 | b<<6 | b<<1
					} else {
						r |= r << 8
						g |= g << 8
						b |= b << 8
					}
					j += sip.pixSize
				} else {
					r, g, b, a = src.At(sp.X+x, sp.Y+y).RGBA()
				}
				if mask != nil {
					a = (a * ma / 0xffff) >> 15 // we are interested in MSbit
					if a != 0 {
						r = r * ma / 0xffff
						g = g * ma / 0xffff
						b = b * ma / 0xffff
					}
				}
				if a != 0 {
					// opaque pixel
					if !drawing {
						drawing = true
						if buf.i != 0 {
							dci.WriteBytes(buf.p[:buf.i])
							buf.i = 0
						}
						r := image.Rectangle{
							image.Pt(x, y),
							image.Pt(x+width, y+1),
						}.Add(dmin)
						startWrite(r)
					}
					if dst.pixSize == 2 {
						r >>= 11
						g >>= 10
						b >>= 11
						buf.p[buf.i+0] = uint8(r<<3 | g>>3)
						buf.p[buf.i+1] = uint8(g<<5 | b)
					} else {
						buf.p[buf.i+0] = uint8(r >> 8)
						buf.p[buf.i+1] = uint8(g >> 8)
						buf.p[buf.i+2] = uint8(b >> 8)
					}
					buf.i += dst.pixSize
					if buf.i == len(buf.p) {
						dci.WriteBytes(buf.p)
						buf.i = 0
					}
					continue
				}
			}
			// transparent pixel
			drawing = false
		}
	}
	if buf.i != 0 {
		dci.WriteBytes(buf.p[:buf.i])
	}
}
