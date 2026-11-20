// Copyright 2026 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package displays

import (
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/fbdrv"
	"github.com/embeddedgo/display/pix/driver/fbdrv/ssd1306"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
)


// GNPE Monochrome 0.41" 72x40 OLED Graphic Display

var GNPE_0i41_72x40_OLED_SSD1315 = Def{
	0,
	ssd1306.MaxSPIWriteClock,
	func(dci tftdrv.DCI) *pix.Display {
		fb := ssd1306.New(dci, 72, 40)
		fb.Init(ssd1306.GNPE72x40)
		return pix.NewDisplay(fbdrv.NewMono(fb))
	},
}
