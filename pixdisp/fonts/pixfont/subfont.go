// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixfont

import "image"

// FontInfo is the interface that wraps the Glyph method. Glyph returns the
// bounds of i-th Glyph in the subfont image together with the origin point and
// advance.
type FontInfo interface {
	Glyph(i int) (bounds image.Rectangle, origin image.Point, advance int)
}

// Image is an image.Image with a SubImage method to obtain the portion of the
// image visible through r.
type Image interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

// Subfont consist of an image that contains N glyphs and metadata that
// describes how to get a subimage containing the glyph for a given rune.
type Subfont struct {
	First  rune     // first character in the subfont
	Last   rune     // last character in the subfont
	Height int16    // interline spacing
	Ascent int16    // height above the baseline
	Info   FontInfo // character descriptions
	Bits   Image    // image holding the glyphs
}

// SubfontLoader is the interface that wraps the LoadSubfont method.
// LoadSubfont loads the subfont containing a given rune. A successful call
// returns the pointer to the loaded subfont. Otherwise the nil pointer is
// returned.
type SubfontLoader interface {
	LoadSubfont(r rune) *Subfont
}
