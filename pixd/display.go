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
	drvBounds image.Rectangle
	tr        image.Point // translation to the driver coordinates
	bounds    image.Rectangle
	lastColor color.Color
}

// NewDisplay returns a new itialized display.
func NewDisplay(drv Driver) *Display {
	d := new(Display)
	d.drv = drv
	d.drvBounds = drv.SetDir(0)
	d.bounds = d.drvBounds
	return d
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

// Rect returns the rectangle set by SetRect.
func (d *Display) Rect() image.Rectangle {
	return d.bounds.Add(d.tr)
}

// SetRect sets the rectangle covered by the display in the frame memory
// provided the driver. The r may exceed the dimensions of the frame memory,
// only the intersection will be used for drawing.
func (d *Display) SetRect(r image.Rectangle) {
	d.tr = r.Min.Sub(d.bounds.Min)
	d.bounds.Max = d.bounds.Min.Add(r.Size())
}

// Bounds returns the bounds of the display
func (d *Display) Bounds() image.Rectangle {
	return d.bounds
}

// SetOrigin sets the coordinate of the upper left corner of the display. It
// translates internal coordinate system of the display in a way that the
// d.Bounds().Min = origin.
func (d *Display) SetOrigin(origin image.Point) {
	d.tr = d.tr.Add(d.bounds.Min).Sub(origin)
	d.bounds.Max = origin.Add(d.bounds.Size())
	d.bounds.Min = origin
}

// SetDir sets the display direction by rotate its default coordinate system
// by dir*90 degrees. The positive number means clockwise rotation, the
// negative one means counterclockwise rotation.
//
// After SetDir the coordinates of all areas that use this display must be
// re-set (usng SetRect method) and their content should be redrawn. Use
// a.SetRect(a.Rect()) if the coordinates should remain the same.
func (d *Display) SetDir(dir int) {
	d.drvBounds = d.drv.SetDir(dir)
}

// NewArea provides a convenient way to create a drawing area on the display.
// Use NewArea for areas that occupy more than one display.
func (d *Display) NewArea(r image.Rectangle) *Area {
	a := new(Area)
	a.bounds = image.Rectangle{Max: r.Size()} // area's default origin is (0,0)
	a.color = color.Alpha{255}
	a.ad.disp = d
	a.SetRect(r)
	return a
}
