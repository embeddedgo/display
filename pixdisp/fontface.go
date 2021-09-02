// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import "image"

// FontFace is an interface used to render text. FontFace represents the
// specific size, style and weight of a font.
type FontFace interface {
	// Size returns the font height (interline spacing) and the ascent (height
	// above the baseline.
	Size() (height, ascent int)

	// Advance returns the glyph advance for the given rune. The advance
	// determines the x-distance on the baseline between the origin point of the
	// current character and the origin point of the next character.
	Advance(r rune) int

	// Glyph returns the data of the glyph for the given rune. The returned
	// image is valid until the next Glyph call. The origin point is given in
	// the img coordinates, can be (and usually is) outside of the glyph image.
	Glyph(r rune) (img image.Image, origin image.Point, advance int)
}
