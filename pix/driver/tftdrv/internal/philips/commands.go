// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package philips

import (
	"image"

	"github.com/embeddedgo/display/pix/driver/tftdrv"
)

// PCF8833 commands
const (
	NOP       = 0x00 // no operation
	SWRESET   = 0x01 // software reset
	BSTROFF   = 0x02 // booster voltage off
	BSTRON    = 0x03 // booster voltage on
	RDDIDIF   = 0x04 // read display identification
	RDDST     = 0x09 // read display status
	SLPIN     = 0x10 // sleep in
	SLPOUT    = 0x11 // sleep out
	PTLON     = 0x12 // partial mode on
	NORON     = 0x13 // normal Display mode on
	INVOFF    = 0x20 // display inversion off
	INVON     = 0x21 // display inversion on
	DALO      = 0x22 // all pixel off
	DAL       = 0x23 // all pixel on
	SETCON    = 0x25 // set contrast
	DISPOFF   = 0x28 // display off
	DISPON    = 0x29 // display on
	CASET     = 0x2A // column address set
	PASET     = 0x2B // page address set
	RAMWR     = 0x2C // memory write
	RGBSET    = 0x2D // colour set
	RAMRD     = 0x2E // not supported by PCF8833 but exists in many derived command sets
	PTLAR     = 0x30 // partial area
	VSCRDEF   = 0x33 // vertical scroll definition
	TEOFF     = 0x34 // tearing line off
	TEON      = 0x35 // tearing line on
	MADCTL    = 0x36 // memory data access control
	SEP       = 0x37 // set Scroll Entry Point
	IDMOFF    = 0x38 // Idle mode off
	IDMON     = 0x39 // Idle mode on
	COLMOD    = 0x3A // interface pixel format
	SETVOP    = 0xB0 // set VOP
	BRS       = 0xB4 // Bottom Row Swap
	TRS       = 0xB6 // Top Row Swap
	FINV      = 0xB9 // super Frame INVersion
	DOR       = 0xBA // Data ORder
	TCDFE     = 0xBD // enable/disable DF temp comp
	TCVOPE    = 0xBF // enable or disable VOP temp comp
	EC        = 0xC0 // internal or external oscillator
	SETMUL    = 0xC2 // set multiplication factor
	TCVOPAB   = 0xC3 // set TCVOP slopes A and B
	TCVOPCD   = 0xC4 // set TCVOP slopes C and D
	TCDF      = 0xC5 // set divider frequency
	DF8COLOUR = 0xC6 // set divider frequency 8-colour mode
	SETBS     = 0xC7 // set bias system
	RDTEMP    = 0xC8 // temperature read back
	NLI       = 0xC9 // N-Line Inversion
	RDID1     = 0xDA // read ID1
	RDID2     = 0xDB // read ID2
	RDID3     = 0xDC // read ID3
	SFD       = 0xEF // select factory defaults
	ECM       = 0xF0 // enter calibration mode
	OTPSHTIN  = 0xF1 // shift data in OTP shift registers
)

// MADCTL arguments
const (
	BGR = 1 << 3 // BGR color order
	LAO = 1 << 4 // line address order (bottom to top)
	V   = 1 << 5 // vertical RAM write; in Y direction
	MX  = 1 << 6 // mirror X
	MY  = 1 << 7 // mirror Y
)

func StartWrite8(dci tftdrv.DCI, reg *tftdrv.Reg, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	reg.Xarg[0] = uint8(r.Min.X)
	reg.Xarg[1] = uint8(r.Max.X)
	dci.WriteBytes(reg.Xarg[:2])
	dci.Cmd(PASET)
	reg.Xarg[0] = uint8(r.Min.Y)
	reg.Xarg[1] = uint8(r.Max.Y)
	dci.WriteBytes(reg.Xarg[:2])
	dci.Cmd(RAMWR)
}

func StartRead8(dci tftdrv.DCI, reg *tftdrv.Reg, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	reg.Xarg[0] = uint8(r.Min.X)
	reg.Xarg[1] = uint8(r.Max.X)
	dci.WriteBytes(reg.Xarg[:2])
	dci.Cmd(PASET)
	reg.Xarg[0] = uint8(r.Min.Y)
	reg.Xarg[1] = uint8(r.Max.Y)
	dci.WriteBytes(reg.Xarg[:2])
	dci.Cmd(RAMRD)
}

func StartWrite16(dci tftdrv.DCI, reg *tftdrv.Reg, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	reg.Xarg[0] = uint8(r.Min.X >> 8)
	reg.Xarg[1] = uint8(r.Min.X)
	reg.Xarg[2] = uint8(r.Max.X >> 8)
	reg.Xarg[3] = uint8(r.Max.X)
	dci.WriteBytes(reg.Xarg[:])
	dci.Cmd(PASET)
	reg.Xarg[0] = uint8(r.Min.Y >> 8)
	reg.Xarg[1] = uint8(r.Min.Y)
	reg.Xarg[2] = uint8(r.Max.Y >> 8)
	reg.Xarg[3] = uint8(r.Max.Y)
	dci.WriteBytes(reg.Xarg[:])
	dci.Cmd(RAMWR)
}

func StartRead16(dci tftdrv.DCI, reg *tftdrv.Reg, r image.Rectangle) {
	r.Max.X--
	r.Max.Y--
	dci.Cmd(CASET)
	reg.Xarg[0] = uint8(r.Min.X >> 8)
	reg.Xarg[1] = uint8(r.Min.X)
	reg.Xarg[2] = uint8(r.Max.X >> 8)
	reg.Xarg[3] = uint8(r.Max.X)
	dci.WriteBytes(reg.Xarg[:])
	dci.Cmd(PASET)
	reg.Xarg[0] = uint8(r.Min.Y >> 8)
	reg.Xarg[1] = uint8(r.Min.Y)
	reg.Xarg[2] = uint8(r.Max.Y >> 8)
	reg.Xarg[3] = uint8(r.Max.Y)
	dci.WriteBytes(reg.Xarg[:])
	dci.Cmd(RAMRD)
}

func SetDir(dci tftdrv.DCI, reg *tftdrv.Reg, dir int) {
	rdir := reg.Dir[0]
	if rdir&V != 0 {
		dir = -dir
	}
	switch dir & 3 {
	case 1:
		rdir ^= (V | MX)
	case 2:
		rdir ^= (MX | MY)
	case 3:
		rdir ^= (V | MY)
	}
	reg.Xarg[0] = rdir
	dci.Cmd(MADCTL)
	dci.WriteBytes(reg.Xarg[:1])
	dci.End()
}
