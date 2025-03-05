// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package displays

import (
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
	"github.com/embeddedgo/display/pix/driver/tftdrv/ili9486"
)

func newMSP4022(dci tftdrv.DCI) *pix.Display {
	drv := ili9486.NewOver(dci, 320, 480)
	drv.Init(ili9486.MSP4022)
	return pix.NewDisplay(drv)
}

// MSP4022 4.0" TFT LCD SPI Module - ILI9486
func MSP4022_4i0_320x480_TFT_ILI9486() Param {
	return Param{
		ili9486.MaxSPIReadClock,
		ili9486.MaxSPIWriteClock,
		newMSP4022,
	}
}
