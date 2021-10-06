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

type GRAM struct {
	DCI     tftdrv.DCI
	Size    image.Point
	PixSize int
}

func DrawSrc(dst GRAM, src image.Image, sp image.Point, sip Image, mask image.Image, mp image.Point, getBuf func() []byte, minChunk int) {
	var buf struct {
		p []byte
		i int
	}
	if sip.PixSize != 0 {
		// known type of source image so we can speed up access to pixel data
		width := dst.Size.X * sip.PixSize
		height := dst.Size.Y
		if sip.PixSize == dst.PixSize {
			// can write sip directly to the graphics RAM
			if len(sip.p) != 0 {
				if width == sip.stride {
					// write the entire sip
					dst.DCI.WriteBytes(sip.p[:height*sip.stride])
					print("P")
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
					print("p")
					return
				}
			} else if w, ok := dst.DCI.(tftdrv.StringWriter); ok {
				if width == sip.stride {
					// write the entire sip
					w.WriteString(sip.s[:height*sip.stride])
					print("S")
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
					print("s")
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
		print("b")
	}
	if buf.i != 0 {
		dst.DCI.WriteBytes(buf.p[:buf.i])
	}
}

/*
	var buf struct {
		p []byte
		i int
	}

		if mask == nil {
			if sip.pixSize != 0 {
				// known image type
				width := r.Dx() * sip.pixSize
				height := r.Dy()
				if sip.pixSize != 4 {
					// RGB or RGB16
					if len(sip.p) != 0 {
						if width == sip.stride {
							// write the entire src
							d.dci.WriteBytes(sip.p[:height*sip.stride])
							return
						}
						if width*4 > len(d.buf)*3 {
							// write line by line directly from src
							for {
								d.dci.WriteBytes(sip.p[:width])
								if height--; height == 0 {
									break
								}
								sip.p = sip.p[sip.stride:]
							}
							return
						}
					} else if w, ok := d.dci.(tftdrv.StringWriter); ok {
						if width == sip.stride {
							// write the entire src
							w.WriteString(sip.s[:height*sip.stride])
							return
						}
						if width*4 > len(d.buf)*3 {
							// write line by line directly from src
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
				buf.p = getBuf(d)
				j := 0
				k := width
				max := height * sip.stride
				dstPixSize := sip.pixSize
				if dstPixSize > 3 {
					dstPixSize = 3
				}
				for {
					if sip.p != nil {
						buf.p[buf.i+0] = sip.p[j+0]
						buf.p[buf.i+1] = sip.p[j+1]
						if dstPixSize == 3 {
							buf.p[buf.i+2] = sip.p[j+2]
						}
					} else {
						buf.p[buf.i+0] = sip.s[j+0]
						buf.p[buf.i+1] = sip.s[j+1]
						if dstPixSize == 3 {
							buf.p[buf.i+2] = sip.s[j+2]
						}
					}
					buf.i += dstPixSize
					j += sip.pixSize
					if buf.i == len(buf.p) {
						d.dci.WriteBytes(buf.p)
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
				// unknown image type
				buf.p = getBuf(d)
				r = r.Add(sp.Sub(r.Min))
				for y := r.Min.Y; y < r.Max.Y; y++ {
					for x := r.Min.X; x < r.Max.X; x++ {
						r, g, b, _ := src.At(x, y).RGBA()
						buf.p[buf.i+0] = uint8(r >> 8)
						buf.p[buf.i+1] = uint8(g >> 8)
						buf.p[buf.i+2] = uint8(b >> 8)
						buf.i += 3
						if buf.i == len(buf.p) {
							d.dci.WriteBytes(buf.p)
							buf.i = 0
						}
					}
				}
			}

	if buf.i != 0 {
		d.dci.WriteBytes(buf.p[:buf.i])
	}
	return

   text    data     bss     dec     hex filename
 756892    2740   20032  779664   be590 ili9341.elf

*/
