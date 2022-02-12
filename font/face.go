// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package font

import "image"

// Face converts an unicode codepoint (rune) to its graphical representation
// (glyph) in some font face. A font face represents a specific size, style and
// weight of a font.
//
// The implementation of Face is allowed to be UNSAFE for concurrent use (see
// description of Glyph method). In such case a font repository should provide
// a way to request multiple independend Face instances that represents the same
// font face in the repository.
type Face interface {
	// Size returns the font face height (interline spacing) and the ascent
	// (height above the baseline).
	Size() (height, ascent int)

	// Advance returns the glyph advance for the given rune. The advance
	// determines the x-distance on the baseline between the origin point of the
	// current character and the origin point of next character.
	Advance(r rune) int

	// Glyph returns the graphical representation of the given rune in the alpha
	// channel of returned image. The image is valid until the next Glyph call
	// which makes Face UNSAFE for concurrent use. The origin point is given in
	// the img coordinates, may be outside of the image bounds.
	Glyph(r rune) (img image.Image, origin image.Point, advance int)
}
