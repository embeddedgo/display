// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/internal/philips"
)

// New returns a new tftdrv.Driver that implements pixd.Driver interface. The
// returned driver works without reading the frame memory so the alpha blending
// is slow and reduced to 1-bit resolution. Use NewOver if the display supports
// reading pixel data and the full-fledged Porter-Duff composition is required.
func New(dci tftdrv.DCI) *tftdrv.Driver {
	return tftdrv.New(
		dci,
		240, 320,
		philips.StartWrite16, pixSet,
		MCU16, MCU18,
	)
}

// NewOver returns a new tftdrv.DriverOver that implements pixd.Driver
// interface. The returned driver fully supports the draw.Over operator but
// requires reading pixel data from the frame memory. If the display has
// write-only interface use New instead.
func NewOver(dci tftdrv.RDCI) *tftdrv.DriverOver {
	return tftdrv.NewOver(
		dci,
		240, 320,
		philips.StartRead16, philips.StartWrite16, pixSet,
		MCU16, MCU18,
	)
}

func pixSet(dci tftdrv.DCI, oldpf *[1]byte, newpf byte) {
	if oldpf[0] != newpf {
		oldpf[0] = newpf
		dci.Cmd(PIXSET)
		dci.WriteBytes(oldpf[:])
	}
}
