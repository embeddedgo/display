// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1351

func New(dci tftdrv.DCI) *tftdrv.Driver {
	return tftdrv.New(dci, 240, 320, tftdrv.W16|tftdrv.W18, ctrl)
}

var (
	ctrl = &tftdrv.Ctrl{
		StartWrite: epson.StartWrite8,
		SetPF:      setPF,
		SetDir:     setDir,
	}
)

func setPF(dci tftdrv.DCI, parg *[1]byte, pixSize int) {
	pf := byte(0)
	if pixSize == 3 {
		pf = MCU18
	}
	if parg[0] != pf {
		parg[0] = pf
		dci.Cmd(PIXSET)
		dci.WriteBytes(parg[:])
	}
}


setDir