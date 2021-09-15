// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

// The commands below can be used to directly interract with a display
// controller..
const (
	NOP       = 0x00
	SWRESET   = 0x01
	RDDMADCTL = 0x0B
	RDDCOLMOD = 0x0C
	SLPIN     = 0x10
	SLPOUT    = 0x11
	GAMMASET  = 0x26
	DISPOFF   = 0x28
	DISPON    = 0x29
	CASET     = 0x2A
	PASET     = 0x2B
	RAMWR     = 0x2C
	RAMRD     = 0x2E
	MADCTL    = 0x36
	VSCRSADD  = 0x37
	PIXSET    = 0x3A
	WRCONT    = 0x3C
	RDCONT    = 0x3E
	FRMCTR1   = 0xB1
	DFUNCTR   = 0xB6
	PWCTR1    = 0xC0
	PWCTR2    = 0xC1
	VMCTR1    = 0xC5
	VMCTR2    = 0xC7
	PWCTRA    = 0xCB
	PWCTRB    = 0xCF
	GMCTRP1   = 0xE0
	GMCTRN1   = 0xE1
	DRVTIM    = 0xE8
	DRVTIMA   = 0xE9
	DRVTIMB   = 0xEA
	PONSEQ    = 0xED
	PUMPRT    = 0xF7
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
	PF16 = 0x55 // 16-bit 565 pixel format
	PF18 = 0x66 // 18-bit 666 pixel format
)
