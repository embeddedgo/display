// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import "image"

// Font is an interface used by TextWriter to render a font.
type Font interface {
	// Bounds returns the string bounds. The returned rectangle is relative to
	// the  "base point" of the first glyph therefore some of its coordinates
	// can be negative.
	Bounds(s string) image.Rectangle

	// Glyph returns a glyph representing the given rune. The image bounds are
	// relative to the "base point". The returned image is valid until the next
	// Glyph call. The line terminators should be handled the same way as
	// unsupported runes.
	Glyph(r rune) image.Image
}
