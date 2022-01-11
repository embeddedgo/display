// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/draw"
	"unicode/utf8"
)

// Wrap determines what TextWriter does if the string does not fit in the area.
type Wrap uint8

const (
	NoWrap Wrap = iota
	WrapNewLine
	WrapSameLine
)

// StringWidth returns the number of horizontal pixels that would be occupied by
// the string if it were drawn using the given font face and NoWrap mode.
func StringWidth(s string, f FontFace) int {
	x := 0
	for _, r := range s {
		x += f.Advance(r)
	}
	return x
}

// Width works like StringWidth but for byte slices.
func Width(s []byte, f FontFace) int {
	x := 0
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRune(s[i:])
		x += f.Advance(r)
		i += size
	}
	return x
}

// TextWriter allows to write a text on the area. At least the Area, Face and
// Color fields must be set before use it.
//
// Notice that the Color field type is image.Image, not color.Color. This gives
// greater flexibility when drawing text. Set it to &image.Uniform{color} for
// traditional uniform color of glyphs.
//
// If Filter is not nil it is used to filter (scale, rotate, etc.) the glyphs
// obtained from Face.
type TextWriter struct {
	Area   *Area
	Face   FontFace
	Color  image.Image
	Pos    image.Point
	Filter func(glyph image.Image) image.Image
	Wrap   Wrap
	_      byte // literals must have keys to allow adding fields in the future
}

func (w *TextWriter) Write(s []byte) (int, error) {
	height, _ := w.Face.Size()
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRune(s[i:])
		drawRune(w, r, height)
		i += size
	}
	return len(s), nil
}

func (w *TextWriter) WriteString(s string) (int, error) {
	height, _ := w.Face.Size()
	for _, r := range s {
		drawRune(w, r, height)
	}
	return len(s), nil
}

func drawRune(w *TextWriter, r rune, height int) {
	mask, origin, advance := w.Face.Glyph(r)
	if mask == nil {
		mask, origin, advance = w.Face.Glyph(0)
		if mask == nil {
			return
		}
		if w.Filter != nil {
			mask = w.Filter(mask)
		}
	}
	nx := w.Pos.X + advance
	if w.Wrap != NoWrap && (nx > w.Area.bounds.Max.X || r == '\n') {
		if w.Wrap == WrapNewLine {
			w.Pos.Y += height
		}
		w.Pos.X = w.Area.bounds.Min.X
		if r == '\n' {
			return
		}
		nx = w.Pos.X + advance
	}
	mr := mask.Bounds()
	dr := mr.Add(w.Pos.Sub(origin))
	w.Area.Draw(dr, w.Color, image.Point{}, mask, mr.Min, draw.Over)
	w.Pos.X = nx
}
