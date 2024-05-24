// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package displays contains initialization functions for the most popular
// displays.
//
// If you do not found a function for your display here do not worry. If your
// display uses a display controller that had a driver in ../driver/tftdrv you
// can easily initialize your display this way:
//
//	drv := driver.New(dci) // or NewOver(dci) if your display supports reading
//	drv.Init(driver.InitCommands)
//	disp := pix.NewDisplay(drv)
//
// which is basically what all functions in this package do internally.
package displays

import (
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/driver/tftdrv"
)

type Param struct {
	MaxReadClk  int
	MaxWriteClk int
	New         func(dci tftdrv.DCI) *pix.Display
}
