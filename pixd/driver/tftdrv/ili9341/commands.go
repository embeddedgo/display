// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

import "github.com/embeddedgo/display/pixd/driver/tftdrv/internal/philips"

// ILI9341 commands
const (
	NOP       = philips.NOP
	SWRESET   = philips.SWRESET
	RDDIDIF   = philips.RDDIDIF
	RDDST     = philips.RDDST
	RDDPM     = 0x0A
	RDDMADCTL = 0x0B
	RDDCOLMOD = 0x0C
	RDDIM     = 0x0D
	RDDSM     = 0x0E
	RDDSDR    = 0x0F
	SLPIN     = philips.SLPIN
	SLPOUT    = philips.SLPOUT
	PTLON     = philips.PTLON
	NORON     = philips.NORON
	DINVOFF   = philips.INVOFF
	DINVON    = philips.INVON
	GAMMASET  = 0x26
	DISPOFF   = philips.DISPOFF
	DISPON    = philips.DISPON
	CASET     = philips.CASET
	PASET     = philips.PASET
	RAMWR     = philips.RAMWR
	RAMRD     = 0x2E
	PLTAR     = philips.PTLAR
	VSCRDEF   = philips.VSCRDEF
	TEOFF     = philips.TEOFF
	TEON      = philips.TEON
	MADCTL    = philips.MADCTL
	VSCRSADD  = philips.SEP
	IDMOFF    = philips.IDMOFF
	IDMON     = philips.IDMON
	PIXSET    = philips.COLMOD
	WRCONT    = 0x3C
	RDCONT    = 0x3E
	TSCLSET   = 0x44
	TSCLGET   = 0x45
	WRDISBV   = 0x51
	RDDISBV   = 0x52
	WRCTRLD   = 0x53
	RDCTRLD   = 0x54
	WRCABC    = 0x55
	RDCABC    = 0x56
	WRCABCMB  = 0x5E
	RDCABCMB  = 0x5F
	IFMODE    = 0xB0
	FRMCTR1   = 0xB1
	FRMCTR2   = 0xB2
	FRMCTR3   = 0xB3
	INVTR     = 0xB4
	PRCTR     = 0xB5
	DFUNCTR   = 0xB6
	ETMOD     = 0xB7
	BLCTR1    = 0xB8
	BLCTR2    = 0xB9
	BLCTR3    = 0xBA
	BLCTR4    = 0xBB
	BLCTR5    = 0xBC
	BLCTR7    = 0xBE
	BLCTR8    = 0xBF
	PWCTR1    = 0xC0
	PWCTR2    = 0xC1
	VMCTR1    = 0xC5
	VMCTR2    = 0xC7
	PWCTRA    = 0xCB
	PWCTRB    = 0xCF
	NVMWR     = 0xD0
	NVMPKEY   = 0xD1
	RDNVM     = 0xD2
	RDID4     = 0xD3
	RDID1     = philips.RDID1
	RDID2     = philips.RDID2
	RDID3     = philips.RDID3
	PGAMCTRL  = 0xE0
	NGAMCTRL  = 0xE1
	DGAMCTRL1 = 0xE2
	DGAMCTRL2 = 0xE3
	DRVTIM    = 0xE8
	DRVTIMA   = 0xE9
	DRVTIMB   = 0xEA
	PONSEQ    = 0xED
	IFCTL     = 0xF6
	EN3G      = 0xF2
	PUMPRT    = 0xF7
	SFD       = philips.SFD // undocumented
)

// MADCTL arguments
const (
	MH  = 1 << 2 // horizontal refresh order
	BGR = 1 << 3 // BGR color filter panel
	ML  = 1 << 4 // vertical refresh order
	MV  = 1 << 5 // row/column exchange
	MX  = 1 << 6 // column address order
	MY  = 1 << 7 // row address order
)

// PIXSET arguments
const (
	MCU16 = 0x05 // set 16-bit 565 pixel format for MCU interface
	MCU18 = 0x06 // set 18-bit 666 pixel format for MCU interface
	RGB16 = 0x50 // set 16-bit 565 pixel format for RGB interface
	RGB18 = 0x60 // set 18-bit 666 pixel format for RGB interface
)
