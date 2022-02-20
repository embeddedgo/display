// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"strings"
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

func (w *TextWriter) Write(s []byte) (int, error) {
	n := len(s)
	bounds := w.Area.Bounds()
	height, _ := w.Face.Size()
	for len(s) != 0 {
		k, m := len(s), len(s)
		if w.Wrap != NoWrap {
			if i := bytes.IndexByte(s[:k], '\n'); i >= 0 {
				k = i
				m = i + 1
			}
			free := bounds.Max.X - w.Pos.X - Width(s[:k], w.Face)
			for free < 0 {
				i := -1
				if w.Break == BreakSpace {
					i = bytes.LastIndexFunc(s[:k], unicode.IsSpace)
					if i < 0 && w.Pos.X > bounds.Min.X {
						k = 0
						m = 0
						break // try one more time with w.Pos.X=bounds.Min.X

					}
					m = i + 1
				}
				if i < 0 {
					i = k - 1
					m = i
				}
				free += Width(s[i:k], w.Face)
				k = i
			}
		}
		for i := 0; i < k; {
			r, size := utf8.DecodeRune(s[i:k])
			drawRune(w, r, height, bounds)
			i += size
		}
		if k < len(s) {
			w.Pos.X = bounds.Min.X
			if w.Wrap == WrapNewLine {
				w.Pos.Y += height
			}
		}
		s = s[m:]
	}
	return n, nil
}

func (w *TextWriter) WriteString(s string) (int, error) {
	n := len(s)
	bounds := w.Area.Bounds()
	height, _ := w.Face.Size()
	for len(s) != 0 {
		k, m := len(s), len(s)
		if w.Wrap != NoWrap {
			if i := strings.IndexByte(s[:k], '\n'); i >= 0 {
				k = i
				m = i + 1
			}
			free := bounds.Max.X - w.Pos.X - StringWidth(s[:k], w.Face)
			for free < 0 {
				i := -1
				if w.Break == BreakSpace {
					i = strings.LastIndexFunc(s[:k], unicode.IsSpace)
					if i < 0 && w.Pos.X > bounds.Min.X {
						// try one more time with w.Pos.X=bounds.Min.X
						k = 0
						m = 0
						break // try one more time with w.Pos.X=bounds.Min.X
					}
					m = i + 1
				}
				if i < 0 {
					i = k - 1
					m = i
				}
				free += StringWidth(s[i:k], w.Face)
				k = i
			}
		}
		for _, r := range s[:k] {
			drawRune(w, r, height, bounds)
		}
		if k < len(s) {
			w.Pos.X = bounds.Min.X
			if w.Wrap == WrapNewLine {
				w.Pos.Y += height
			}
		}
		s = s[m:]
	}
	return n, nil
}

func drawRune(w *TextWriter, r rune, height int, bounds image.Rectangle) {
	mask, origin, advance := w.Face.Glyph(r)
	if mask == nil {
		mask, origin, advance = w.Face.Glyph(0)
		if mask == nil {
			return
		}
	}
	mr := mask.Bounds()
	dr := mr.Add(w.Pos.Add(w.Offset).Sub(origin))
	w.Area.Draw(dr, w.Color, image.Point{}, mask, mr.Min, draw.Over)
	w.Pos.X += advance
}
