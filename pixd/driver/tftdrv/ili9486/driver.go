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
	return tftdrv.New(
		dci,
		320, 480,
		tftdrv.W18,
		philips.StartWrite16,
		nil,
		nil,
	)
}

// NewOver returns a new tftdrv.DriverOver that implements pixd.Driver
// interface. The returned driver fully supports the draw.Over operator but
// requires reading pixel data from the frame memory. If the display has
// write-only interface use New instead.
func NewOver(dci tftdrv.RDCI) *tftdrv.DriverOver {
	return tftdrv.NewOver(
		dci,
		320, 480,
		tftdrv.W18|tftdrv.R18,
		philips.StartWrite16,
		read,
		nil,
		nil,
	)
}

func read(dci tftdrv.RDCI, xarg *[4]byte, r image.Rectangle, buf []byte) {
	philips.StartRead16(dci, xarg, r)
	dci.ReadBytes(buf)
	dci.End()

	// Workaround for the undocumented behavior. The read command probably does
	// not stop immediately after deaserting CSN (need as much as 1.4ms).
	time.Sleep(1400 * time.Microsecond)

	// In case of SYSTICK based system timer the above sleep lasts 2 ms and the
	// drawing is 1.34 times slower compared with the busy wait below:
	//
	//	t0 := time.Now()
	//	for time.Now().Sub(t0) < 1400*time.Microsecond {
	//		runtime.Gosched()
	//	}
}
