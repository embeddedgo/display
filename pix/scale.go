// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/color"
)

type ScaledUp struct {
	Img image.Image
	Mul int
}

func NewScaledUp(img image.Image, mul int) *ScaledUp {
	return &ScaledUp{img, mul}
}

func (p *ScaledUp) ColorModel() color.Model {
	return p.Img.ColorModel()
}

func (p *ScaledUp) Bounds() image.Rectangle {
	r := p.Img.Bounds()
	r.Min = r.Min.Mul(p.Mul)
	r.Max = r.Max.Mul(p.Mul)
	return r
}

func (p *ScaledUp) At(x, y int) color.Color {
	println(x, y)
	return p.Img.At(x/p.Mul, y/p.Mul)
}

type ScaledFont struct {
	face FontFace
	img  ScaledUp
}

func NewScaledFont(face FontFace, mul int) *ScaledFont {
	sf := new(ScaledFont)
	sf.face = face
	sf.img.Mul = mul
	return sf
}

func (sf *ScaledFont) Size() (height, ascent int) {
	height, ascent = sf.face.Size()
	return height * sf.img.Mul, ascent * sf.img.Mul
}

func (sf *ScaledFont) Advance(r rune) int {
	return sf.face.Advance(r) * sf.img.Mul
}

func (sf *ScaledFont) Glyph(r rune) (img image.Image, origin image.Point, advance int) {
	img, origin, advance = sf.face.Glyph(r)
	sf.img.Img = img
	return &sf.img, origin.Mul(sf.img.Mul), advance * sf.img.Mul

}
