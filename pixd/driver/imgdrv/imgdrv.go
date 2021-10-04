// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package imgdrv provides a driver based on Go image and image/draw packages.
// It allows to use pix package without any real display, mainly for tests.
package imgdrv

import (
	"image"
	"image/color"
	"image/draw"
)

type Driver struct {
	img  draw.Image
	fill image.Uniform
}

func New(img draw.Image) *Driver {
	return &Driver{img: img}
}

func (d *Driver) Size() image.Point {
	return d.img.Bounds().Size()
}

func (d *Driver) Draw(r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point, op draw.Op) {
	draw.DrawMask(d.img, r.Add(d.img.Bounds().Min), src, sp, mask, mp, op)
}

func (d *Driver) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	d.fill.C = color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (d *Driver) Fill(r image.Rectangle) {
	d.Draw(r, &d.fill, image.Point{}, nil, image.Point{}, draw.Over)
}

func (d *Driver) Flush()               {}
func (d *Driver) Err(clear bool) error { return nil }
