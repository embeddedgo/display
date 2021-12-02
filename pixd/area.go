// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import (
	"image"
	"image/color"
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
	ad     areaDisplay
	color  color.Color
	bounds image.Rectangle // area coordinates
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
	a.ad.disp.mt.Lock()
	dispTr := a.ad.disp.tr
	a.ad.disp.mt.Unlock()
	tr := a.ad.tr.Sub(dispTr) // translation to the display coordinates
	return a.bounds.Add(tr)
}

// SetRect sets the rectangle covered by the area on the displays.
func (a *Area) SetRect(r image.Rectangle) {
	tr := r.Min.Sub(a.bounds.Min)
	a.bounds.Max = a.bounds.Min.Add(r.Size())
	for ad := &a.ad; ad != nil; ad = ad.link {
		ad.disp.mt.Lock()
		dispTr := ad.disp.tr
		drvBounds := ad.disp.drvBounds
		ad.disp.mt.Unlock()
		ad.tr = tr.Add(dispTr)
		ad.visible = r.Add(dispTr).Intersect(drvBounds)
	}
}

func (a *Area) Flush() {
	for ad := &a.ad; ad != nil; ad = ad.link {
		ad.disp.Flush()
	}
}

// Bounds return the area bounds.
func (a *Area) Bounds() image.Rectangle {
	return a.bounds
}

// SetOrigin sets the coordinate of the upper left corner of the area. It does
// not affect the rectangle covered by the area but translates the area own
// coordinate system in a way that the a.Bounds().Min = origin.
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
func (a *Area) NewTextWriter(f FontFace) *TextWriter {
	_, ascent := f.Size()
	return &TextWriter{
		Area:  a,
		Face:  f,
		Color: &image.Uniform{a.color},
		Pos:   image.Pt(0, ascent),
	}
}
