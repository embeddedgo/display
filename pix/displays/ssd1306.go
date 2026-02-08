// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package displays

import (
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/fbdrv"
	"github.com/embeddedgo/display/pix/driver/fbdrv/ssd1306"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
)

// Adafruit Monochrome 0.96" 128x64 OLED Graphic Display
var Adafruit_0i96_128x64_OLED_SSD1306 = Def{
	0,
	ssd1306.MaxSPIWriteClock,
	func(dci tftdrv.DCI) *pix.Display {
		fb := ssd1306.New(dci)
		fb.Init(ssd1306.GFX128x64)
		return pix.NewDisplay(fbdrv.NewMono(fb))
	},
}
