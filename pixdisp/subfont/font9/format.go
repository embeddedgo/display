// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package font9

import "image"

// Mono can hold data of a monospace font.
type Mono struct {
	Left  int8  // offset to the origin
	Adv   uint8 // distance between two successive glyph origins when drawing
	Width uint8 // width of the glyph subimage
	Bits  Image // image holding the glyphs, the baseline is at y = 0
}

// Advance implements subfont.Data interface
func (d *Mono) Advance(i int) int {
	return int(d.Width)
}

// Glyph implements subfont.Data interface
func (d *Mono) Glyph(i int) (img image.Image, origin image.Point, advance int) {
	r := d.Bits.Bounds()
	r.Min.X += int(d.Width) * i
	r.Max.X = r.Min.X + int(d.Width)
	img = d.Bits.SubImage(r)
	origin = image.Point{r.Min.X + int(d.Left), 0}
	advance = int(d.Adv)
	return
}

// Prop can hold data of a proportional font. The difference to the original
// Plan 9 subfont is the lack of the top and bottom information in Info. Info is
// stored in a string to avoid copying it from Flash to RAM in case of embedded
// systems. Use strings.Builder to efficiently load/build it at runtime.
type Prop struct {
	Info string // character descriptions N x 4 bytes (advance, left, xl, xh)
	Bits Image  // image holding the glyphs, the baseline is at y = 0
}

// Advance implements subfont.Data interface
func (d *Prop) Advance(i int) int {
	return int(d.Info[i*4])
}

// Glyph implements subfont.Data interface
func (d *Prop) Glyph(i int) (img image.Image, origin image.Point, advance int) {
	r := d.Bits.Bounds()
	info := d.Info[i*4:]
	if len(info) >= 8 {
		r.Max.X = r.Min.X + (int(info[6]) | int(info[7])<<8)
	}
	r.Min.X += int(info[2]) | int(info[3])<<8
	img = d.Bits.SubImage(r)
	origin = image.Point{r.Min.X + int(int8(info[1])), 0}
	advance = int(info[0])
	return
}
