// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9486

import (
	"image"
	"time"

	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/internal/philips"
)

// New returns a new tftdrv.Driver that implements pixd.Driver interface. The
// returned driver works without reading the frame memory so the alpha blending
// is slow and reduced to 1-bit resolution. Use NewOver if the display supports
// reading pixel data and the full-fledged Porter-Duff composition is required.
func New(dci tftdrv.DCI) *tftdrv.Driver {
	return tftdrv.New(dci, 320, 480, tftdrv.W18, ctrl)
}

// NewOver returns a new tftdrv.DriverOver that implements pixd.Driver
// interface. The returned driver fully supports the draw.Over operator but
// requires reading pixel data from the frame memory. If the display has
// write-only interface use New instead.
//
// ILI9486 controller seems to require a 2 ms delay before start writing to
// frame memory just after reading from it. It's undocumented behavior and
// slows down the draw.Over operation much.
func NewOver(dci tftdrv.DCI) *tftdrv.DriverOver {
	return tftdrv.NewOver(dci, 320, 480, tftdrv.W18|tftdrv.R18, ctrlOver)
}

var (
	ctrl = &tftdrv.Ctrl{
		StartWrite: philips.StartWrite16,
		SetDir:     philips.SetDir,
	}

	ctrlOver = &tftdrv.Ctrl{
		StartWrite: philips.StartWrite16,
		Read:       read,
		SetDir:     philips.SetDir,
	}
)

func read(dci tftdrv.DCI, reg *tftdrv.Reg, r image.Rectangle, buf []byte) {
	philips.StartRead16(dci, reg, r)
	dci.ReadBytes(buf)
	dci.End()

	// Workaround for the undocumented behavior. Deaserting CSN probably does
	// not stop reading immediately (it takes as much as 2 ms).
	time.Sleep(2 * time.Millisecond)
}
