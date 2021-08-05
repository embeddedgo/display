// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9341

// These commands can be used to directly interract with display controller
// using DCI.
const (
	NOP      = 0x00
	SWRESET  = 0x01
	SLPIN    = 0x10
	SLPOUT   = 0x11
	GAMMASET = 0x26
	DISPOFF  = 0x28
	DISPON   = 0x29
	CASET    = 0x2A
	PASET    = 0x2B
	RAMWR    = 0x2C
	MADCTL   = 0x36
	VSCRSADD = 0x37
	PIXSET   = 0x3A
	FRMCTR1  = 0xB1
	DFUNCTR  = 0xB6
	PWCTR1   = 0xC0
	PWCTR2   = 0xC1
	VMCTR1   = 0xC5
	VMCTR2   = 0xC7
	PWCTRB   = 0xCF
	GMCTRP1  = 0xE0
	GMCTRN1  = 0xE1
)

// MADCTL arguments
const (
	MH  = 1 << 2 // Horizontal refresh order.
	BGR = 1 << 3 // RGB-BGR order.
	ML  = 1 << 4 // Vertical refresh order.
	MV  = 1 << 5 // Row/column exchange.
	MX  = 1 << 6 // Column address order.
	MY  = 1 << 7 // Row address order.
)

// PIXSET arguments
const (
	PF16 = 0x55 // 16-bit 565 pixel format.
	PF18 = 0x66 // 18-bit 666 pixel format.
)
