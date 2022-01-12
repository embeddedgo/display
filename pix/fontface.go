// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import "image"

// FontFace allows to convert an unicode codepoint (rune) to its graphical
// representation (glyph) in some font face. A font face represents a specific
// size, style and weight of a font.
//
// Multiple goroutines can use the same font face at the same time. If the
// implementation does not allow to share one font face instance by multiple
// goroutines it should provide a way to obtain multiple independent instances
// of it.
type FontFace interface {
	// Size returns the font height (interline spacing) and the ascent (height
	// above the baseline.
	Size() (height, ascent int)

	// Advance returns the glyph advance for the given rune. The advance
	// determines the x-distance on the baseline between the origin point of the
	// current character and the origin point of the next character.
	Advance(r rune) int

	// Glyph returns the graphical representation of the given rune in the alpha
	// channel of returned image. The image is valid until the next Glyph call.
	// The origin point is given in the img coordinates, can be outside of the
	// image bounds.
	Glyph(r rune) (img image.Image, origin image.Point, advance int)
}
