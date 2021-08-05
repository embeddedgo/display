// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tft

import (
	"image"
)

type Driver interface {
	FillRect(x0, y0, x1, y1, color uint16)
	DrawImage(x, y int, img image.Image)
}

type FillRect func(dci DCI, x0, y0, x1, y1, wxh int, color uint16)

type Display struct {
	dci      DCI
	fillRect FillRect
	width    uint16
	height   uint16
	swapWH   bool
}

// NewDisplay returns new nitialized display struct..
func NewDisplay(dci DCI, width, height int, fillRect FillRect) *Display {
	return &Display{
		dci:      dci,
		fillRect: fillRect,
		width:    uint16(width),
		height:   uint16(height),
	}
}

// DCI allows to direct access to the internal DCI.
func (d *Display) DCI() DCI {
	return d.dci
}

// Err returns and clears internal error variable.
func (d *Display) Err(clear bool) error {
	return d.dci.Err(clear)
}

// Bounds returns the bounds of the display
func (d *Display) Bounds() image.Rectangle {
	if d.swapWH {
		return image.Rectangle{Max: image.Pt(int(d.height), int(d.width))}
	}
	return image.Rectangle{Max: image.Pt(int(d.width), int(d.height))}
}

func (d *Display) Area(r image.Rectangle) Area {
	a := Area{disp: d, rect: r.Canon()}
	a.updateBounds()
	return a
}

func (d *Display) NewArea(r image.Rectangle) *Area {
	a := new(Area)
	*a = d.Area(r)
	return a
}
