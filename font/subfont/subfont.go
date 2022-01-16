// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package subfont

import "image"

// Subfont provides Last-First+1 font glyphs for runes form First to Last. The
// glyphs are stored in Data starting from Offset.
type Subfont struct {
	First  rune // first character in the subfont
	Last   rune // last character in the subfont
	Offset int  // offset in Data to the first character
	Data   Data // character data
}

// Loader is the interface that wraps the Load method.
//
// Load loads the subfont containing a given rune. A successful call returns the
// pointer to the loaded subfont, otherwise the nil pointer is returned.
type Loader interface {
	Load(r rune) *Subfont
}

// Data represents a glyph storage.
type Data interface {
	// Advance returns the advance for the i-th glyph.
	Advance(i int) int

	// Glyph returns the data of the i-th glyph. The returned image is valid
	// until the next Glyph call.
	Glyph(i int) (img image.Image, origin image.Point, advance int)
}
