// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixfont

import "image"

// Font is a pixmap based font.
type Font struct {
	Name     string
	Height   int        // interline spacing
	Ascent   int        // height above the baseline
	Subfonts []*Subfont // ordered subfonts that make up the font
	Loader   FontLoader // used to load missing subfonts
}

// Subfont consist of an image that contains N glyphs and metadata that
// describes how to get a subimage containing the glyph for a given rune.
type Subfont struct {
	First  rune        // first character in the subfont
	Last   rune        // last character in the subfont
	Ascent int         // height above the baseline
	Info   FontInfo    // character descriptions
	Bits   image.Image // image holding the glyphs
}

// FontInfo is the interface that wraps the Glyph method. Glyph returns the i-th
// Glyph in the subfont.
type FontInfo interface {
	Glyph(i int) Glyph
}

type Glyph struct {
	X      int // x position in the image holding the glyphs
	Top    int // first non-zero scan line
	Bottom int // last non-zero scan line
	Left   int // offset of baseline
	Width  int // width of baseline
}

// FontLoader is the interface that wraps the Load method. Load loads the
// subfont containing a given rune. A successful call returns the pointer to
// the loaded subfont. Otherwise the nil pointer is returned.
type FontLoader interface {
	Load(r rune) *Subfont
}
