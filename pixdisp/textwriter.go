// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

import (
	"image"
	"image/color"
	"image/draw"
	"unicode/utf8"
)

// WrapMode determines what the TextWriter does when the written text does not
// fit in the area.
type WrapMode uint8

const (
	WrapNewLine WrapMode = iota
	WrapSameLine
	NoWrap
)

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
	for len(s) > 0 {
		r, size := utf8.DecodeRune(s)
		drawRune(w, r)
		s = s[size:]
	}
	return 0, nil
}

func (w *TextWriter) WriteString(s string) (int, error) {
	for _, r := range s {
		drawRune(w, r)
	}
	return 0, nil
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
	if w.Wrap != NoWrap && (nx > w.Area.width || r == '\n') {
		if w.Wrap == WrapNewLine {
			h, _ := w.Face.Size()
			w.Pos.Y += h
		}
		w.Pos.X = 0
		nx = advance
		if r == '\n' {
			return
		}
	}
	img := &image.Uniform{w.Color}
	mr := mask.Bounds()
	dr := mr.Add(w.Pos.Sub(origin))
	w.Area.DrawMask(dr, img, image.Point{}, mask, mr.Min, draw.Over)
	w.Pos.X = nx
}
