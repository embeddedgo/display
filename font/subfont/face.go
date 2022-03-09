// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package subfont

import "image"

// Face is a collection of subfonts from one font with the same size, style and
// weight.
type Face struct {
	Height   int16      // interline spacing (sumarizes all subfonts)
	Ascent   int16      // height above the baseline (sumarizes all subfonts)
	Subfonts []*Subfont // ordered subfonts that make up the face
	Loader   Loader     // used to load missing subfonts
}

// Size implements font.Face interface.
func (f *Face) Size() (height, ascent int) {
	return int(f.Height), int(f.Ascent)
}

// Advance implements font.Face interface.
func (f *Face) Advance(r rune) int {
	sf := getSubfont(f, r)
	if sf == nil {
		return 0
	}
	return sf.Data.Advance(int(r - sf.First))
}

// Glyph implements font.Face interface.
func (f *Face) Glyph(r rune) (img image.Image, origin image.Point, advance int) {
	sf := getSubfont(f, r)
	if sf == nil {
		return
	}
	return sf.Data.Glyph(int(r - sf.First))
}

func getSubfont(f *Face, r rune) (sf *Subfont) {
	// TODO: binary search
	for _, sf = range f.Subfonts {
		if sf != nil && sf.First <= r && r <= sf.Last {
			return sf
		}
	}
	if f.Loader != nil {
		sf, f.Subfonts = f.Loader.Load(r, f.Subfonts)
	}
	return sf
}
