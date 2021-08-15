// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package font9

// Data implements subfont.Data interface storing characters as a single
// image with the glyphs placed side-by-side on a common baseline.
type Data struct {
	Info Info  // character descriptions
	Bits Image // image holding the glyphs, the baseline is at y = 0
}

// Info is the interface that wraps the Fontchar method.
//
// Fontchar descibes i-th character glyph in the subfont. The difference to the
// original Plan 9 fontchar is the lack of the top and bottom information.
type Info interface {
	Fontchar(i int) (x, left, advance int)
}

// Mono implements Info interface for monospace font.
type Mono uint8

// Prop implements Info interface for proportional font. Each character is
// described using 4 bytes (xh, xl, left, advance) where the x = uint16(xh)<<8 +
// uint16(xl) is the glyph position in the image, int8(left) is the offset
// to the glyph origin, uint8(advance) is the distance between two successive
// glyph origins when drawing.
type Prop []byte

// ImmProp is immutable version of Prop.
type ImmProp string
