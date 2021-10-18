// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import (
	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/internal/philips"
)

// New returns a new driver that implements pixd.Driver interface. It uses
// write-only DCI so the alpha blending is slow and reduced to 1-bit resolution.
// Use NewOver if the full-fledged Porter-Duff composition is required and the
// display supports reading from its frame memory.
func New(dci tftdrv.DCI) *tftdrv.Driver {
	return tftdrv.New(dci, 240, 320, philips.StartWrite16, pixSet, MCU16, MCU18)
}
