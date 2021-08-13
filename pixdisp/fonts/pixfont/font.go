// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixfont

import "image"

// Font is a pixmap based font.
type Font struct {
	Name     string
	Height   int16         // interline spacing (sumarizes all subfonts)
	Ascent   int16         // height above the baseline (sumarizes all subfonts)
	Subfonts []*Subfont    // ordered subfonts that make up the font
	Loader   SubfontLoader // used to load missing subfonts
}

func (f *Font) Size() (height, ascent int) {
	return int(f.Height), int(f.Ascent)
}

func getSubfont(f *Font, r rune) *Subfont {
	// TODO: binary search in ordered subfonts
	for _, sf := range f.Subfonts {
		if sf.First <= r && r <= sf.Last {
			return sf
		}
	}
	if f.Loader == nil {
		return nil
	}
	new := f.Loader.LoadSubfont(r)
	if new == nil {
		return nil
	}
	// TODO: binary search in ordered subfonts
	for i, sf := range f.Subfonts {
		if new.Last < sf.First {
			f.Subfonts = append(f.Subfonts[:i+1], f.Subfonts[i:]...)
			f.Subfonts[i] = new
			return new
		}
	}
	f.Subfonts = append(f.Subfonts, new)
	return new
}

func (f *Font) Advance(r rune) int {
	sf := getSubfont(f, r)
	if sf == nil {
		return 0
	}
	_, _, advance := sf.Info.Glyph(int(r - sf.First))
	return advance
}

func (f *Font) Glyph(r rune) (img image.Image, origin image.Point, advance int) {
	sf := getSubfont(f, r)
	if sf == nil {
		return
	}
	bounds, origin, advance := sf.Info.Glyph(int(r - sf.First))
	return sf.Bits.SubImage(bounds), origin, advance
}
