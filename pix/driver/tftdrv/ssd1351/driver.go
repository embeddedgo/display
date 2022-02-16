// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1351

import (
	"github.com/embeddedgo/display/pix/driver/tftdrv"
	"github.com/embeddedgo/display/pix/driver/tftdrv/internal/epson"
)

func New(dci tftdrv.DCI) *tftdrv.Driver {
	return tftdrv.New(dci, 128, 128, tftdrv.W16|tftdrv.W18L, ctrl)
}

var (
	ctrl = &tftdrv.Ctrl{
		StartWrite: epson.StartWrite8,
		SetPF:      setPF,
		SetDir:     setDir,
	}
)

func setPF(dci tftdrv.DCI, reg *tftdrv.Reg, pixSize int) {
	rmcd := reg.PF[0] &^ RGB18
	if pixSize == 3 {
		rmcd |= RGB18
	}
	if reg.PF[0] != rmcd {
		reg.PF[0] = rmcd
		reg.Xarg[0] = RMCD
		dci.Cmd(reg.Xarg[:1])
		dci.WriteBytes(reg.PF[:])
	}
}

func setDir(dci tftdrv.DCI, reg *tftdrv.Reg, dir int) {
	rmcd := reg.Dir[0]
	if dir&2 != 0 { // SSD1351 in hardware supports only 180 deg rotation
		rmcd ^= (C127_SEG0 | COMn_COM0)

	}
	reg.PF[0] = rmcd
	reg.Xarg[0] = RMCD
	dci.Cmd(reg.Xarg[:1])
	dci.WriteBytes(reg.PF[:])
	dci.End()
}
