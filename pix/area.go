// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pix

import (
	"image"
	"image/color"

	"github.com/embeddedgo/display/font"
	"github.com/embeddedgo/display/images"
)

type areaDisplay struct {
	disp    *Display
	tr      image.Point     // translation to the driver coordinates
	visible image.Rectangle // driver coordinates
	link    *areaDisplay
}

// Area is the drawing area on the display. It has its own coordinates
// independent of its position on the display. Only one goroutine can use an
// area at the same time.
type Area struct {
	bounds image.Rectangle // area coordinates; always use Bounds() for drawing
	tr     image.Point     // translation to the display coordinates
	color  color.Color     // drawing color
	misrc  images.Mirror
	mimask images.Mirror
	ad     areaDisplay // linked list of displays covered by this area
}

// NewArea creates new area that covers the given rectangle on the displays.
func NewArea(r image.Rectangle, displays ...*Display) *Area {
	a := new(Area)
	a.bounds = image.Rectangle{Max: r.Size()} // default origin is (0,0)
	a.color = color.Alpha{255}
	ads := make([]areaDisplay, len(displays)-1) // panics if len(displays) == 0
	for i, ad := 0, &a.ad; ; i, ad = i+1, ad.link {
		ad.disp = displays[i]
		if i >= len(ads) {
			break
		}
		ad.link = &ads[i]
	}
	a.SetRect(r)
	return a
}

// Rect returns the rectangle set by SetRect.
func (a *Area) Rect() image.Rectangle {
	return a.bounds.Add(a.tr)
}

// SetRect sets the rectangle covered by the area on the displays.
func (a *Area) SetRect(r image.Rectangle) {
	a.tr = r.Min.Sub(a.bounds.Min)
	a.bounds.Max = a.bounds.Min.Add(r.Size())
	for ad := &a.ad; ad != nil; ad = ad.link {
		ad.disp.mt.Lock()
		dispTr := ad.disp.tr
		drvBounds := ad.disp.drvBounds
		ad.disp.mt.Unlock()
		ad.tr = a.tr.Add(dispTr)
		ad.visible = r.Add(dispTr).Intersect(drvBounds)
	}
}

// Flush runs through all the displays covered by a and calls thier Flush
// methods.
func (a *Area) Flush() {
	for ad := &a.ad; ad != nil; ad = ad.link {
		ad.disp.Flush()
	}
}

// Err runs through all the displays covered by a and calls thier Err
// methods until it encounters non-nil error.
func (a *Area) Err(clear bool) error {
	for ad := &a.ad; ad != nil; ad = ad.link {
		if err := ad.disp.Err(clear); err != nil {
			return err
		}
	}
	return nil
}

const (
	MI = images.MI // identity (no operation)
	MV = images.MV // swap X with Y
	MX = images.MX // mirror X axis
	MY = images.MY // mirror Y axis
)

// Mirror returns current mirror drawing mode.
func (a *Area) Mirror() int {
	return a.misrc.Mode
}

// SetMirror sets mirror drawing mode. It affects all drawing methods and
// Bounds. Other Area's methods are unaffected.
func (a *Area) SetMirror(mvxy int) {
	a.misrc.Mode = mvxy
	a.mimask.Mode = mvxy
}

// Bounds return the area bounds.
func (a *Area) Bounds() image.Rectangle {
	if a.misrc.Mode == MI {
		return a.bounds
	}
	r := a.bounds
	if a.misrc.Mode&MX != 0 {
		r.Min.X, r.Max.X = -r.Max.X, -r.Min.X
	}
	if a.misrc.Mode&MY != 0 {
		r.Min.Y, r.Max.Y = -r.Max.Y, -r.Min.Y
	}
	if a.misrc.Mode&MV != 0 {
		r.Min.X, r.Min.Y = r.Min.Y, r.Min.X
		r.Max.X, r.Max.Y = r.Max.Y, r.Max.X
	}
	return r
}

// Returns the coordinate of the upper left corner of the area in the area's own
// coordinate system.  The origin equals to a.Bounds().Min in MI mirror mode.
func (a *Area) Origin() image.Point {
	return a.bounds.Min
}

// SetOrigin sets the coordinate of the upper left corner of the area. It does
// not affect the rectangle covered by the area but translates the area's own
// coordinate system in a way that the a.Bounds().Min = origin in MI mirror
// mode.
func (a *Area) SetOrigin(origin image.Point) {
	delta := a.bounds.Min.Sub(origin)
	a.bounds.Max = origin.Add(a.bounds.Size())
	a.bounds.Min = origin
	for ad := &a.ad; ad != nil; ad = ad.link {
		ad.tr = ad.tr.Add(delta)
	}
}

// SetColor sets the color used by drawing methods.
func (a *Area) SetColor(c color.Color) {
	a.color = c
}

// SetColorRGBA is equivalent of SetColor(color.RGBA{r, g, b, alpha}). Notice
// that r, g, b must be alpha-premultiplied, e.g. they must be less than or
// equal to alpha.
func (a *Area) SetColorRGBA(r, g, b, alpha uint8) {
	a.color = color.RGBA{r, g, b, alpha}
}

// Color returns the color used by drawing methods.
func (a *Area) Color() color.Color {
	return a.color
}

// NewTextWriter provides a conveniet way to create new TextWriter. It can be
// used in place of the following set of statements:
//	w := new(TextWriter)
//	w.Area = a
//	w.Face = f
//	w.Color = &image.Uniform{a.Color()}
//	_, w.Pos.Y = f.Size() // ascent
func (a *Area) NewTextWriter(f font.Face) *TextWriter {
	_, ascent := f.Size()
	return &TextWriter{
		Area:   a,
		Face:   f,
		Color:  &image.Uniform{a.color},
		Pos:    a.Bounds().Min,
		Offset: image.Pt(0, ascent),
		Wrap:   WrapNewLine,
	}
}
