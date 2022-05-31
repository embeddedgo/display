// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/color"
	"image/draw"
	"unicode"
	"unicode/utf8"

	"github.com/embeddedgo/display/font"
)

// Wrapping modes determine what TextWriter does if the string does not
// fit in the drawing area.
const (
	NoWrap byte = iota
	WrapNewLine
	WrapSameLine
)

// Line breaking mode.
const (
	BreakAny   byte = iota // break at any rune
	BreakSpace             // break at unicode White Space
)

// StringWidth returns the number of horizontal pixels that would be occupied by
// the string if it were drawn using the given font face and NoWrap mode.
func StringWidth(s string, f font.Face) int {
	x := 0
	for _, r := range s {
		x += f.Advance(r)
	}
	return x
}

// Width works like StringWidth but for byte slices.
func Width(s []byte, f font.Face) int {
	x := 0
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRune(s[i:])
		x += f.Advance(r)
		i += size
	}
	return x
}

// TextWriter allows to write a text on an area. At least the Area, Face and
// Color fields must be set before use it.
//
// Notice that the Color field type is image.Image, not color.Color. Set it to
// &image.Uniform{color} for traditional uniform color of glyphs.
type TextWriter struct {
	Area   *Area       // drawing area
	Face   font.Face   // source of glyphs
	Color  image.Image // glyph color (foreground image)
	Pos    image.Point // position of the next glyph
	Offset image.Point // offset from the Pos to the glyph origin
	Wrap   byte        // wrapping mode
	Break  byte        // line breaking mode
	_      byte        // require keys in literals
}

// SetColor provides a convenient way to set drawing color. If the w.Color
// contains value of type *image.Uniform it modifies the color of the current
// image. Otherwise it sets w.Color to &image.Uniform{c}.
func (w *TextWriter) SetColor(c color.Color) {
	switch p := c.(type) {
	case color.RGBA:
		if (int(p.A)-int(p.R))|(int(p.A)-int(p.G))|(int(p.A)-int(p.B)) < 0 {
			panic(badAlphaPremul)
		}
	case color.RGBA64:
		if (int(p.A)-int(p.R))|(int(p.A)-int(p.G))|(int(p.A)-int(p.B)) < 0 {
			panic(badAlphaPremul)
		}
	}
	if img, ok := w.Color.(*image.Uniform); ok {
		img.C = c
	} else {
		w.Color = image.NewUniform(c)
	}
}

// SetFace provides a convenient way to set w.Face and at the same time modify
// w.Offset in such a way that the current baseline is unaffected.
func (w *TextWriter) SetFace(f font.Face) {
	_, ascent := f.Size()
	if w.Face != nil {
		_, a := w.Face.Size()
		ascent -= a
	}
	w.Face = f
	w.Offset.Y += ascent
}

func writeNoWrap(w *TextWriter, s []byte) {
	bounds := w.Area.Bounds()
	for i := 0; i < len(s) && w.Pos.X < bounds.Max.X; {
		r, m := utf8.DecodeRune(s[i:])
		i += m
		if r == utf8.RuneError {
			continue
		}
		drawRune(w, r)
	}
}

func (w *TextWriter) Write(s []byte) (int, error) {
	if w.Wrap == NoWrap {
		writeNoWrap(w, s)
		return len(s), nil
	}
	bounds := w.Area.Bounds()
	height, _ := w.Face.Size()
	posX := w.Pos.X
	lastSpace := -1
	done := 0
	for i, m := 0, 0; i < len(s); i += m {
		var r rune
		r, m = utf8.DecodeRune(s[i:])
		if r == utf8.RuneError {
			continue
		}
		var advance, skip int
		if r == '\n' {
			skip = m
			goto wrap
		}
		if w.Break == BreakSpace && unicode.IsSpace(r) {
			lastSpace = i
		}
		advance = w.Face.Advance(r)
		posX += advance
		if posX > bounds.Max.X {
			if w.Break == BreakSpace && lastSpace >= 0 {
				i = lastSpace
				lastSpace = -1
				_, m = utf8.DecodeRune(s[i:])
				advance = 0
				skip = m
			}
			goto wrap
		}
		continue
	wrap:
		writeNoWrap(w, s[done:i])
		done = i + skip
		w.Pos.X = bounds.Min.X
		posX = w.Pos.X + advance
		if w.Wrap == WrapNewLine {
			w.Pos.Y += height
		}
	}
	writeNoWrap(w, s[done:])
	return len(s), nil
}

func writeStringNoWrap(w *TextWriter, s string) {
	bounds := w.Area.Bounds()
	for i := 0; i < len(s) && w.Pos.X < bounds.Max.X; {
		r, m := utf8.DecodeRuneInString(s[i:])
		i += m
		if r == utf8.RuneError {
			continue
		}
		drawRune(w, r)
	}
}

func (w *TextWriter) WriteString(s string) (int, error) {
	if w.Wrap == NoWrap {
		writeStringNoWrap(w, s)
		return len(s), nil
	}
	bounds := w.Area.Bounds()
	height, _ := w.Face.Size()
	posX := w.Pos.X
	lastSpace := -1
	done := 0
	for i, m := 0, 0; i < len(s); i += m {
		var r rune
		r, m = utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError {
			continue
		}
		var advance, skip int
		if r == '\n' {
			skip = m
			goto wrap
		}
		if w.Break == BreakSpace && unicode.IsSpace(r) {
			lastSpace = i
		}
		advance = w.Face.Advance(r)
		posX += advance
		if posX > bounds.Max.X {
			if w.Break == BreakSpace && lastSpace >= 0 {
				i = lastSpace
				lastSpace = -1
				_, m = utf8.DecodeRuneInString(s[i:])
				advance = 0
				skip = m
			}
			goto wrap
		}
		continue
	wrap:
		writeStringNoWrap(w, s[done:i])
		done = i + skip
		w.Pos.X = bounds.Min.X
		posX = w.Pos.X + advance
		if w.Wrap == WrapNewLine {
			w.Pos.Y += height
		}
	}
	writeStringNoWrap(w, s[done:])
	return len(s), nil
}

func drawRune(w *TextWriter, r rune) {
	mask, origin, advance := w.Face.Glyph(r)
	if mask == nil {
		return
	}
	mr := mask.Bounds()
	dr := mr.Add(w.Pos.Add(w.Offset).Sub(origin))
	w.Area.Draw(dr, w.Color, image.Point{}, mask, mr.Min, draw.Over)
	w.Pos.X += advance
}
