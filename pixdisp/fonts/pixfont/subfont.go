// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixfont

import "image"

type FontData interface {
	// Advance returns the advance for the i-th glyph.
	Advance(i int) int

	// Glyph returns the data of the i-th glyph. The returned image is valid
	// until the next Glyph call.
	Glyph(i int) (img image.Image, origin image.Point, advance int)
}

// Subfont consist of an image that contains N glyphs and metadata that
// describes how to get a subimage containing the glyph for a given rune.
type Subfont struct {
	First rune     // first character in the subfont
	Last  rune     // last character in the subfont
	Data  FontData // character data
}

// SubfontLoader is the interface that wraps the LoadSubfont method.
// LoadSubfont loads the subfont containing a given rune. A successful call
// returns the pointer to the loaded subfont, otherwise the nil pointer is
// returned.
type SubfontLoader interface {
	LoadSubfont(r rune) *Subfont
}
