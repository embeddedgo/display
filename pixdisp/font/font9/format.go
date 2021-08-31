// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package font9

import "image"

// Fixed can hold data of a monospace font where all character subimages has
// the same width.
type Fixed struct {
	Left  int8
	Adv   uint8 // distance between two successive glyph origins when drawing
	Width uint8 // width of the glyph subimage
	Bits  Image // image holding the glyphs, the baseline is at y = 0
}

// Advance implements subfont.Data interface
func (d *Fixed) Advance(i int) int {
	return int(d.Adv)
}

// Glyph implements subfont.Data interface
func (d *Fixed) Glyph(i int) (img image.Image, origin image.Point, advance int) {
	r := d.Bits.Bounds()
	r.Min.X += int(d.Width)*i + int(d.Left)
	r.Max.X = r.Min.X + int(d.Width)
	img = d.Bits.SubImage(r)
	origin = image.Point{r.Min.X, 0}
	advance = int(d.Adv)
	return
}

// Num returns the number of characters covered by d.
func (d *Fixed) Num() int {
	return (d.Bits.Bounds().Dx() + int(d.Width) - 1) / int(d.Width)
}

// Variable can hold data of a monospace or proportional font. The character
// images can have different widths and origins.
type Variable struct {
	// Info stores information about N character subimagas in the N*4 + 2 bytes:
	// xlo0, xhi0, left0, advance0, ... , xloN, xhiN. The difference to the
	// original Plan 9 subfont is the lack of the top and bottom information.
	// Info is stored in a string to avoid being copied to RAM in the case of
	// systems that can leave read-only data in Flash. Use strings.Builder to
	// efficiently load/build it at runtime.
	Info string

	// Bits is an image holding the glyphs. The baseline is at y = 0.
	Bits Image
}

// Advance implements subfont.Data interface
func (d *Variable) Advance(i int) int {
	return int(d.Info[i*4+3])
}

// Glyph implements subfont.Data interface
func (d *Variable) Glyph(i int) (img image.Image, origin image.Point, advance int) {
	r := d.Bits.Bounds()
	info := d.Info[i*4:]
	r.Min.X = int(info[0]) | int(info[1])<<8
	r.Max.X = int(info[4]) | int(info[5])<<8
	img = d.Bits.SubImage(r)
	origin = image.Point{r.Min.X - int(int8(info[2])), 0}
	advance = int(info[3])
	return
}

// Num returns the number of characters covered by d.
func (d *Variable) Num() int {
	return len(d.Info) / 4
}
