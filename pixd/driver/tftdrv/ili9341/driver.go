// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"image"

	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/internal/philips"
)

// New returns a new tftdrv.Driver that implements pixd.Driver interface. The
// returned driver works without reading the frame memory so the alpha blending
// is slow and reduced to 1-bit resolution. Use NewOver if the display supports
// reading pixel data and the full-fledged Porter-Duff composition is required.
func New(dci tftdrv.DCI) *tftdrv.Driver {
	return tftdrv.New(dci, 240, 320, tftdrv.W16|tftdrv.W18H, ctrl)
}

// NewOver returns a new tftdrv.DriverOver that implements pixd.Driver
// interface. The returned driver fully supports the draw.Over operator but
// requires reading pixel data from the frame memory. If the display has
// write-only interface use New instead.
func NewOver(dci tftdrv.DCI) *tftdrv.DriverOver {
	return tftdrv.NewOver(dci, 240, 320, tftdrv.W16|tftdrv.W18H|tftdrv.R18H, ctrlOver)
}

func read(dci tftdrv.DCI, reg *tftdrv.Reg, r image.Rectangle, buf []byte) {
	philips.StartRead16(dci, reg, r)
	dci.ReadBytes(buf)
	dci.End()
}

func setPF(dci tftdrv.DCI, reg *tftdrv.Reg, pixSize int) {
	pf := byte(MCU16)
	if pixSize == 3 {
		pf = MCU18
	}
	if reg.PF[0] != pf {
		reg.PF[0] = pf
		dci.Cmd(PIXSET)
		dci.WriteBytes(reg.PF[:])
	}
}

var (
	ctrl = &tftdrv.Ctrl{
		StartWrite: philips.StartWrite16,
		SetPF:      setPF,
		SetDir:     philips.SetDir,
	}

	ctrlOver = &tftdrv.Ctrl{
		StartWrite: philips.StartWrite16,
		Read:       read,
		SetPF:      setPF,
		SetDir:     philips.SetDir,
	}
)
