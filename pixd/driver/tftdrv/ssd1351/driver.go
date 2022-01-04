// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssd1351

import (
	"github.com/embeddedgo/display/pixd/driver/tftdrv"
	"github.com/embeddedgo/display/pixd/driver/tftdrv/internal/epson"
)

func New(dci tftdrv.DCI) *tftdrv.Driver {
	return tftdrv.New(dci, 128, 128, tftdrv.W18, ctrl)
}

var (
	ctrl = &tftdrv.Ctrl{
		StartWrite: epson.StartWrite8,
		SetPF:      setPF,
		SetDir:     setDir,
	}
)

func setPF(dci tftdrv.DCI, reg *tftdrv.Reg, pixSize int) {
	dirpf := reg.Dir[0] ^ RGB18
	if pixSize == 3 {
		dirpf |= RGB18
	}
	if reg.Dir[0] != dirpf {
		reg.Dir[0] = dirpf
		dci.Cmd(RMDCLM)
		dci.WriteBytes(reg.Dir[:])
	}
}

func setDir(dci tftdrv.DCI, reg *tftdrv.Reg, dir int) {

}
