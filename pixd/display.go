// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/color"
)

// Display represents a pixmap based display.
type Display struct {
	drv       Driver
	lastColor color.Color
	origin    image.Point
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
	return image.Rectangle{Min: d.origin, Max: d.origin.Add(d.drv.Size())}
}

// SetOrigin sets the coordinate of the upper left corner of the display. It
// translates internal coordinate system of the display in a way that the
// d.Bounds().Min = origin.
func (d *Display) SetOrigin(origin image.Point) {
	d.origin = origin
}

// SetDir sets the display direction by rotate its default coordinate system
// by n*90 degrees. The positive number means clockwise rotation, the
// negative one means counterclockwise rotation.
//
// After SetDir the coordinates of all areas that use this display should be
// re-set (usng SetRect method) and their content should be redrawn. Use
// a.SetRect(a.Rect()) if the coordinates of the area on the rotated display
// should remain the same.
func (d *Display) SetDir(n int) {
	d.drv.SetDir(n)
}

func (d *Display) NewArea(r image.Rectangle) *Area {
	a := new(Area)
	a.bounds = image.Rectangle{Max: r.Size()} // default origin is (0,0)
	a.color = color.Alpha{255}
	a.ad.disp = d
	a.SetRect(r)
	return a
}
