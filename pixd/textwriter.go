// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/color"
	"image/draw"
	"unicode/utf8"
)

// WrapMode determines what the TextWriter does if the string does not fit in
// the area.
type WrapMode uint8

const (
	WrapNewLine WrapMode = iota
	WrapSameLine
	NoWrap
)

// StringWidth returns the number of horizontal pixels that would be occupied by
// the string if it were drawn using the given font face and NoWrap mode.
func StringWidth(f FontFace, s string) int {
	x := 0
	for _, r := range s {
		x += f.Advance(r)
	}
	return x
}

// TextWriter allows to write a text on the area. The Area, Face and Color
// fields must be set before use it.
type TextWriter struct {
	Area  *Area
	Face  FontFace
	Color color.Color
	Pos   image.Point
	Wrap  WrapMode
}

func (w *TextWriter) Write(s []byte) (int, error) {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRune(s[i:])
		drawRune(w, r)
		i += size
	}
	return len(s), nil
}

func (w *TextWriter) WriteString(s string) (int, error) {
	for _, r := range s {
		drawRune(w, r)
	}
	return len(s), nil
}

func drawRune(w *TextWriter, r rune) {
	mask, origin, advance := w.Face.Glyph(r)
	if mask == nil {
		mask, origin, advance = w.Face.Glyph(0)
		if mask == nil {
			return
		}
	}
	nx := w.Pos.X + advance
	if w.Wrap != NoWrap && (nx > w.Area.bounds.Max.X || r == '\n') {
		if w.Wrap == WrapNewLine {
			h, _ := w.Face.Size()
			w.Pos.Y += h
		}
		w.Pos.X = w.Area.bounds.Min.X
		if r == '\n' {
			return
		}
		nx = w.Pos.X + advance
	}
	img := &image.Uniform{w.Color}
	mr := mask.Bounds()
	dr := mr.Add(w.Pos.Sub(origin))
	w.Area.Draw(dr, img, image.Point{}, mask, mr.Min, draw.Over)
	w.Pos.X = nx
}
