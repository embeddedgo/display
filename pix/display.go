// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/color"
)

// Display represents a pixmap based display.
type Display struct {
	drv       Driver
	lastColor color.Color
	swapWH    bool
}

// NewDisplay returns a new itialized display.
func NewDisplay(drv Driver) *Display {
	return &Display{drv: drv}
}

// Driver returns the used driver
func (d *Display) Driver() Driver {
	return d.drv
}

// Flush ensures that everything drawn so far has been sent to the display
// controller. Using it is not required if the driver draws directly to the
// display without any buffering. Even then it is worth using Flush for code
// portability.
func (d *Display) Flush() {
	d.drv.Flush()
}

// Err returns the saved error and clears it if the clear is true.
func (d *Display) Err(clear bool) error {
	return d.drv.Err(clear)
}

// Bounds returns the bounds of the display
func (d *Display) Bounds() image.Rectangle {
	width, height := d.drv.Dim()
	if d.swapWH {
		return image.Rectangle{Max: image.Point{height, width}}
	}
	return image.Rectangle{Max: image.Point{width, height}}
}

func (d *Display) NewArea(r image.Rectangle) *Area {
	a := &Area{disp: d, color: color.Alpha{255}, rect: r.Canon()}
	a.updateBounds()
	return a
}
